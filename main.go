package main

import (
	"flag"
	"github.com/fsouza/go-dockerclient"
	"log"
	"gopkg.in/AlecAivazis/survey.v1"
	"bytes"
	"fmt"
	"os"
	"strings"
	"gopkg.in/AlecAivazis/survey.v1/core"
)

func main() {
	core.ErrorIcon = "X"
	core.HelpIcon = "????"
	core.QuestionIcon = "?"
	core.SelectFocusIcon = ">"
	core.MarkedOptionIcon = "[x]"
	core.UnmarkedOptionIcon = "[ ]"

	var dockerEndpoint = flag.String("endpoint", "unix:///var/run/docker.sock", "Docker endpoint to use")
	flag.Parse()
	client, err := docker.NewClient(*dockerEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Fatal(err)
	}

	if len(containers) == 0 {
		fmt.Printf("No containers are running, nothing to do here!")
		os.Exit(0)
	}

	var ids []string
	for _, c := range containers {
		ids = append(ids, c.ID[:12])
	}

	chosenShortId := ""
	prompt := &survey.Select{
		Message: "Select a container to listen to:",
		Options: ids,
	}
	err = survey.AskOne(prompt, &chosenShortId, nil)
	if err != nil {
		log.Fatal(err)
	}

	if chosenShortId == "" {
		log.Fatal("No container chosen")
	}

	exec, err := client.CreateExec(docker.CreateExecOptions{
		AttachStderr: true,
		AttachStdout: true,
		Tty: true,
		Cmd:       []string{"ls", "/sys/class/net"},
		Container: chosenShortId,
	})
	if err != nil {
		log.Fatalf("Couldn't create Exec: %v", err)
	}
	bufout := bytes.NewBufferString("")
	buferr := bytes.NewBufferString("")
	err = client.StartExec(exec.ID, docker.StartExecOptions{OutputStream: bufout, ErrorStream: buferr})
	if err != nil {
		log.Fatalf("Couldn't start Exec: %v", err)
	}
	if buferr.Len() > 0 {
		log.Fatalf("Couldn't read container's interfaces: %s", buferr.String())
	}

	ifaces := strings.Split(strings.Replace(bufout.String(), "  ", " ", -1), " ")
	for idx, i := range ifaces {
		ifaces[idx] = strings.Trim(i, "\n")
	}

	chosenIface := ""
	prompt = &survey.Select{
		Message: "Select an interface to listen to:",
		Options: ifaces,
	}
	err = survey.AskOne(prompt, &chosenIface, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Chosen interface: %s\n", chosenIface)
}
