package handler

import (
	"log/slog"

	docker "github.com/malsuke/seccube-back/api/docker/container"
	"github.com/malsuke/seccube-back/utils"

	"github.com/docker/docker/api/types/container"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var (
	ssh = []*docker.ContainerService{
		docker.NewContainerWithConfig(
			&container.Config{
				Image: "password-crack-attack:latest",
				Tty:   true,
			},
			&container.HostConfig{
				AutoRemove: true,
				// PortBindings: nat.PortMap{
				// 	"22/tcp": []nat.PortBinding{
				// 		{
				// 			HostPort: "0",
				// 		},
				// 	},
				// },
				Resources: container.Resources{
					Memory: 1024 * 1024 * 1024,
				},
			},
			nil,
			nil,
		),
		docker.NewContainerWithConfig(
			&container.Config{
				Image: "password-crack-defense:latest",
				Tty:   true,
			},
			&container.HostConfig{
				AutoRemove: true,
				// PortBindings: nat.PortMap{
				// 	"22/tcp": []nat.PortBinding{
				// 		{
				// 			HostPort: "0",
				// 		},
				// 	},
				// },
				Resources: container.Resources{
					Memory: 1024 * 1024 * 1024,
				},
			},
			nil,
			nil,
		),
	}
	ContainerList = map[string][]*docker.ContainerService{
		"sshBrute": ssh,
	}
)

func Create(c echo.Context) error {
	tag := c.Param("tag")
	slog.Info(tag)
	ctx := c.Request().Context()
	cli, err := docker.CreateDockerClient()
	if err != nil {
		return err
	}

	var ids []map[string]string

	nid := utils.GenerateUUID()
	nid, err = docker.CreateNetwork(ctx, cli, nid)
	if err != nil {
		return err
	}
	log.Debug().Str("network", nid).Msg("network created")

	for _, container := range ContainerList[tag] {
		if tag == "sqli" {
			container.SetNetworkEndpointConfigWithAlias(nid)
		} else {
			container.SetNetworkEndpointConfig(nid)
		}
		container.SetNetworkEndpointConfig(nid)
		log.Debug().Str("network", nid).Msg("network attached")
		id, err := container.CreateContainer(ctx, cli)
		if err != nil {
			return err
		}
		ids = append(ids, map[string]string{
			"id": *id,
		})
		c.Logger().Debug(id)
	}
	return c.JSON(200, ids)
}
