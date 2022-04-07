package service

import "github.com/racecarparts/dashster/model"

func Docker() model.DockerStat {
	cmd := "docker ps --format \"table {{.Names}}\\t{{.Status}}\" | cut -c-$(tput cols)"
	dockerStat := string(runcmd(cmd, true))

	return model.DockerStat{Stat: dockerStat}
}
