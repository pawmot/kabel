package sniffer

import (
	"fmt"
	"strconv"
	"time"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"syscall"
	"github.com/docker/docker/pkg/stdcopy"
	"context"
)

func (sa *Actor) Close() error {
	ch := make(chan struct{})
	sa.quitc <- ch
	select {
	case <-ch:
		return nil

	case err := <-sa.err:
		return err
	}
}

func (sa *Actor) Connect(request ConnectionRequest) (ConnectionResponse, error) {
	respC := make(chan ConnectionResponse)
	errC := make(chan error)
	sa.actionc <- func() {
		spec := request.getConnSpec()
		var effectiveEndpoint string
		if spec.sshTunnelSpec != "" {
			// TODO handle password
			local, err := sa.sshClient.CreateTunnel(spec.sshTunnelSpec)
			if err != nil {
				errC <- err
				return
			}

			effectiveEndpoint = fmt.Sprintf("http://localhost:%d", local)
		} else {
			effectiveEndpoint = spec.dockerEndpoint
		}

		err := sa.dockerIface.Connect(effectiveEndpoint)
		if err != nil {
			sa.close()
			errC <- err
			return
		}

		sa.state = connected
		respC <- Connected
	}
	select {
	case resp := <-respC:
		return resp, nil
	case err := <-errC:
		return Error, err
	}
}

func (sa* Actor) PullImage() error {
	errC := make(chan error)
	sa.actionc <- func() {
		err := sa.ensureState(connected)
		if err != nil {
			errC <- err
			return
		}
		err = sa.dockerIface.PullTcpDumpImage()
		if err != nil {
			errC <- err
		} else {
			close(errC)
		}
	}
	return <-errC
}

func (sa* Actor) GetContainers() ([]Container, error) {
	respC := make(chan []Container)
	errC := make(chan error)
	sa.actionc <- func() {
		err := sa.ensureState(connected)
		if err != nil {
			errC <- err
			return
		}
		resp, err := sa.dockerIface.GetContainers()
		if err != nil {
			errC <- err
		} else {
			respC <- resp
		}
	}
	select {
	case resp := <-respC:
		return resp, nil
	case err := <-errC:
		return nil, err
	}
}

func (sa* Actor) GetNetworkInterfaces(containerId string) ([]NetworkInterface, error) {
	respC := make(chan []NetworkInterface)
	errC := make(chan error)
	sa.actionc <- func() {
		err := sa.ensureState(connected)
		if err != nil {
			errC <- err
			return
		}
		resp, err := sa.dockerIface.GetNetworkInterfaces(containerId)
		if err != nil {
			errC <- err
		} else {
			respC <- resp
		}
	}
	select {
	case resp := <-respC:
		return resp, nil
	case err := <-errC:
		return nil, err
	}
}

func (sa *Actor) Sniff(containerIdToSniff string, interfaceName string) (error) {
	errC := make(chan error)
	sa.actionc <- func() {
		err := sa.ensureState(connected)
		if err != nil {
			errC <- err
			return
		}

		name := "tcpdump-" + containerIdToSniff + "-" + interfaceName + "-" + strconv.FormatInt(time.Now().Unix(), 10)
		snifferId, err := sa.dockerIface.CreateTcpDumpContainer(name, containerIdToSniff, interfaceName)
		sa.snifferId = snifferId

		fifoPath := "/tmp/" + name
		unix.Mkfifo(fifoPath, 0666)

		sa.fifoPath = fifoPath
		closedC := make(chan struct{})
		sa.wiresharkClient.Open(fifoPath, closedC)

		go func() {
			<-closedC
			log.Println("Wireshark closed, tearing everything down!")
			sa.Close()
		}()

		fifo, err := os.OpenFile(fifoPath, syscall.O_WRONLY, 0600)

		if err != nil {
			errC <- err
			return
		}

		log.Println("WireShark connected!")

		log.Println("Attaching...")
		ctx := context.Background()
		resp, err := sa.dockerIface.ContainerAttach(ctx, snifferId)
		if err != nil {
			errC <- err
			return
		}
		log.Println("Attach goroutine finished...")

		sa.hijackedResponse = &resp
		go func() {
			stdcopy.StdCopy(fifo, os.Stderr, resp.Reader)
		}()

		err = sa.dockerIface.ContainerStart(ctx, snifferId)

		if err != nil {
			errC <- err
			return
		}

		log.Println("Continuing!")
		sa.state = sniffing
		close(errC)
	}
	return <-errC
}