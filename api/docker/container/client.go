package container

import "github.com/docker/docker/client"

const (
	dockerClientVersion = "1.42"
)

func CreateDockerClient() (cli *client.Client, err error) {
	cli, err = client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion(dockerClientVersion),
	)
	return
}
