package sshHandler

import (
	"errors"
	"fmt"
	"github.com/pawmot/kabel/sniffer"
	"github.com/phayes/freeport"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
)

type state int

const (
	disconnected state = iota
	connected
)

var (
	urnRegexp = regexp.MustCompile("^[a-z]*://(?P<urn>.*)$")
)

type Actor struct {
	sshPid  int
	state   state
	actionC chan func()
}

func NewSshActor() *Actor {
	actionC := make(chan func())
	actor := Actor{
		sshPid:  -1,
		state:   disconnected,
		actionC: actionC,
	}

	go actor.handleMessages()

	return &actor
}

func (a *Actor) handleMessages() {
	for {
		f := <-a.actionC
		f()
	}
}

func (a *Actor) CreateTunnel(remoteSpec string, dockerEndpoint string) (local sniffer.SshTunnelLocalPort, err error) {
	portC := make(chan int)
	errC := make(chan error)
	a.actionC <- func() {
		if a.state == disconnected {
			p, err := a.tunnel(remoteSpec, dockerEndpoint)
			if err != nil {
				errC <- err
			} else {
				portC <- p
			}
		} else {
			errC <- errors.New("already connected")
		}
	}
	select {
	case port := <-portC:
		return sniffer.SshTunnelLocalPort(port), nil
	case err := <-errC:
		return -1, err
	}
}

func (a *Actor) Close() error {
	errC := make(chan error)
	a.actionC <- func() {
		if a.state == connected {
			err := syscall.Kill(a.sshPid, syscall.SIGTERM)
			if err != nil {
				errC <- err
				return
			}
			log.Println("SSH tunnel closed")
		}

		close(errC)
	}
	return <-errC
}

func (a *Actor) tunnel(remoteSpec string, dockerEndpoint string) (int, error) {
	localPort, err := freeport.GetFreePort()
	if err != nil {
		return -1, err
	}

	log.Println("Running SSH to on local port " + strconv.Itoa(localPort) + "!")
	// TODO use golang ssh client if possible
	// TODO allow for password access
	// TODO allow ssh port setting
	match := urnRegexp.FindStringSubmatch(dockerEndpoint)
	if match == nil {
		return -1, errors.New(fmt.Sprintf("cannot extract a URN from the docker endpoint: %s", dockerEndpoint))
	}
	cmd := exec.Command("/usr/bin/ssh", "-Llocalhost:"+strconv.Itoa(localPort)+":" + match[1], remoteSpec, "-N")
	err = cmd.Start()
	if err != nil {
		return -1, err
	}

	a.sshPid = cmd.Process.Pid
	a.state = connected

	go func() {
		// TODO handle broken connections
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

	return localPort, nil
}
