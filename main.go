package main

import (
	"flag"
	"log"
	"gopkg.in/AlecAivazis/survey.v1"
	"fmt"
	"os"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"github.com/pawmot/kabel/dockerHandler"
	"github.com/pawmot/kabel/sshHandler"
	"github.com/pawmot/kabel/wiresharkHandler"
	"github.com/pawmot/kabel/sniffer"
)

func main() {
	configSurveyIconsCompat()

	docker := dockerHandler.NewDockerHandler()
	ssh := sshHandler.NewSshActor()
	wireshark := wiresharkHandler.NewWiresharkClient()
	sniff := sniffer.NewSnifferActor(docker, ssh, wireshark)

	connect(sniff)
	err := sniff.PullImage()
	if err != nil {
		log.Fatal(err)
	}

	containers := getContainers(sniff)
	chosenName, chosenShortId := promptUserForContainer(containers)
	ifaces := getInterfacesInContainer(sniff, chosenShortId)
	chosenIface := promptUserForInterface(ifaces)

	log.Printf("Chosen container: %s (id: %s)\n", chosenName, chosenShortId)
	log.Printf("Chosen interface: %s\n", chosenIface)

	closed, err := sniff.Sniff(chosenShortId, chosenIface)

	<- closed

	log.Println("Bye!")
}

func connect(s *sniffer.Actor) error {
	var ssh = flag.String("ssh", "", "user@host of the machine that the docker daemon runs on")
	// TODO use the port!
	var _ = flag.Int("P", 22, "Port to use with ssh")
	var dockerEndpoint = flag.String("endpoint", "unix:///var/run/docker.sock", "Docker endpoint to use")
	flag.Parse()

	var req sniffer.ConnectionRequest
	if *ssh == "" {
		req = sniffer.DirectConnectionRequest(*dockerEndpoint)
	} else {
		req = sniffer.TunneledConnectionRequest(*dockerEndpoint, *ssh)
	}

	_, err := s.Connect(req)

	return err
}

func getContainers(s *sniffer.Actor) []sniffer.Container {
	containers, err := s.GetContainers()
	if err != nil {
		log.Fatal(err)
	}
	if len(containers) == 0 {
		fmt.Println("No containers are running, nothing to do here!")
		os.Exit(0)
	}

	return containers
}

func promptUserForContainer(containers []sniffer.Container) (name string, id string) {
	var names []string

	for _, c := range containers {
		names = append(names, c.Name)
	}

	chosenName := ""
	prompt := &survey.Select{
		Message: "Select a container to listen to:",
		Options: names,
	}
	err := survey.AskOne(prompt, &chosenName, nil)
	if err != nil {
		log.Fatal(err)
	}
	if chosenName == "" {
		log.Fatal("No container chosen")
	}

	var chosenShortId string
	for _, c := range containers {
		if c.Name == chosenName {
			chosenShortId = c.Id[:12]
			break
		}
	}

	return chosenName, chosenShortId
}

func getInterfacesInContainer(s *sniffer.Actor, chosenShortId string) []string {
	nis, err := s.GetNetworkInterfaces(chosenShortId)
	if err != nil {
		log.Fatal(err)
	}

	var ifaces []string
	for _, i := range nis {
		ifaces = append(ifaces, i.Name)
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
