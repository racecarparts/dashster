package service

import (
	"fmt"
	"testing"

	"github.com/racecarparts/dashster/model"
)

func TestGitHubPRs(t *testing.T) {
	model.AppConfig = &model.Config{}
	cfg := model.AppConfig
	cfg.GithubPulls.Enabled = true
	cfg.GithubPulls.Organizations = []model.GithubOrg{{
		Name:         "",
		TeamNameSlug: "",
		Username:     "",
		AccessKey:    "",
		Team: model.GithubTeam{
			Id:   0,
			Name: "",
			Slug: "",
			Members: []model.GithubUser{{
				Login: "",
			}},
		},
	}}

	// this isn't a real test, it's just an easy place to invoke the service.
	sprs := PullRequests()
	fmt.Println(sprs)
}
