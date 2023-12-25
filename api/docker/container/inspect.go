package container

import (
	"context"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
)

func (i *ContainerInformation) SetContainerInformation(ctx context.Context, cli *client.Client) error {
	info, err := cli.ContainerInspect(ctx, i.ID)
	if err != nil {
		return err
	}
	i.ContainerIP = info.NetworkSettings.IPAddress

	for port := range info.NetworkSettings.Ports {
		portStr := port.Port()
		portNumber := strings.Split(portStr, "/")[0]
		if port, err := strconv.ParseUint(portNumber, 10, 16); err == nil {
			i.ContainerPorts = append(i.ContainerPorts, uint16(port))
		}
	}

	for _, portBindings := range info.NetworkSettings.Ports {
		for _, binding := range portBindings {
			if hostPort, err := strconv.ParseUint(binding.HostPort, 10, 16); err == nil {
				i.HostPorts = append(i.HostPorts, uint16(hostPort))
			}
		}
	}
	return err
}
