package sniffer

import (
	"github.com/docker/docker/api/types"
	"context"
	"time"
)

//go:generate moq -out depsInterfaces_gen_test.go . DockerClient SshClient WiresharkClient

type SshTunnelLocalPort int

type SshClient interface {
	CreateTunnel(remoteSpec string, dockerEndpoint string) (local SshTunnelLocalPort, err error)
	Close() error
}

type DockerClient interface {
	Connect(endpoint string) error
	Close() error
	PullTcpDumpImage() error
	GetContainers() ([]Container, error)
	GetNetworkInterfaces(containerId string) ([]NetworkInterface, error)
	CreateTcpDumpContainer(name string, containerIdToSniff, interfaceName string) (id string, err error)
	ContainerAttach(ctx context.Context, containerId string) (types.HijackedResponse, error)
	ContainerStart(ctx context.Context, containerId string) error
	ContainerStop(ctx context.Context, snifferId string, duration time.Duration) error
}

type WiresharkClient interface {
	Open(fifoPath string, closedC chan<- struct{}) error
	Close() error
}
