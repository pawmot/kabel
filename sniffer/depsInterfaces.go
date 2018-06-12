package sniffer

//go:generate moq -out depsInterfaces_test.go . DockerClient SshClient

type SshTunnelLocalPort int

type SshClient interface {
	CreateTunnel(remoteSpec string) (local SshTunnelLocalPort, err error)
	Close() error
}

type DockerClient interface {
	Connect(endpoint string) error
	PullImage(imageName string) error
	GetContainers() ([]Container, error)
	GetNetworkInterfaces(containerId string) ([]NetworkInterface, error)
}