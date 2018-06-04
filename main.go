package main

import (
	"flag"
	"log"
	"gopkg.in/AlecAivazis/survey.v1"
	"bytes"
	"fmt"
	"os"
	"strings"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"time"
	"strconv"
	"syscall"
	"github.com/docker/docker/client"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"os/exec"
)

func main() {
	configSurveyIconsCompat()

	client := createDockerClient()
	ch := make(chan struct{})
	go func() {
		imageName := "pawmot/tcpdump"
		ctx := context.Background()
		resp, err := client.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			log.Fatal(err)
		}
		resp.Close()
		ch <- struct{}{}
	}()
	ids := getContainerIds(client)
	chosenShortId := promptUserForContainerId(ids)
	ifaces := getInterfacesInContainer(client, chosenShortId)
	chosenIface := promptUserForInterface(ifaces)

	log.Printf("Chosen container id: %s\n", chosenShortId)
	log.Printf("Chosen interface: %s\n", chosenIface)

	<-ch

	ctx := context.Background()
	name := "tcpdump-" + chosenShortId + "-" + chosenIface + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	tdContainer, err := client.ContainerCreate(ctx, &container.Config{
		Image:        "pawmot/tcpdump",
		Env: []string {
			"IF=" + chosenIface,
		},
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		NetworkMode:   container.NetworkMode("container:" + chosenShortId),
		AutoRemove:    true,
		DNS:           []string{},
		DNSOptions:    []string{},
		DNSSearch:     []string{},
		RestartPolicy: container.RestartPolicy{Name: "no", MaximumRetryCount: 0},
	}, &network.NetworkingConfig{

	}, name)

	if err != nil {
		log.Fatal(err)
	}

	fifoName := "/tmp/" + name
	err = syscall.Mkfifo(fifoName, 0666)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Will dump TCP to '" + fifoName + "'")

	wsClosed := make(chan struct{})
	go func() {
		log.Println("Running WireShark on '" + fifoName + "'!")
		cmd := exec.Command("/usr/bin/wireshark", "-k", "-i", fifoName)
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		cmd.Wait()
		wsClosed <- struct{}{}
	}()

	fifo, err := os.OpenFile(fifoName, syscall.O_WRONLY, 0600)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("WireShark connected!")

	log.Println("Attaching...")
	attCtx := context.Background()
	resp, err := client.ContainerAttach(attCtx, tdContainer.ID, types.ContainerAttachOptions{
		//Logs:   true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	log.Println("Attach goroutine finished...")
	defer resp.Close()

	go func() {
		stdcopy.StdCopy(fifo, os.Stderr, resp.Reader)
	}()

	err = client.ContainerStart(ctx, tdContainer.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Continuing!")
	<-wsClosed

	log.Println("WireShark closed, cleaning up!")

	dur := 30 * time.Second
	err = client.ContainerStop(ctx, tdContainer.ID, &dur)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(fifoName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Cleanup completed, bye!")
}

func createDockerClient() (*client.Client) {
	var dockerEndpoint = flag.String("endpoint", "unix:///var/run/docker.sock", "Docker endpoint to use")
	flag.Parse()
	apiClient, err := client.NewClientWithOpts(client.WithHost(*dockerEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	return apiClient
}

func getContainerIds(client *client.Client) []string {
	ctx := context.Background()
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	if len(containers) == 0 {
		fmt.Printf("No containers are running, nothing to do here!")
		os.Exit(0)
	}
	var ids []string
	// TODO that list should contain names, not only IDs
	for _, c := range containers {
		ids = append(ids, c.ID[:12])
	}
	return ids
}

func promptUserForContainerId(ids []string) (string) {
	chosenShortId := ""
	prompt := &survey.Select{
		Message: "Select a container to listen to:",
		Options: ids,
	}
	err := survey.AskOne(prompt, &chosenShortId, nil)
	if err != nil {
		log.Fatal(err)
	}
	if chosenShortId == "" {
		log.Fatal("No container chosen")
	}
	return chosenShortId
}

func getInterfacesInContainer(client *client.Client, chosenShortId string) []string {
	ctx := context.Background()
	exec, err := client.ContainerExecCreate(ctx, chosenShortId, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Tty:          true,
		Cmd:          []string{"ls", "/sys/class/net"},
	})
	if err != nil {
		log.Fatalf("Couldn't create Exec: %v", err)
	}
	bufout := bytes.NewBufferString("")
	buferr := bytes.NewBufferString("")
	resp, err := client.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{Detach: false, Tty: false})
	if err != nil {
		log.Fatalf("Couldn't start Exec: %v", err)
	}
	defer resp.Close()
	stdcopy.StdCopy(bufout, buferr, resp.Reader)
	if buferr.Len() > 0 {
		log.Fatalf("Couldn't read container's interfaces: %s", buferr.String())
	}
	ifaces := strings.Split(strings.Replace(bufout.String(), "  ", " ", -1), " ")
	for idx, i := range ifaces {
		ifaces[idx] = strings.Trim(i, "\n")
	}
	return ifaces
}

func promptUserForInterface(ifaces []string) string {
	chosenIface := ""
	prompt := &survey.Select{
		Message: "Select an interface to listen to:",
		Options: ifaces,
	}
	err := survey.AskOne(prompt, &chosenIface, nil)
	if err != nil {
		log.Fatal(err)
	}
	return chosenIface
}

func configSurveyIconsCompat() {
	core.ErrorIcon = "X"
	core.HelpIcon = "????"
	core.QuestionIcon = "?"
	core.SelectFocusIcon = ">"
	core.MarkedOptionIcon = "[x]"
	core.UnmarkedOptionIcon = "[ ]"
}
