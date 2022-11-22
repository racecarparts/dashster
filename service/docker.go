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
	maxStatusLength := 0
	for i, container := range containers {
		name := fmt.Sprintf("no-name-%d", i)
		status := container.Status

		if len(container.Names) > 0 {
			name = container.Names[0]
		}

		stats = append(stats, model.DockerStat{
			Name:   name,
			Status: status,
			Ports:  container.Ports,
		})

		nameLen := len(name)
		if nameLen > maxNameLength {
			maxNameLength = nameLen
		}

		statusLen := len(status)
		if statusLen > maxStatusLength {
			maxStatusLength = statusLen
		}

	}

	if maxNameLength == 0 {
		maxNameLength = 5
	}

	if maxStatusLength == 0 {
		maxStatusLength = 6
	}

	nameHeader := padRight("Names", " ", maxNameLength)
	statusHeader := padRight("Status", " ", maxStatusLength)
	portsHeader := "Ports (host:container)"
	out := fmt.Sprintf("%s\t%s\t%s\n", nameHeader, statusHeader, portsHeader)
	for i, stat := range stats {
		portsStat := ""
		for i, port := range stat.Ports {
			portsStat += fmt.Sprintf("%d:%d", port.PublicPort, port.PrivatePort)
			if len(stat.Ports) > 1 && i < len(stat.Ports)-1 {
				portsStat += ", "
			}
		}

		name := padRight(stat.Name, " ", maxNameLength)
		status := padRight(stat.Status, " ", maxStatusLength)
		out += fmt.Sprintf("%s\t%s\t%s", name, status, portsStat)
		if i < len(stats)-1 {
			out += "\n"
		}
	}

	return model.DockerStatOut{Stat: out}
}
