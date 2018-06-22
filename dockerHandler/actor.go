package dockerHandler

import (
	"github.com/pawmot/kabel/sniffer"
	"github.com/docker/docker/api/types"
	"time"
	"context"
	"github.com/docker/docker/client"
	"errors"
	"io/ioutil"
	"bytes"
	"github.com/docker/docker/pkg/stdcopy"
	"strings"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

var (
	notConnectedErr = errors.New("not connected")
)

type Actor struct {
	api     *client.Client
	actionC chan func()
}

func NewDockerHandler() *Actor {
	actionC := make(chan func())
	actor := Actor{
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

func (a *Actor) Connect(endpoint string) error {
	errC := make(chan error)
	a.actionC <- func() {
		if a.api == nil {
			api, err := client.NewClientWithOpts(client.WithHost(endpoint))
			if err != nil {
				errC <- err
			} else {
				a.api = api
				close(errC)
			}
		} else {
			errC <- errors.New("already connected")
		}
	}
	return <-errC
}

func (a *Actor) Close() error {
	errC := make(chan error)
	a.actionC <- func() {
		if a.api != nil {
			err := a.api.Close()
			if err != nil {
				errC <- err
			}
		}

		close(errC)
	}
	return <-errC
}

func (a *Actor) PullTcpDumpImage() error {
	errC := make(chan error)
	a.actionC <- func() {
		if a.api == nil {
			errC <- notConnectedErr
			return
		}

		imageName := "pawmot/tcpdump"
		ctx := context.Background()
		resp, err := a.api.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			errC <- err
			return
		}
		defer resp.Close()
		// TODO report the progress (by layer?)
		_, err = ioutil.ReadAll(resp)
		if err != nil {
			errC <- err
			return
		}

		close(errC)
	}
	return <-errC
}

func (a *Actor) GetContainers() ([]sniffer.Container, error) {
	containersC := make(chan []sniffer.Container)
	errC := make(chan error)
	a.actionC <- func() {
		if a.api == nil {
			errC <- notConnectedErr
			return
		}

		ctx := context.Background()
		containers, err := a.api.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			errC <- err
			return
		}

		var cs []sniffer.Container
		for _, c := range containers {
			container := sniffer.Container{
				Name: c.Names[0],
				Id:   c.ID,
			}
			cs = append(cs, container)
		}
		containersC <- cs
	}
	select {
	case cs := <-containersC:
		return cs, nil
	case err := <-errC:
		return nil, err
	}
}

func (a *Actor) GetNetworkInterfaces(containerId string) ([]sniffer.NetworkInterface, error) {
	netIfacesC := make(chan []sniffer.NetworkInterface)
	errC := make(chan error)
	a.actionC <- func() {
		if a.api == nil {
			errC <- notConnectedErr
			return
		}

		ctx := context.Background()
		// TODO provide running stats of rx/tx
		exec, err := a.api.ContainerExecCreate(ctx, containerId, types.ExecConfig{
			AttachStderr: true,
			AttachStdout: true,
			Tty:          true,
			Cmd:          []string{"ls", "--color=none", "/sys/class/net"},
		})
		if err != nil {
			errC <- fmt.Errorf("couldn't create Exec: %v", err)
			return
		}
		bufout := bytes.NewBufferString("")
		buferr := bytes.NewBufferString("")
		resp, err := a.api.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{Detach: false, Tty: false})
		if err != nil {
			errC <- fmt.Errorf("couldn't start Exec: %v", err)
			return
		}
		defer resp.Close()
		stdcopy.StdCopy(bufout, buferr, resp.Reader)
		if buferr.Len() > 0 {
			errC <- fmt.Errorf("couldn't read container's interfaces: %s", buferr.String())
			return
		}

		ifacesStr := strings.Split(strings.Replace(bufout.String(), "  ", " ", -1), " ")
		var ifaces []sniffer.NetworkInterface

		for _, i := range ifacesStr {
			iface := sniffer.NetworkInterface{
				Name: strings.Trim(i, "\n"),
			}
			ifaces = append(ifaces, iface)
		}
		netIfacesC <- ifaces
	}

	select {
	case nis := <-netIfacesC:
		return nis, nil
	case err := <-errC:
		return nil, err
	}
}

func (a *Actor) CreateTcpDumpContainer(name string, containerIdToSniff, interfaceName string) (id string, err error) {
	idC := make(chan string)
	errC := make(chan error)
	a.actionC <- func() {
		if a.api == nil {
			errC <- notConnectedErr
			return
		}

		ctx := context.Background()
		tdContainer, err := a.api.ContainerCreate(ctx, &container.Config{
			Image: "pawmot/tcpdump",
			Env: []string{
				"IF=" + interfaceName,
			},
			AttachStdout: true,
			AttachStderr: true,
		}, &container.HostConfig{
			NetworkMode:   container.NetworkMode("container:" + containerIdToSniff),
			AutoRemove:    true,
			DNS:           []string{},
			DNSOptions:    []string{},
			DNSSearch:     []string{},
			RestartPolicy: container.RestartPolicy{Name: "no", MaximumRetryCount: 0},
		}, &network.NetworkingConfig{

		}, name)

		if err != nil {
			errC <- err
		} else {
			idC <- tdContainer.ID
		}
	}

	select {
	case id := <-idC:
		return id, nil
	case err := <-errC:
		return "", err
	}
}

func (a *Actor) ContainerAttach(ctx context.Context, containerId string) (types.HijackedResponse, error) {
	respC := make(chan types.HijackedResponse)
	errC := make(chan error)

	a.actionC <- func() {
		if a.api == nil {
			errC <- notConnectedErr
			return
		}

		resp, err := a.api.ContainerAttach(ctx, containerId, types.ContainerAttachOptions{
			//Logs:   true,
			Stdout: true,
			Stderr: true,
			Stream: true,
		})

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
		return types.HijackedResponse{}, err
	}
}

func (a *Actor) ContainerStart(ctx context.Context, containerId string) error {
	errC := make(chan error)

	a.actionC <- func() {
		if a.api == nil {
			errC <- notConnectedErr
			return
		}

		err := a.api.ContainerStart(ctx, containerId, types.ContainerStartOptions{})

		if err != nil {
			errC <- err
		} else {
			close(errC)
		}
	}

	return <-errC
}

func (a *Actor) ContainerStop(ctx context.Context, snifferId string, duration time.Duration) error {
	errC := make(chan error)

	a.actionC <- func() {
		if a.api == nil {
			errC <- notConnectedErr
			return
		}

		err := a.api.ContainerStop(ctx, snifferId, &duration)

		if err != nil {
			errC <- err
		} else {
			close(errC)
		}
	}

	return <-errC
}
