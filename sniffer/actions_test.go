package sniffer

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"errors"
)

func Test_ActorShouldReturnErrorWhenClosedSecondTime(t *testing.T) {
	// given
	dc := DockerClientMock{}
	sc := SshClientMock{
		CloseFunc: func() error {
			return nil
		},
	}
	sa := NewSnifferActor(&dc, &sc)

	// when
	err := sa.Close()

	// then
	assert.Nil(t, err)

	// when
	err = sa.Close()

	// then
	assert.Error(t, err)
}

func Test_ConnectShouldForwardErrorFromDocker(t *testing.T) {
	// given
	expectedErr := errors.New("sample expectedErr")
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return expectedErr
		},
	}
	sc := SshClientMock{
		CloseFunc: func() error {
			return nil
		},
	}
	sa := NewSnifferActor(&dc, &sc)

	// when
	resp, err := sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))

	// then
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, Error, resp)
}

func Test_ConnectShouldCloseSshWhenDockerErrorOccurs(t *testing.T) {
	//given
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return errors.New("Fek")
		},
	}
	sc := SshClientMock{
		CloseFunc: func() error {
			return nil
		},
	}
	sa := NewSnifferActor(&dc, &sc)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))

	// then
	assert.Equal(t, 1, len(sc.CloseCalls()))
}

func Test_ConnectShouldForwardErrorFromSsh(t *testing.T) {
	// given
	expectedErr := errors.New("sample expectedErr")
	dc := DockerClientMock{}
	sc := SshClientMock{
		CreateTunnelFunc: func(remoteSpec string) (SshTunnelLocalPort, error) {
			return -1, expectedErr
		},
	}
	sa := NewSnifferActor(&dc, &sc)

	// when
	resp, err := sa.Connect(TunneledConnectionRequest("unix:///someDockerEndpoint/lol.sock", "user@host.com"))

	// then
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, Error, resp)
}

func Test_DirectConnectionRequestShouldCallDockerInterface(t *testing.T) {
	// given
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			if endpoint == "unix:///someDockerEndpoint/lol.sock" {
				return nil
			} else {
				t.Logf("Wrong endpoint!\nExpected: %s\nGot: %s", "unix:///someDockerEndpoint/lol.sock", endpoint)
				t.Fail()
				return nil
			}
		},
	}
	sa := NewSnifferActor(&dc, nil)

	// when
	resp, err := sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))

	// then
	assert.Nil(t, err)
	assert.Equal(t, Connected, resp)
	assert.Equal(t, 1, len(dc.ConnectCalls()))
}

func Test_TunneledConnectionShouldOpenSshTunnel(t *testing.T) {
	// given
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
	}
	sc := SshClientMock{
		CreateTunnelFunc: func(remoteSpec string) (SshTunnelLocalPort, error) {
			if remoteSpec == "user@host.com" {
				return 12345, nil
			} else {
				t.Logf("Wrong remote spec!\nExpected: %s\nGot: %s", "user@host.com", remoteSpec)
				t.Fail()
				return -1, nil
			}
		},
	}
	sa := NewSnifferActor(&dc, &sc)

	// when
	resp, err := sa.Connect(TunneledConnectionRequest("unix:///someDockerEndpoint/lol.sock", "user@host.com"))

	// then
	assert.Nil(t, err)
	assert.Equal(t, Connected, resp)
	assert.Equal(t, 1, len(dc.ConnectCalls()))
	assert.Equal(t, 1, len(sc.CreateTunnelCalls()))
}

func Test_TunneledConnectionShouldUseLocalEndpoint(t *testing.T) {
	// given
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
	}
	sc := SshClientMock{
		CreateTunnelFunc: func(remoteSpec string) (SshTunnelLocalPort, error) {
			if remoteSpec == "user@host.com" {
				return 12345, nil
			} else {
				t.Logf("Wrong remote spec!\nExpected: %s\nGot: %s", "user@host.com", remoteSpec)
				t.Fail()
				return -1, nil
			}
		},
	}
	sa := NewSnifferActor(&dc, &sc)

	// when
	resp, err := sa.Connect(TunneledConnectionRequest("unix:///someDockerEndpoint/lol.sock", "user@host.com"))

	// then
	assert.Nil(t, err)
	assert.Equal(t, Connected, resp)
	assert.Equal(t, "http://localhost:12345", dc.ConnectCalls()[0].Endpoint)
}

func Test_PullImageShouldCheckIfConnected(t *testing.T) {
	// given
	dc := DockerClientMock{}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	err := sa.PullImage()

	// then
	assert.Error(t, err)
}

func Test_PullImageShouldPullTcpDumpImage(t *testing.T) {
	// given
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
		PullImageFunc: func(imageName string) error {
			return nil
		},
	}
	sa := NewSnifferActor(&dc, nil)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))
	err := sa.PullImage()

	// then
	assert.Nil(t, err)
	assert.Equal(t, "pawmot/tcpdump", dc.PullImageCalls()[0].ImageName)
}

