package sshHandler

import (
	"github.com/pawmot/kabel/sniffer"
	"github.com/phayes/freeport"
	"log"
	"strconv"
	"os/exec"
	"syscall"
	"errors"
)

type state int

const (
	disconnected state = iota
	connected
)

type Actor struct {
	sshPid  int
	state   state
	actionc chan func()
}

func NewSshActor() *Actor {
	actor := Actor{
		sshPid: -1,
		state:  disconnected,
	}

	go actor.handleMessages()

	return &actor
}

func (a *Actor) handleMessages() {
	for {
		f := <-a.actionc
		f()
	}
}

func (a *Actor) CreateTunnel(remoteSpec string) (local sniffer.SshTunnelLocalPort, err error) {
	portC := make(chan int)
	errC := make(chan error)
	a.actionc <- func() {
		if a.state == disconnected {
			p, err := a.tunnel(remoteSpec)
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
	a.actionc <- func() {
		if a.state == connected {
			err := syscall.Kill(a.sshPid, syscall.SIGTERM)
			if err != nil {
				errC <- err
				return
			}
		}

		close(errC)
	}
	return <-errC
}

func (a *Actor) tunnel(remoteSpec string) (int, error) {
	localPort, err := freeport.GetFreePort()
	if err != nil {
		return -1, err
	}

	log.Println("Running SSH to on local port " + strconv.Itoa(localPort) + "!")
	// TODO use golang ssh client if possible
	// TODO allow for password access
	// TODO allow ssh port setting
	// TODO allow docker socket setting
	cmd := exec.Command("/usr/bin/ssh", "-Llocalhost:"+strconv.Itoa(localPort)+":/var/run/docker.sock", remoteSpec, "-N")
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
