package sniffer

import "fmt"

func (sa *Actor) Close() error {
	sa.quitc <- struct{}{}
	select {
	case <-sa.quitc:
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
		err = sa.dockerIface.PullImage("pawmot/tcpdump")
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

func (sa *Actor) CreateSnifferContainer(containerId string, interfaceName string) (error) {
	return nil
}