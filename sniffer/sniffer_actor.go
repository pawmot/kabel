package sniffer

import (
	"errors"
	"github.com/docker/docker/api/types"
	"time"
	"log"
	"os"
	"context"
)

type state int

const (
	created   state = iota
	connected
	sniffing
	closed
)

type Actor struct {
	dockerIface      DockerClient
	sshClient        SshClient
	wiresharkClient  WiresharkClient
	actionc          chan func()
	quitc            chan chan struct{}
	closedC          chan struct{}
	Logs             <-chan string
	logs             chan<- string
	err              chan error
	state            state
	hijackedResponse *types.HijackedResponse
	snifferId        string
	fifoPath         string
}

func NewSnifferActor(dockerInterface DockerClient, sshClient SshClient, wiresharkClient WiresharkClient) *Actor {
	logs := make(chan string)
	sa := Actor{
		dockerIface:      dockerInterface,
		sshClient:        sshClient,
		wiresharkClient:  wiresharkClient,
		actionc:          make(chan func()),
		quitc:            make(chan chan struct{}),
		logs:             logs,
		Logs:             logs,
		err:              make(chan error),
		state:            created,
		hijackedResponse: nil,
		snifferId:        "",
		fifoPath:         "",
	}

	go sa.handleMessages()

	return &sa
}

func (sa *Actor) handleMessages() {
	for {
		select {
		case f := <-sa.actionc:
			sa.guardedExec(f)
		case ch := <-sa.quitc:
			sa.guardedExec(sa.close)
			close(ch)
		}
	}
}

func (sa *Actor) guardedExec(f func()) {
	if sa.state == closed {
		sa.err <- errors.New("closed")
		return
	}

	f()
}

func (sa *Actor) close() {
	log.Println("Cleanup in progress...")
	err := sa.sshClient.Close()
	if err != nil {
		log.Fatal(err)
	}
	if sa.state == sniffing {
		sa.hijackedResponse.Close()
		log.Println("Hijacked response closed")

		if sa.snifferId != "" {
			dur := 30 * time.Second
			err := sa.dockerIface.ContainerStop(context.Background(), sa.snifferId, dur)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println("Sniffer container stopped")

		if sa.fifoPath != "" {
			err := os.Remove(sa.fifoPath)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println("Fifo " + sa.fifoPath + " removed")
	}
	sa.state = closed
	log.Println("Cleanup complete!")
	close(sa.closedC)
}

func (sa *Actor) ensureState(s state) error {
	if sa.state != s {
		return errors.New("wrong state")
	} else {
		return nil
	}
}
