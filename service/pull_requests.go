package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/racecarparts/dashster/model"
	"sync"
)

//func openPRListAsStr(orgs []model.GithubOrg, nextInterval time.Duration) string {
//    openPRs, err := getOpenGithubPRs(orgs)
//    if err != nil {
//        return "Problem getting open PRs: " + err.Error()
//    }
//    return strings.Join(openPRs, "\n")
//}

func PullRequests() model.SimplePullRequests {
	return getGithubPRs(model.AppConfig.GithubPulls.Organizations)
}

func populateGithubTeams(orgs []model.GithubOrg) []model.GithubOrg {
	for i, org := range orgs {
		authToken := org.AccessKey

		team := model.GithubTeam{}

		teamUrl := fmt.Sprintf("https://api.github.com/orgs/%s/teams/%s", org.Name, org.TeamNameSlug)
		teamBody, err := getRequestBasicAuth(teamUrl, authToken)
		if err != nil {
			return orgs
		}
		if len(teamBody) < 10 {
			continue
		}

		err = json.Unmarshal(teamBody, &team)
		org.Team.Members = make([]model.GithubUser, 0)

		page := 1
		membersUrl := ""
		for {
			var membersBody []byte
			membersUrl = fmt.Sprintf("https://api.github.com/orgs/%s/teams/%s/members?page=%d&per_page=100", org.Name, org.TeamNameSlug, page)

			membersBody, err = getRequestBasicAuth(membersUrl, authToken)
			if err != nil {
				break
			}
			if len(membersBody) < 10 {
				break
			}

			teamMembers := make([]model.GithubUser, 0)
			err = json.Unmarshal(membersBody, &teamMembers)
			if err != nil {
				break
			}
			team.Members = append(team.Members, teamMembers...)

			page += 1
		}
		org.Team = team
		orgs[i] = org
	}

	return orgs
}

// getGithubOrgRepos get the repos for a teams in the orgs, then all the repos in the orgs.  The lists are de-duped
func getGithubOrgRepos(orgs []model.GithubOrg) (teamRepos map[int]model.Repo, otherRepos map[int]model.Repo, err error) {

	for _, org := range orgs {
		// if this org has a team, then we need to create another version without a team, so the repos can be found
		if org.TeamNameSlug != "" {
			org.TeamNameSlug = ""
			orgs = append(orgs, org)
		}
	}

	teamRepos = make(map[int]model.Repo)
	otherRepos = make(map[int]model.Repo)
	for _, org := range orgs {
		page := 1
		url := ""
		hasTeamRepos := org.TeamNameSlug != ""
		authToken := org.AccessKey
		for {
			var moreRepos []model.Repo
			var body []byte

			if hasTeamRepos {
				url = fmt.Sprintf("https://api.github.com/orgs/%s/teams/%s/repos?page=%d&per_page=100", org.Name, org.TeamNameSlug, page)
			} else {
				url = fmt.Sprintf("https://api.github.com/orgs/%s/repos?page=%d&per_page=100", org.Name, page)
			}

			body, err = getRequestBasicAuth(url, authToken)
			if err != nil {
				return
			}
			if len(body) < 10 {
				break
			}

			err = json.Unmarshal(body, &moreRepos)
			if err != nil {
				return
			}

			// prevent duplicates
			for _, repo := range moreRepos {
				_, okTeam := teamRepos[repo.Id]
				_, okOther := otherRepos[repo.Id]

				if okTeam || okOther {
					continue
				}
				repo.Username = org.Username
				repo.OrgName = org.Name
				repo.AuthToken = authToken

				if hasTeamRepos {
					repo.IsTeamRepo = true
					teamRepos[repo.Id] = repo
				} else {
					otherRepos[repo.Id] = repo
				}

			}

			page += 1
		}
	}

	return
}

