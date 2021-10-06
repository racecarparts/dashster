package service

import (
    "github.com/racecarparts/dashster/model"
    "time"
)

func Calendar() model.DisplayCalendar {
    cmd := "/usr/local/bin/gcal .+ | sed -e '1,4d'"
    //cmd := "/usr/bin/cal -A2"
    titleTimeFormat := "Mon, 02 Jan 2006"
    calTitle := time.Now().Format(titleTimeFormat)
    cal := string(runcmd(cmd, true))

    return model.DisplayCalendar{
        Title: calTitle,
        Cal:   cal,
    }
}