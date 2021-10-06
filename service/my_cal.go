package service

import "github.com/racecarparts/dashster/model"

func MyCal() model.MyCal {
    cmd := "/usr/local/bin/icalBuddy -sd -eep \"url\",\"notes\",\"attendees\",\"location\" -ec \"United States holidays\",\"Holidays in United States\" -n -nrd eventsToday"

    agenda := string(runcmd(cmd, true))

    return model.MyCal{Agenda: agenda}
}