func getGithubPullReviews(repo model.Repo, pr model.PullRequest) ([]model.PullReview, error) {
	reviews := make([]model.PullReview, 0)
	page := 1
	for {
		body, err := getRequestBasicAuth(fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/reviews?page=%d&per_page=100", repo.OrgName, repo.Name, pr.Number, page), repo.AuthToken)
		if err != nil {
			return reviews, err
		}
		if len(body) < 10 {
			break
		}

		reviewPage := make([]model.PullReview, 0)
		err = json.Unmarshal(body, &reviewPage)
		if err != nil {
			return reviews, err
		}

		reviews = append(reviews, reviewPage...)

		page += 1
	}

	return reviews, nil

}

func createSimplePR(repo model.Repo, pr model.PullRequest) (model.SimplePullRequest, error) {
	simplePR := model.SimplePullRequest{
		RepositoryName: repo.Name,
		Number:         pr.Number,
		User:           pr.User.Login,
		Title:          pr.Title,
		ReviewUrl:      pr.HtmlUrl,
		IsDraft:        pr.Draft,
		SHA:            pr.Head.SHA[0:7],
	}

	reviews, err := getGithubPullReviews(repo, pr)
	if err != nil {
		return simplePR, err
	}

	simplePR.Reviews = reviews
	return simplePR, nil
}

func getGithubPRs(orgs []model.GithubOrg) model.SimplePullRequests {
	prs := model.SimplePullRequests{}

	teamRepos, otherRepos, err := getGithubOrgRepos(orgs)
	orgs = populateGithubTeams(orgs)
	if err != nil {
		prs.Message = err.Error()
		return prs
	}

	iterateRepos := func(repos map[int]model.Repo) (myPRs []model.SimplePullRequest, requestedPRs []model.SimplePullRequest, err error) {
		myPRs = make([]model.SimplePullRequest, 0)
		requestedPRs = make([]model.SimplePullRequest, 0)

		for _, repo := range repos {
			page := 1
			for {
				body, err := getRequestBasicAuth(fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?page=%d&per_page=100", repo.OrgName, repo.Name, page), repo.AuthToken)
				if err != nil {
					return myPRs, requestedPRs, err
				}
				if len(body) < 10 {
					break
				}
				repoPRs := make([]model.PullRequest, 0)
				err = json.Unmarshal(body, &repoPRs)
				if err != nil {
					return myPRs, requestedPRs, err
				}

			PRLoop:
				for _, pr := range repoPRs {
					if pr.Assignee.Login == repo.Username {
						simplePR, err := createSimplePR(repo, pr)
						if err != nil {
							return myPRs, requestedPRs, errors.New(fmt.Sprintf("error adding review info to PR %s: %s", pr.HtmlUrl, err.Error()))
						}
						myPRs = append(myPRs, simplePR)
						continue PRLoop
					}

					for _, reviewer := range pr.RequestedReviewers {
						for _, org := range orgs {
							if reviewer.Login == org.Username {
								simplePR, err := createSimplePR(repo, pr)
								if err != nil {
									return myPRs, requestedPRs, errors.New(fmt.Sprintf("error adding review info to PR %s: %s", pr.HtmlUrl, err.Error()))
								}
								requestedPRs = append(requestedPRs, simplePR)
								continue PRLoop
							}
						}
					}

					for _, team := range pr.RequestedTeams {
						for _, org := range orgs {
							if team.Slug == org.TeamNameSlug {
								simplePR, err := createSimplePR(repo, pr)
								if err != nil {
									return myPRs, requestedPRs, errors.New(fmt.Sprintf("error adding review info to PR %s: %s", pr.HtmlUrl, err.Error()))
								}
								requestedPRs = append(requestedPRs, simplePR)
								continue PRLoop
							}
						}
					}

					for _, org := range orgs {
						for _, teamMember := range org.Team.Members {
							if teamMember.Login == pr.Assignee.Login {
								simplePR, err := createSimplePR(repo, pr)
								if err != nil {
									return myPRs, requestedPRs, errors.New(fmt.Sprintf("error adding review info to PR %s: %s", pr.HtmlUrl, err.Error()))
								}
								requestedPRs = append(requestedPRs, simplePR)
								continue PRLoop
							}
						}
					}
				}
				page += 1
			}
		}

		return
	}

	var myTeamPRs []model.SimplePullRequest
	var requestedTeamPRs []model.SimplePullRequest
	var myOtherPRs []model.SimplePullRequest
	var requestedOtherPRs []model.SimplePullRequest

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		myTeamPRs, requestedTeamPRs, err = iterateRepos(teamRepos)

	}()

	go func() {
		defer wg.Done()
		myOtherPRs, requestedOtherPRs, err = iterateRepos(otherRepos)
	}()

	if err != nil {
		prs.Message = err.Error()
		return prs
	}

	wg.Wait()

	myPRs := append(myTeamPRs, myOtherPRs...)
	requestedPRs := append(requestedTeamPRs, requestedOtherPRs...)

	prs.MyPRs = myPRs
	prs.RequestedPRs = requestedPRs

	return prs
}

func getOpenGithubPRs(orgs []model.GithubOrg) ([]string, error) {
	pullRequests := make([]string, 0)

	teamRepos, otherRepos, err := getGithubOrgRepos(orgs)
	orgs = populateGithubTeams(orgs)
	if err != nil {
		pullRequests = append(pullRequests, err.Error())
		return pullRequests, err
	}

	appendPRInfo := func(prInfos []string, repo model.Repo, pr model.PullRequest) []string {
		draftStatus := ""
		if pr.Draft {
			draftStatus = "(DRAFT) "
		}
		reviews, err := getGithubPullReviews(repo, pr)
		reviewStates := ""
		if err != nil {
			reviewStates = fmt.Sprintf("!! %s: %s", "error getting pr reviews", err.Error())
		}

		for _, review := range reviews {
			reviewStates = fmt.Sprintf("%s %s:%s ", reviewStates, review.User.Login, review.State)
		}
		return append(prInfos, fmt.Sprintf("  * %s: \n    %s%s\n    Reviews:%s\n    %s\n", repo.Name, draftStatus, pr.Title, reviewStates, pr.HtmlUrl))
	}

	iterateRepos := func(repos map[int]model.Repo) (myPRs []string, requestedPRs []string, err error) {
		myPRs = make([]string, 0)
		requestedPRs = make([]string, 0)

		for _, repo := range repos {
			page := 1
			for {
				body, err := getRequestBasicAuth(fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?page=%d&per_page=100", repo.OrgName, repo.Name, page), repo.AuthToken)
				if err != nil {
					pullRequests = append(pullRequests, err.Error())
					return myPRs, requestedPRs, err
				}
				if len(body) < 10 {
					break
				}
				repoPRs := make([]model.PullRequest, 0)
				err = json.Unmarshal(body, &repoPRs)
				if err != nil {
					pullRequests = append(pullRequests, err.Error())
					return myPRs, requestedPRs, err
				}

			PRLoop:
				for _, pr := range repoPRs {
					if pr.Assignee.Login == repo.Username {
						myPRs = appendPRInfo(myPRs, repo, pr)
						continue PRLoop
					}

					for _, reviewer := range pr.RequestedReviewers {
						for _, org := range orgs {
							if reviewer.Login == org.Username {
								requestedPRs = appendPRInfo(requestedPRs, repo, pr)
								continue PRLoop
							}
						}
					}

					for _, team := range pr.RequestedTeams {
						for _, org := range orgs {
							if team.Slug == org.TeamNameSlug {
								requestedPRs = appendPRInfo(requestedPRs, repo, pr)
								continue PRLoop
							}
						}
					}

					for _, org := range orgs {
						for _, teamMember := range org.Team.Members {
							if teamMember.Login == pr.Assignee.Login {
								requestedPRs = appendPRInfo(requestedPRs, repo, pr)
								continue PRLoop
							}
						}
					}
				}
				page += 1
			}
		}

		return
	}

	myTeamPRs, requestedTeamPRs, err := iterateRepos(teamRepos)
	myOtherPRs, requestedOtherPRs, err := iterateRepos(otherRepos)

	myPRs := append(myTeamPRs, myOtherPRs...)
	requestedPRs := append(requestedTeamPRs, requestedOtherPRs...)

	pullRequests = append(pullRequests, "  My PRs\n  ------")
	pullRequests = append(pullRequests, myPRs...)
	pullRequests = append(pullRequests, "  Requested PRs\n  -------------")
	pullRequests = append(pullRequests, requestedPRs...)

	return pullRequests, nil
}
