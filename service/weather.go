package service

import (
    "encoding/json"
    "fmt"
    "github.com/racecarparts/dashster/model"
    "math"
    "net/url"
)

func Weather() []model.SimpleWeather {
    wCfg := model.AppConfig.Weather
    simpleWeatherReports := make([]model.SimpleWeather, len(wCfg.Locations))

    for i, location := range wCfg.Locations {
        resp, err := getRequest(fmt.Sprintf(wCfg.WeatherDataUrl, url.QueryEscape(location.Location)))
        if err != nil {
            simpleWeatherReports[i].Location = fmt.Sprintf("Problem getting weather report for '%s': %s", location, err.Error())
            continue
        }

        report := model.WeatherReport{}
        err = json.Unmarshal(resp, &report)
        if err != nil {
            simpleWeatherReports[i].Location = "Problem getting weather report: " + err.Error()
            continue
        }

        temp := fmt.Sprintf("%d°F", int64(math.Round(report.Main.Temp)))
        locName := report.Name
        if len(location.DisplayName) > 0 {
            locName = location.DisplayName
        }
        simpleWeather := model.SimpleWeather{
            Location:      locName,
            CurConditions: report.Weather[len(report.Weather)-1].Description,
            CurTemp:       temp,
        }
        simpleWeatherReports[i] = simpleWeather
    }
    return simpleWeatherReports
}