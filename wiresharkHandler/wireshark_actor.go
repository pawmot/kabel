package wiresharkHandler

import (
	"log"
	"os/exec"
)

type state int

const (
	ready   state = iota
	running
)

type Actor struct {
	state   state
	actionC chan func()
}

func NewWiresharkClient() *Actor {
	actionC := make(chan func())
	actor := Actor{
		state:   ready,
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

func (a *Actor) Open(fifoPath string, closedC chan<- struct{}) error {
	errC := make(chan error)
	a.actionC <- func() {
		if a.state == ready {
			log.Println("Running WireShark on '" + fifoPath + "'!")
			cmd := exec.Command("wireshark", "-k", "-i", fifoPath)
			err := cmd.Start()
			if err != nil {
				errC <- err
			} else {
				go func() {
					cmd.Wait()
					// TODO check the exit status
					a.Close()
					close(closedC)
				}()
				a.state = running
				close(errC)
			}
		}
	}
	return <-errC
}

func (a *Actor) Close() error {
	errC := make(chan error)
	a.actionC <- func() {
		if a.state == running {
			// TODO actually kill wireshark if it's still running
			a.state = ready
		}

		close(errC)
	}
	return <-errC
}
