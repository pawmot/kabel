// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package sniffer

import (
	"context"
	"github.com/docker/docker/api/types"
	"sync"
	"time"
)

var (
	lockDockerClientMockClose                  sync.RWMutex
	lockDockerClientMockConnect                sync.RWMutex
	lockDockerClientMockContainerAttach        sync.RWMutex
	lockDockerClientMockContainerStart         sync.RWMutex
	lockDockerClientMockContainerStop          sync.RWMutex
	lockDockerClientMockCreateTcpDumpContainer sync.RWMutex
	lockDockerClientMockGetContainers          sync.RWMutex
	lockDockerClientMockGetNetworkInterfaces   sync.RWMutex
	lockDockerClientMockPullTcpDumpImage       sync.RWMutex
)

// DockerClientMock is a mock implementation of DockerClient.
//
//     func TestSomethingThatUsesDockerClient(t *testing.T) {
//
//         // make and configure a mocked DockerClient
//         mockedDockerClient := &DockerClientMock{
//             CloseFunc: func() error {
// 	               panic("TODO: mock out the Close method")
//             },
//             ConnectFunc: func(endpoint string) error {
// 	               panic("TODO: mock out the Connect method")
//             },
//             ContainerAttachFunc: func(ctx context.Context, containerId string) (types.HijackedResponse, error) {
// 	               panic("TODO: mock out the ContainerAttach method")
//             },
//             ContainerStartFunc: func(ctx context.Context, containerId string) error {
// 	               panic("TODO: mock out the ContainerStart method")
//             },
//             ContainerStopFunc: func(ctx context.Context, snifferId string, duration time.Duration) error {
// 	               panic("TODO: mock out the ContainerStop method")
//             },
//             CreateTcpDumpContainerFunc: func(name string, containerIdToSniff string, interfaceName string) (string, error) {
// 	               panic("TODO: mock out the CreateTcpDumpContainer method")
//             },
//             GetContainersFunc: func() ([]Container, error) {
// 	               panic("TODO: mock out the GetContainers method")
//             },
//             GetNetworkInterfacesFunc: func(containerId string) ([]NetworkInterface, error) {
// 	               panic("TODO: mock out the GetNetworkInterfaces method")
//             },
//             PullTcpDumpImageFunc: func() error {
// 	               panic("TODO: mock out the PullTcpDumpImage method")
//             },
//         }
//
//         // TODO: use mockedDockerClient in code that requires DockerClient
//         //       and then make assertions.
//
//     }
type DockerClientMock struct {
	// CloseFunc mocks the Close method.
	CloseFunc func() error

	// ConnectFunc mocks the Connect method.
	ConnectFunc func(endpoint string) error

	// ContainerAttachFunc mocks the ContainerAttach method.
	ContainerAttachFunc func(ctx context.Context, containerId string) (types.HijackedResponse, error)

	// ContainerStartFunc mocks the ContainerStart method.
	ContainerStartFunc func(ctx context.Context, containerId string) error

	// ContainerStopFunc mocks the ContainerStop method.
	ContainerStopFunc func(ctx context.Context, snifferId string, duration time.Duration) error

	// CreateTcpDumpContainerFunc mocks the CreateTcpDumpContainer method.
	CreateTcpDumpContainerFunc func(name string, containerIdToSniff string, interfaceName string) (string, error)

	// GetContainersFunc mocks the GetContainers method.
	GetContainersFunc func() ([]Container, error)

	// GetNetworkInterfacesFunc mocks the GetNetworkInterfaces method.
	GetNetworkInterfacesFunc func(containerId string) ([]NetworkInterface, error)

	// PullTcpDumpImageFunc mocks the PullTcpDumpImage method.
	PullTcpDumpImageFunc func() error

	// calls tracks calls to the methods.
	calls struct {
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// Connect holds details about calls to the Connect method.
		Connect []struct {
			// Endpoint is the endpoint argument value.
			Endpoint string
		}
		// ContainerAttach holds details about calls to the ContainerAttach method.
		ContainerAttach []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ContainerId is the containerId argument value.
			ContainerId string
		}
		// ContainerStart holds details about calls to the ContainerStart method.
		ContainerStart []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ContainerId is the containerId argument value.
			ContainerId string
		}
		// ContainerStop holds details about calls to the ContainerStop method.
		ContainerStop []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// SnifferId is the snifferId argument value.
			SnifferId string
			// Duration is the duration argument value.
			Duration time.Duration
		}
		// CreateTcpDumpContainer holds details about calls to the CreateTcpDumpContainer method.
		CreateTcpDumpContainer []struct {
			// Name is the name argument value.
			Name string
			// ContainerIdToSniff is the containerIdToSniff argument value.
			ContainerIdToSniff string
			// InterfaceName is the interfaceName argument value.
			InterfaceName string
		}
		// GetContainers holds details about calls to the GetContainers method.
		GetContainers []struct {
		}
		// GetNetworkInterfaces holds details about calls to the GetNetworkInterfaces method.
		GetNetworkInterfaces []struct {
			// ContainerId is the containerId argument value.
			ContainerId string
		}
		// PullTcpDumpImage holds details about calls to the PullTcpDumpImage method.
		PullTcpDumpImage []struct {
		}
	}
}

