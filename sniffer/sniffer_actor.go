package sniffer

import (
	"errors"
)

type state int

const (
	created state = iota
	connected
	closed
)

type Actor struct {
	dockerIface DockerClient
	sshClient   SshClient
	actionc     chan func()
	quitc       chan struct{}
	Logs        <-chan string
	logs        chan<- string
	err         chan error
	state       state
}

func NewSnifferActor(dockerInterface DockerClient, sshClient SshClient) *Actor {
	logs := make(chan string)
	sa := Actor{
		dockerIface: dockerInterface,
		sshClient:   sshClient,
		actionc:     make(chan func()),
		quitc:       make(chan struct{}),
		logs:        logs,
		Logs:        logs,
		err:         make(chan error),
		state:       created,
	}

	go sa.handleMessages()

	return &sa
}

func (sa *Actor) handleMessages() {
	for {
		select {
		case f := <-sa.actionc:
			sa.guardedExec(f)
		case <-sa.quitc:
			sa.guardedExec(sa.close)
			sa.quitc <- struct{}{}
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
	sa.state = closed
	sa.sshClient.Close()
}

func (sa *Actor) ensureState(s state) error {
	if sa.state != s {
		return errors.New("Wrong state")
	} else {
		return nil
	}
}
