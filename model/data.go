package model

import "time"

// Response types

type ClockTime struct {
	Group         int       `json:"group"`
	Time          time.Time `json:"time"`
	TZ            string    `json:"tz"`
	ShortTZ       string    `json:"short_tz"`
	UtcOffset     string    `json:"utc_offset"`
	IsCurrentZone bool      `json:"current_zone"`
}

type DisplayCalendar struct {
	Title string `json:"title"`
	Cal   string `json:"cal"`
}

type MainWeather struct {
	Temp float64 `json:"temp"`
}

type WeatherReport struct {
	Weather []Weather   `json:"weather"`
	Main    MainWeather `json:"main"`
	Name    string      `json:"name"`
}

type Weather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type SimpleWeather struct {
	Location      string `json:"location"`
	CurConditions string `json:"cur_conditions"`
	CurTemp       string `json:"cur_temp"`
}

type DockerStat struct {
	Stat string `json:"stat"`
}

type SimplePullRequests struct {
	Message      string              `json:"message"`
	MyPRs        []SimplePullRequest `json:"my_prs"`
	RequestedPRs []SimplePullRequest `json:"requested_prs"`
}

type SimplePullRequest struct {
	RepositoryName string       `json:"repository_name"`
	Number         int          `json:"number"`
	User           string   `json:"user"`
	Title          string       `json:"title"`
	Reviews        []PullReview `json:"reviews"`
	ReviewUrl      string       `json:"review_url"`
	IsDraft        bool         `json:"is_draft"`
}

type MyCal struct {
	Agenda string `json:"agenda"`
}

// Intermediate types for processing requests

type Repo struct {
	IsTeamRepo bool   `json:"-"` // does not exist from github, used for internal purposes
	Username   string `json:"-"` // does not exist from github, used for internal purposes
	OrgName    string `json:"-"` // does not exist from github, used for internal purposes
	AuthToken  string `json:"-"` // does not exist from github, used for internal purposes
	Id         int    `json:"id"`
	Name       string `json:"name"`
}

type PullRequest struct {
	Url                string       `json:"url"`
	HtmlUrl            string       `json:"html_url"`
	Number             int          `json:"number"`
	Title              string       `json:"title"`
	User               GithubUser   `json:"user"`
	Assignee           GithubUser   `json:"assignee"`
	Assignees          []GithubUser `json:"assignees"`
	RequestedReviewers []GithubUser `json:"requested_reviewers"`
	RequestedTeams     []GithubTeam `json:"requested_teams"`
	State              string       `json:"state"`
	Draft              bool         `json:"draft"`
}

type PullReview struct {
	User  GithubUser `json:"user"`
	State string     `json:"state"`
}

type GithubUser struct {
	Login string `json:"login"`
}

type GithubTeam struct {
	Id      int          `json:"id"`
	Name    string       `json:"name"`
	Slug    string       `json:"slug"`
	Members []GithubUser `json:"members"`
}
