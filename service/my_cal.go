package service

import (
	"fmt"
	"github.com/racecarparts/dashster/model"
	"strings"
)

func MyCal() model.MyCal {
	calNames := model.AppConfig.MyCalendar.ExcludedCalendars
	excludedCals := ""
	comma := ","
	for i, calName := range calNames {
		if i+1 == len(calNames) {
			comma = ""
		}
		excludedCals += fmt.Sprintf("\"%s\"%s", calName, comma)
	}
	icalBuddyCmdBytes := runcmd("which icalBuddy", true)
	if len(icalBuddyCmdBytes) < 0 {
		return model.MyCal{}
	}
	icalBuddyCmd := strings.TrimSuffix(string(icalBuddyCmdBytes), "\n")
	//cmd := fmt.Sprintf("/usr/local/bin/icalBuddy -sd -eep \"url\",\"notes\",\"attendees\",\"location\" -ec %s -n -nrd eventsToday", excludedCals)
	//cmd2 := "/usr/local/bin/icalBuddy -sd -eep \"url\",\"notes\",\"attendees\",\"location\" -ec \"Radio Nets\" -n -nrd eventsToday"
	cmd := icalBuddyCmd + " -sd -eep \"url\",\"notes\",\"attendees\",\"location\" -ec \"United States holidays\",\"Holidays in United States\",\"Radio Nets\" -n -nrd eventsToday"

	agenda := string(runcmd(cmd, true))

	return model.MyCal{Agenda: agenda}
}
