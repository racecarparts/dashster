package service

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/racecarparts/dashster/model"
)

func Docker() model.DockerStatOut {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	stats := []model.DockerStat{}
	maxNameLength := 0
	for i, container := range containers {
		name := fmt.Sprintf("no-name-%d", i)
		if len(container.Names) > 0 {
			name = container.Names[0]
		}

		stats = append(stats, model.DockerStat{
			Name: name,
			Status: container.Status,
		})

		nameLen := len(name)
		if nameLen > maxNameLength {
			maxNameLength = nameLen
		}
	}

	nameHeader := padRight("Names", " ", maxNameLength)
	out := fmt.Sprintf("%s\t%s\n", nameHeader, "Status")
	for i, stat := range stats {
		name := padRight(stat.Name, " ", maxNameLength)
		out += fmt.Sprintf("%s\t%s", name, stat.Status)
		if i < len(stats) - 1 {
			out += "\n"
		}
	}

	return model.DockerStatOut{Stat: out}
}