func Test_PullImageShouldReportErrors(t *testing.T) {
	// given
	errExp := errors.New("expected")
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
		PullImageFunc: func(imageName string) error {
			return errExp
		},
	}
	sa := NewSnifferActor(&dc, nil)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))
	err := sa.PullImage()

	// then
	assert.Equal(t, errExp, err)
}

func Test_GetContainersShouldCheckIfConnected(t *testing.T) {
	// given
	dc := DockerClientMock{}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	_, err := sa.GetContainers()

	// then
	assert.Error(t, err)
}

func Test_GetContainersShouldCallDockerIface(t *testing.T) {
	// given
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
		GetContainersFunc: func() ([]Container, error) {
			return []Container{}, nil
		},
	}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))
	sa.GetContainers()

	// then
	assert.Equal(t, 1, len(dc.GetContainersCalls()))
}

func Test_GetContainersShouldReturnDockerIfaceResponse(t *testing.T) {
	// given
	exp := []Container{
		{id: "1", name: "3"},
	}
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
		GetContainersFunc:
		func() ([]Container, error) {
			return exp, nil
		},
	}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))
	resp, _ := sa.GetContainers()

	// then
	assert.Equal(t, exp, resp)
}

func Test_GetContainersShouldForwardDockerIfaceError(t *testing.T) {
	// given
	exp := errors.New("Fek")
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
		GetContainersFunc:
		func() ([]Container, error) {
			return nil, exp
		},
	}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))
	_, err := sa.GetContainers()

	// then
	assert.Equal(t, exp, err)
}

func Test_GetNetworkInterfacesShouldCheckIfConnected(t *testing.T) {
	// given
	dc := DockerClientMock{}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	_, err := sa.GetNetworkInterfaces("id")

	// then
	assert.Error(t, err)
}

func Test_GetNetworkInterfacesShouldCallDockerIface(t *testing.T) {
	// given
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
		GetNetworkInterfacesFunc:
		func(id string) ([]NetworkInterface, error) {
			return []NetworkInterface{}, nil
		},
	}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))
	sa.GetNetworkInterfaces("id")

	// then
	assert.Equal(t, 1, len(dc.GetNetworkInterfacesCalls()))
}

func Test_GetNetworkInterfacesShouldReturnDockerIfaceResponse(t *testing.T) {
	// given
	exp := []NetworkInterface{
		{name: "3"},
	}
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
		GetNetworkInterfacesFunc:
		func(id string) ([]NetworkInterface, error) {
			return exp, nil
		},
	}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))
	resp, _ := sa.GetNetworkInterfaces("id")

	// then
	assert.Equal(t, exp, resp)
}

func Test_GetNetworkInterfacesShouldForwardDockerIfaceError(t *testing.T) {
	// given
	exp := errors.New("Fek")
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
		GetNetworkInterfacesFunc:
		func(id string) ([]NetworkInterface, error) {
			return nil, exp
		},
	}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))
	_, err := sa.GetNetworkInterfaces("id")

	// then
	assert.Equal(t, exp, err)
}

func Test_CreateSnifferContainerShouldCheckIfConnected(t *testing.T) {
	// given
	dc := DockerClientMock{}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	_, err := sa.GetNetworkInterfaces("id")

	// then
	assert.Error(t, err)
}

func Test_CreateSnifferContainerShouldCallDockerIface(t *testing.T) {
	// given
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
		GetNetworkInterfacesFunc:
		func(id string) ([]NetworkInterface, error) {
			return []NetworkInterface{}, nil
		},
	}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))
	sa.GetNetworkInterfaces("id")

	// then
	assert.Equal(t, 1, len(dc.GetNetworkInterfacesCalls()))
}

func Test_CreateSnifferContainerShouldReturnDockerIfaceResponse(t *testing.T) {
	// given
	exp := []NetworkInterface{
		{name: "3"},
	}
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
		GetNetworkInterfacesFunc:
		func(id string) ([]NetworkInterface, error) {
			return exp, nil
		},
	}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))
	resp, _ := sa.GetNetworkInterfaces("id")

	// then
	assert.Equal(t, exp, resp)
}

func Test_CreateSnifferContainerShouldForwardDockerIfaceError(t *testing.T) {
	// given
	exp := errors.New("Fek")
	dc := DockerClientMock{
		ConnectFunc: func(endpoint string) error {
			return nil
		},
		GetNetworkInterfacesFunc:
		func(id string) ([]NetworkInterface, error) {
			return nil, exp
		},
	}
	sc := SshClientMock{}
	sa := NewSnifferActor(&dc, &sc)

	// when
	sa.Connect(DirectConnectionRequest("unix:///someDockerEndpoint/lol.sock"))
	_, err := sa.GetNetworkInterfaces("id")

	// then
	assert.Equal(t, exp, err)
}