// Close calls CloseFunc.
func (mock *DockerClientMock) Close() error {
	if mock.CloseFunc == nil {
		panic("moq: DockerClientMock.CloseFunc is nil but DockerClient.Close was just called")
	}
	callInfo := struct {
	}{}
	lockDockerClientMockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	lockDockerClientMockClose.Unlock()
	return mock.CloseFunc()
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//     len(mockedDockerClient.CloseCalls())
func (mock *DockerClientMock) CloseCalls() []struct {
} {
	var calls []struct {
	}
	lockDockerClientMockClose.RLock()
	calls = mock.calls.Close
	lockDockerClientMockClose.RUnlock()
	return calls
}

// Connect calls ConnectFunc.
func (mock *DockerClientMock) Connect(endpoint string) error {
	if mock.ConnectFunc == nil {
		panic("moq: DockerClientMock.ConnectFunc is nil but DockerClient.Connect was just called")
	}
	callInfo := struct {
		Endpoint string
	}{
		Endpoint: endpoint,
	}
	lockDockerClientMockConnect.Lock()
	mock.calls.Connect = append(mock.calls.Connect, callInfo)
	lockDockerClientMockConnect.Unlock()
	return mock.ConnectFunc(endpoint)
}

// ConnectCalls gets all the calls that were made to Connect.
// Check the length with:
//     len(mockedDockerClient.ConnectCalls())
func (mock *DockerClientMock) ConnectCalls() []struct {
	Endpoint string
} {
	var calls []struct {
		Endpoint string
	}
	lockDockerClientMockConnect.RLock()
	calls = mock.calls.Connect
	lockDockerClientMockConnect.RUnlock()
	return calls
}

// ContainerAttach calls ContainerAttachFunc.
func (mock *DockerClientMock) ContainerAttach(ctx context.Context, containerId string) (types.HijackedResponse, error) {
	if mock.ContainerAttachFunc == nil {
		panic("moq: DockerClientMock.ContainerAttachFunc is nil but DockerClient.ContainerAttach was just called")
	}
	callInfo := struct {
		Ctx         context.Context
		ContainerId string
	}{
		Ctx:         ctx,
		ContainerId: containerId,
	}
	lockDockerClientMockContainerAttach.Lock()
	mock.calls.ContainerAttach = append(mock.calls.ContainerAttach, callInfo)
	lockDockerClientMockContainerAttach.Unlock()
	return mock.ContainerAttachFunc(ctx, containerId)
}

// ContainerAttachCalls gets all the calls that were made to ContainerAttach.
// Check the length with:
//     len(mockedDockerClient.ContainerAttachCalls())
func (mock *DockerClientMock) ContainerAttachCalls() []struct {
	Ctx         context.Context
	ContainerId string
} {
	var calls []struct {
		Ctx         context.Context
		ContainerId string
	}
	lockDockerClientMockContainerAttach.RLock()
	calls = mock.calls.ContainerAttach
	lockDockerClientMockContainerAttach.RUnlock()
	return calls
}

// ContainerStart calls ContainerStartFunc.
func (mock *DockerClientMock) ContainerStart(ctx context.Context, containerId string) error {
	if mock.ContainerStartFunc == nil {
		panic("moq: DockerClientMock.ContainerStartFunc is nil but DockerClient.ContainerStart was just called")
	}
	callInfo := struct {
		Ctx         context.Context
		ContainerId string
	}{
		Ctx:         ctx,
		ContainerId: containerId,
	}
	lockDockerClientMockContainerStart.Lock()
	mock.calls.ContainerStart = append(mock.calls.ContainerStart, callInfo)
	lockDockerClientMockContainerStart.Unlock()
	return mock.ContainerStartFunc(ctx, containerId)
}

// ContainerStartCalls gets all the calls that were made to ContainerStart.
// Check the length with:
//     len(mockedDockerClient.ContainerStartCalls())
func (mock *DockerClientMock) ContainerStartCalls() []struct {
	Ctx         context.Context
	ContainerId string
} {
	var calls []struct {
		Ctx         context.Context
		ContainerId string
	}
	lockDockerClientMockContainerStart.RLock()
	calls = mock.calls.ContainerStart
	lockDockerClientMockContainerStart.RUnlock()
	return calls
}

// ContainerStop calls ContainerStopFunc.
func (mock *DockerClientMock) ContainerStop(ctx context.Context, snifferId string, duration time.Duration) error {
	if mock.ContainerStopFunc == nil {
		panic("moq: DockerClientMock.ContainerStopFunc is nil but DockerClient.ContainerStop was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		SnifferId string
		Duration  time.Duration
	}{
		Ctx:       ctx,
		SnifferId: snifferId,
		Duration:  duration,
	}
	lockDockerClientMockContainerStop.Lock()
	mock.calls.ContainerStop = append(mock.calls.ContainerStop, callInfo)
	lockDockerClientMockContainerStop.Unlock()
	return mock.ContainerStopFunc(ctx, snifferId, duration)
}

// ContainerStopCalls gets all the calls that were made to ContainerStop.
// Check the length with:
//     len(mockedDockerClient.ContainerStopCalls())
func (mock *DockerClientMock) ContainerStopCalls() []struct {
	Ctx       context.Context
	SnifferId string
	Duration  time.Duration
} {
	var calls []struct {
		Ctx       context.Context
		SnifferId string
		Duration  time.Duration
	}
	lockDockerClientMockContainerStop.RLock()
	calls = mock.calls.ContainerStop
	lockDockerClientMockContainerStop.RUnlock()
	return calls
}

// CreateTcpDumpContainer calls CreateTcpDumpContainerFunc.
func (mock *DockerClientMock) CreateTcpDumpContainer(name string, containerIdToSniff string, interfaceName string) (string, error) {
	if mock.CreateTcpDumpContainerFunc == nil {
		panic("moq: DockerClientMock.CreateTcpDumpContainerFunc is nil but DockerClient.CreateTcpDumpContainer was just called")
	}
	callInfo := struct {
		Name               string
		ContainerIdToSniff string
		InterfaceName      string
	}{
		Name:               name,
		ContainerIdToSniff: containerIdToSniff,
		InterfaceName:      interfaceName,
	}
	lockDockerClientMockCreateTcpDumpContainer.Lock()
	mock.calls.CreateTcpDumpContainer = append(mock.calls.CreateTcpDumpContainer, callInfo)
	lockDockerClientMockCreateTcpDumpContainer.Unlock()
	return mock.CreateTcpDumpContainerFunc(name, containerIdToSniff, interfaceName)
}

// CreateTcpDumpContainerCalls gets all the calls that were made to CreateTcpDumpContainer.
// Check the length with:
//     len(mockedDockerClient.CreateTcpDumpContainerCalls())
func (mock *DockerClientMock) CreateTcpDumpContainerCalls() []struct {
	Name               string
	ContainerIdToSniff string
	InterfaceName      string
} {
	var calls []struct {
		Name               string
		ContainerIdToSniff string
		InterfaceName      string
	}
	lockDockerClientMockCreateTcpDumpContainer.RLock()
	calls = mock.calls.CreateTcpDumpContainer
	lockDockerClientMockCreateTcpDumpContainer.RUnlock()
	return calls
}

// GetContainers calls GetContainersFunc.
func (mock *DockerClientMock) GetContainers() ([]Container, error) {
	if mock.GetContainersFunc == nil {
		panic("moq: DockerClientMock.GetContainersFunc is nil but DockerClient.GetContainers was just called")
	}
	callInfo := struct {
	}{}
	lockDockerClientMockGetContainers.Lock()
	mock.calls.GetContainers = append(mock.calls.GetContainers, callInfo)
	lockDockerClientMockGetContainers.Unlock()
	return mock.GetContainersFunc()
}

// GetContainersCalls gets all the calls that were made to GetContainers.
// Check the length with:
//     len(mockedDockerClient.GetContainersCalls())
func (mock *DockerClientMock) GetContainersCalls() []struct {
} {
	var calls []struct {
	}
	lockDockerClientMockGetContainers.RLock()
	calls = mock.calls.GetContainers
	lockDockerClientMockGetContainers.RUnlock()
	return calls
}

// GetNetworkInterfaces calls GetNetworkInterfacesFunc.
func (mock *DockerClientMock) GetNetworkInterfaces(containerId string) ([]NetworkInterface, error) {
	if mock.GetNetworkInterfacesFunc == nil {
		panic("moq: DockerClientMock.GetNetworkInterfacesFunc is nil but DockerClient.GetNetworkInterfaces was just called")
	}
	callInfo := struct {
		ContainerId string
	}{
		ContainerId: containerId,
	}
	lockDockerClientMockGetNetworkInterfaces.Lock()
	mock.calls.GetNetworkInterfaces = append(mock.calls.GetNetworkInterfaces, callInfo)
	lockDockerClientMockGetNetworkInterfaces.Unlock()
	return mock.GetNetworkInterfacesFunc(containerId)
}

// GetNetworkInterfacesCalls gets all the calls that were made to GetNetworkInterfaces.
// Check the length with:
//     len(mockedDockerClient.GetNetworkInterfacesCalls())
func (mock *DockerClientMock) GetNetworkInterfacesCalls() []struct {
	ContainerId string
} {
	var calls []struct {
		ContainerId string
	}
	lockDockerClientMockGetNetworkInterfaces.RLock()
	calls = mock.calls.GetNetworkInterfaces
	lockDockerClientMockGetNetworkInterfaces.RUnlock()
	return calls
}

// PullTcpDumpImage calls PullTcpDumpImageFunc.
func (mock *DockerClientMock) PullTcpDumpImage() error {
	if mock.PullTcpDumpImageFunc == nil {
		panic("moq: DockerClientMock.PullTcpDumpImageFunc is nil but DockerClient.PullTcpDumpImage was just called")
	}
	callInfo := struct {
	}{}
	lockDockerClientMockPullTcpDumpImage.Lock()
	mock.calls.PullTcpDumpImage = append(mock.calls.PullTcpDumpImage, callInfo)
	lockDockerClientMockPullTcpDumpImage.Unlock()
	return mock.PullTcpDumpImageFunc()
}

// PullTcpDumpImageCalls gets all the calls that were made to PullTcpDumpImage.
// Check the length with:
//     len(mockedDockerClient.PullTcpDumpImageCalls())
func (mock *DockerClientMock) PullTcpDumpImageCalls() []struct {
} {
	var calls []struct {
	}
	lockDockerClientMockPullTcpDumpImage.RLock()
	calls = mock.calls.PullTcpDumpImage
	lockDockerClientMockPullTcpDumpImage.RUnlock()
	return calls
}

var (
	lockSshClientMockClose        sync.RWMutex
	lockSshClientMockCreateTunnel sync.RWMutex
)

// SshClientMock is a mock implementation of SshClient.
//
//     func TestSomethingThatUsesSshClient(t *testing.T) {
//
//         // make and configure a mocked SshClient
//         mockedSshClient := &SshClientMock{
//             CloseFunc: func() error {
// 	               panic("TODO: mock out the Close method")
//             },
//             CreateTunnelFunc: func(remoteSpec string) (SshTunnelLocalPort, error) {
// 	               panic("TODO: mock out the CreateTunnel method")
//             },
//         }
//
//         // TODO: use mockedSshClient in code that requires SshClient
//         //       and then make assertions.
//
//     }
type SshClientMock struct {
	// CloseFunc mocks the Close method.
	CloseFunc func() error

	// CreateTunnelFunc mocks the CreateTunnel method.
	CreateTunnelFunc func(remoteSpec string) (SshTunnelLocalPort, error)

	// calls tracks calls to the methods.
	calls struct {
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// CreateTunnel holds details about calls to the CreateTunnel method.
		CreateTunnel []struct {
			// RemoteSpec is the remoteSpec argument value.
			RemoteSpec string
		}
	}
}

// Close calls CloseFunc.
func (mock *SshClientMock) Close() error {
	if mock.CloseFunc == nil {
		panic("moq: SshClientMock.CloseFunc is nil but SshClient.Close was just called")
	}
	callInfo := struct {
	}{}
	lockSshClientMockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	lockSshClientMockClose.Unlock()
	return mock.CloseFunc()
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//     len(mockedSshClient.CloseCalls())
func (mock *SshClientMock) CloseCalls() []struct {
} {
	var calls []struct {
	}
	lockSshClientMockClose.RLock()
	calls = mock.calls.Close
	lockSshClientMockClose.RUnlock()
	return calls
}

// CreateTunnel calls CreateTunnelFunc.
func (mock *SshClientMock) CreateTunnel(remoteSpec string) (SshTunnelLocalPort, error) {
	if mock.CreateTunnelFunc == nil {
		panic("moq: SshClientMock.CreateTunnelFunc is nil but SshClient.CreateTunnel was just called")
	}
	callInfo := struct {
		RemoteSpec string
	}{
		RemoteSpec: remoteSpec,
	}
	lockSshClientMockCreateTunnel.Lock()
	mock.calls.CreateTunnel = append(mock.calls.CreateTunnel, callInfo)
	lockSshClientMockCreateTunnel.Unlock()
	return mock.CreateTunnelFunc(remoteSpec)
}

// CreateTunnelCalls gets all the calls that were made to CreateTunnel.
// Check the length with:
//     len(mockedSshClient.CreateTunnelCalls())
func (mock *SshClientMock) CreateTunnelCalls() []struct {
	RemoteSpec string
} {
	var calls []struct {
		RemoteSpec string
	}
	lockSshClientMockCreateTunnel.RLock()
	calls = mock.calls.CreateTunnel
	lockSshClientMockCreateTunnel.RUnlock()
	return calls
}

var (
	lockWiresharkClientMockClose sync.RWMutex
	lockWiresharkClientMockOpen  sync.RWMutex
)

// WiresharkClientMock is a mock implementation of WiresharkClient.
//
//     func TestSomethingThatUsesWiresharkClient(t *testing.T) {
//
//         // make and configure a mocked WiresharkClient
//         mockedWiresharkClient := &WiresharkClientMock{
//             CloseFunc: func() error {
// 	               panic("TODO: mock out the Close method")
//             },
//             OpenFunc: func(fifoPath string, closedC chan<- struct{}) error {
// 	               panic("TODO: mock out the Open method")
//             },
//         }
//
//         // TODO: use mockedWiresharkClient in code that requires WiresharkClient
//         //       and then make assertions.
//
//     }
type WiresharkClientMock struct {
	// CloseFunc mocks the Close method.
	CloseFunc func() error

	// OpenFunc mocks the Open method.
	OpenFunc func(fifoPath string, closedC chan<- struct{}) error

	// calls tracks calls to the methods.
	calls struct {
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// Open holds details about calls to the Open method.
		Open []struct {
			// FifoPath is the fifoPath argument value.
			FifoPath string
			// ClosedC is the closedC argument value.
			ClosedC chan<- struct{}
		}
	}
}

// Close calls CloseFunc.
func (mock *WiresharkClientMock) Close() error {
	if mock.CloseFunc == nil {
		panic("moq: WiresharkClientMock.CloseFunc is nil but WiresharkClient.Close was just called")
	}
	callInfo := struct {
	}{}
	lockWiresharkClientMockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	lockWiresharkClientMockClose.Unlock()
	return mock.CloseFunc()
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//     len(mockedWiresharkClient.CloseCalls())
func (mock *WiresharkClientMock) CloseCalls() []struct {
} {
	var calls []struct {
	}
	lockWiresharkClientMockClose.RLock()
	calls = mock.calls.Close
	lockWiresharkClientMockClose.RUnlock()
	return calls
}

// Open calls OpenFunc.
func (mock *WiresharkClientMock) Open(fifoPath string, closedC chan<- struct{}) error {
	if mock.OpenFunc == nil {
		panic("moq: WiresharkClientMock.OpenFunc is nil but WiresharkClient.Open was just called")
	}
	callInfo := struct {
		FifoPath string
		ClosedC  chan<- struct{}
	}{
		FifoPath: fifoPath,
		ClosedC:  closedC,
	}
	lockWiresharkClientMockOpen.Lock()
	mock.calls.Open = append(mock.calls.Open, callInfo)
	lockWiresharkClientMockOpen.Unlock()
	return mock.OpenFunc(fifoPath, closedC)
}

// OpenCalls gets all the calls that were made to Open.
// Check the length with:
//     len(mockedWiresharkClient.OpenCalls())
func (mock *WiresharkClientMock) OpenCalls() []struct {
	FifoPath string
	ClosedC  chan<- struct{}
} {
	var calls []struct {
		FifoPath string
		ClosedC  chan<- struct{}
	}
	lockWiresharkClientMockOpen.RLock()
	calls = mock.calls.Open
	lockWiresharkClientMockOpen.RUnlock()
	return calls
}