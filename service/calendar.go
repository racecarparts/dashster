package service

import (
	"github.com/racecarparts/dashster/model"
	"strings"
	"time"
)

func Calendar() model.DisplayCalendar {
	gcalCmdBytes := runcmd("which gcal", true)
	if len(gcalCmdBytes) < 0 {
		return model.DisplayCalendar{}
	}
	gcalCmd := strings.TrimSuffix(string(gcalCmdBytes), "\n")

	cmd := gcalCmd + " . | sed -e '1,4d'"
	//cmd := "/usr/bin/cal -A2"
	titleTimeFormat := "Mon, 02 Jan 2006"
	calTitle := time.Now().Format(titleTimeFormat)
	cal := string(runcmd(cmd, true))

	return model.DisplayCalendar{
		Title: calTitle,
		Cal:   cal,
	}
}
