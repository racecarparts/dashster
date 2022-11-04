package model

var AppConfig *Config

type Config struct {
	WorldClock  WorldClock    `json:"world_clock"`
	Calendar    Calendar      `json:"calendar"`
	Weather     WeatherConfig `json:"weather"`
	MyCalendar  MyCalendar    `json:"my_calendar"`
	Docker      Docker        `json:"docker"`
	GithubPulls GithubPulls   `json:"github_pulls"`
	Gitlab      Gitlab        `json:"gitlab"`
}

type WorldClock struct {
	HighlightCurrentTZ bool       `json:"highlight_current_tz"`
	TimeZones          []TimeZone `json:"time_zones"`
}

type TimeZone struct {
	TimeZone string `json:"time_zone"`
	Group    int    `json:"group"`
}

type Calendar struct{}

type WeatherConfig struct {
	Locations      []WeatherLocation `json:"locations"`
	WeatherDataUrl string            `json:"weatherDataUrl"`
}

type WeatherLocation struct {
	Location    string `json:"location"`
	DisplayName string `json:"display_name"`
}

type MyCalendar struct {
	ExcludedCalendars []string `json:"excluded_calendars"`
}

type Docker struct{}

type GithubPulls struct {
	Enabled       bool        `json:"enabled"`
	Organizations []GithubOrg `json:"organizations"`
}

type GithubOrg struct {
	Name         string     `json:"name"`
	TeamNameSlug string     `json:"team_name_slug"`
	Username     string     `json:"username"`
	AccessKey    string     `json:"access_key"`
	Team         GithubTeam `json:"team"`
}

type Gitlab struct {
	Enabled       bool        `json:"enabled"`
	Organizations []GitlabOrg `json:"organizations"`
}

type GitlabOrg struct {
	Name         string `json:"name"`
	BaseUrl      string `json:"base_url"`
	PrivateToken string `json:"private_token"`
	Username     string `json:"username"`
}
