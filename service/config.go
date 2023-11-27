package service

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/racecarparts/dashster/model"
)

var configFilename = ".dashster_config.json"

func ReadOrCreateConfig() error {
	userFolder, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	confFilePath := userFolder + string(os.PathSeparator) + configFilename

	confFile, err := os.OpenFile(confFilePath, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer confFile.Close()

	confData, err := io.ReadAll(confFile)

	if len(confData) == 0 {
		model.AppConfig = &model.Config{
			WorldClock: model.WorldClock{
				HighlightCurrentTZ: false,
				TimeZones: []model.TimeZone{
					{
						TimeZone: "",
						Group:    0,
						People:   []string{},
					},
				},
			},
			Calendar: model.Calendar{},
			Weather: model.WeatherConfig{
				Locations: []model.WeatherLocation{
					{
						Location:    "",
						DisplayName: "",
					},
				},
				WeatherDataUrl: "",
			},
			MyCalendar: model.MyCalendar{
				ExcludedCalendars: []string{""},
			},
			Docker: model.Docker{},
			GithubPulls: model.GithubPulls{
				Enabled: false,
				Organizations: []model.GithubOrg{
					{
						Name:         "",
						TeamNameSlug: "",
						Username:     "",
						AccessKey:    "",
						Team: model.GithubTeam{
							Id:   0,
							Name: "",
							Slug: "",
							Members: []model.GithubUser{
								{
									Login: "",
								},
							},
						},
					},
				},
			},
			Gitlab: model.Gitlab{
				Enabled: false,
				Organizations: []model.GitlabOrg{
					{
						Name:             "",
						BaseUrl:          "",
						PrivateToken:     "",
						Username:         "",
						GroupNames:       []string{},
						FilterMRsByGroup: false,
					},
				},
			},
		}

		confData, err := json.MarshalIndent(model.AppConfig, "", "  ")
		if err != nil {
			return err
		}
		err = os.WriteFile(confFilePath, confData, 0644)
		if err != nil {
			return err
		}
		return errors.New("config file was not present, but has been created, finish filling it out at " + confFilePath + ".")
	}

	err = json.Unmarshal(confData, &model.AppConfig)
	if err != nil {
		return err
	}

	return nil
}
