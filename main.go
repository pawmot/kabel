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
	"io/ioutil"
	"github.com/phayes/freeport"
)

func main() {
	configSurveyIconsCompat()

	dockerClient, sshPidCh := createDockerClient()
	sshPid := <-sshPidCh
	log.Printf("Using SSH pid %d\n", sshPid)
	if sshPid != -1 {
		time.Sleep(20 * time.Second)
	}
	ch := make(chan struct{})
	go func() {
		imageName := "pawmot/tcpdump"
		ctx := context.Background()
		resp, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Close()
		str, err := ioutil.ReadAll(resp)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(string(str))
		ch <- struct{}{}
	}()
	ids := getContainerIds(dockerClient)
	chosenShortId := promptUserForContainerId(ids)
	ifaces := getInterfacesInContainer(dockerClient, chosenShortId)
	chosenIface := promptUserForInterface(ifaces)

	log.Printf("Chosen container id: %s\n", chosenShortId)
	log.Printf("Chosen interface: %s\n", chosenIface)

	<-ch

	ctx := context.Background()
	name := "tcpdump-" + chosenShortId + "-" + chosenIface + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	tdContainer, err := dockerClient.ContainerCreate(ctx, &container.Config{
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
	resp, err := dockerClient.ContainerAttach(attCtx, tdContainer.ID, types.ContainerAttachOptions{
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

	err = dockerClient.ContainerStart(ctx, tdContainer.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Continuing!")
	<-wsClosed

	log.Println("WireShark closed, cleaning up!")

	dur := 30 * time.Second
	err = dockerClient.ContainerStop(ctx, tdContainer.ID, &dur)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(fifoName)
	if err != nil {
		log.Fatal(err)
	}

	if sshPid != -1 {
		err := syscall.Kill(sshPid, syscall.SIGTERM)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Cleanup completed, bye!")
}

func createDockerClient() (*client.Client, <-chan int) {
	var ssh = flag.String("ssh", "", "user@host of the machine that the docker daemon runs on")
	var sshPort = flag.Int("P", 22, "Port to use with ssh")
	var dockerEndpoint = flag.String("endpoint", "unix:///var/run/docker.sock", "Docker endpoint to use")
	flag.Parse()
	var effectiveEnpoint string
	var ch <-chan int
	if len(*ssh) > 0 {
		effectiveEnpoint, ch = createSshTunnel(*ssh, *sshPort, *dockerEndpoint)
	} else {
		var ch1 = make(chan int, 1)
		ch1 <- -1
		ch = ch1
		effectiveEnpoint = *dockerEndpoint
	}
	apiClient, err := client.NewClientWithOpts(client.WithHost(effectiveEnpoint))
	if err != nil {
		log.Fatal(err)
	}
	return apiClient, ch
}

func createSshTunnel(sshSpec string, sshPort int, dockerEndpoint string) (string, <-chan int) {
	localPort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan int)

	go func() {
		log.Println("Running SSH to on local port " + strconv.Itoa(localPort) + "!")
		cmd := exec.Command("/usr/bin/ssh", "-p", strconv.Itoa(sshPort), "-Llocalhost:" + strconv.Itoa(localPort) + ":/var/run/docker.sock", sshSpec, "-N")
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		ch <- cmd.Process.Pid

		if err := cmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					log.Printf("Exit Status: %d", status.ExitStatus())
				}
			} else {
				log.Fatalf("cmd.Wait: %v", err)
			}
		} else {
			log.Println("SSH exited normally")
		}
	}()

	return "tcp://localhost:" + strconv.Itoa(localPort), ch
}

func getContainerIds(dockerClient *client.Client) []string {
	ctx := context.Background()
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
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

func getInterfacesInContainer(dockerClient *client.Client, chosenShortId string) []string {
	ctx := context.Background()
	exec, err := dockerClient.ContainerExecCreate(ctx, chosenShortId, types.ExecConfig{
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
	resp, err := dockerClient.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{Detach: false, Tty: false})
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
