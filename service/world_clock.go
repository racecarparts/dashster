package service

import (
	"fmt"
	"time"

	"github.com/racecarparts/dashster/model"
)

const (
	offsetFmt  = "%5s"
	shortTZFmt = "%4s"
)

func WorldClock() []model.ClockTime {
	clockCfg := model.AppConfig.WorldClock
	clockTimes := make([]model.ClockTime, len(clockCfg.TimeZones))
	now := time.Now()
	myZone, _ := now.Zone()
	for i, tz := range clockCfg.TimeZones {
		clockTime := model.ClockTime{
			TZ:            tz.TimeZone,
			Group:         tz.Group,
			Time:          time.Time{},
			UtcOffset:     "+0000",
			IsCurrentZone: false,
			People:        tz.People,
		}
		loc, err := time.LoadLocation(tz.TimeZone)
		if err != nil {
			continue
		}

		nowTZ := now.In(loc)
		zone, _ := nowTZ.Zone()
		if myZone == zone {
			clockTime.IsCurrentZone = true
		}
		clockTime.Time = now
		clockTime.UtcOffset = fmt.Sprintf(offsetFmt, nowTZ.Format("-0700"))
		clockTime.ShortTZ = fmt.Sprintf(shortTZFmt, nowTZ.Format(zone))

		clockTimes[i] = clockTime
	}

	return clockTimes
}
