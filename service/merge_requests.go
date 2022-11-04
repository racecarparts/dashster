package service

import (
	"fmt"
	"sync"

	"github.com/goccy/go-json"

	"github.com/racecarparts/dashster/model"
)

func MergeRequests() model.SimplePullRequests {
	if !model.AppConfig.Gitlab.Enabled {
		return model.SimplePullRequests{}
	}

	return gitlabMRs(model.AppConfig.Gitlab.Organizations)
}

func gitlabMRs(orgs []model.GitlabOrg) model.SimplePullRequests {
	mrs := model.SimplePullRequests{}

	for _, org := range orgs {
		var err error
		mrs, err = orgMRs(org)

		if err != nil {
			err := fmt.Errorf("problem with getting org MRs for %s: %s", org.Name, err.Error())
			fmt.Println(err)
			mrs.Message = err.Error()
			return mrs
		}
	}

	return mrs
}

func orgMRs(org model.GitlabOrg) (model.SimplePullRequests, error) {

	mrs := model.SimplePullRequests{}

	projectUrl := fmt.Sprintf("%s%s", org.BaseUrl, "/projects?starred=true")
	projectsBody, err := getRequestBearerAuth(projectUrl, org.PrivateToken)
	if err != nil {
		fmt.Println(projectUrl, err)
		return mrs, err
	}

	projects := []model.GLProject{}

	err = json.Unmarshal(projectsBody, &projects)
	if err != nil {
		fmt.Println("could not unmarshall project bytes", err)
		return mrs, err
	}

	projectMRs := make(map[int]model.SimplePullRequests, 0)
	var wg sync.WaitGroup
	for _, project := range projects {

		wg.Add(1)

		go func(prj model.GLProject) {
			defer wg.Done()
			mrs := getProjectMRs(prj, org)
			projectMRs[prj.Id] = mrs
		}(project)
	}

	wg.Wait()

	mrs.MyPRs = []model.SimplePullRequest{}
	mrs.RequestedPRs = []model.SimplePullRequest{}

	for _, prjMRs := range projectMRs {
		mrs.MyPRs = append(mrs.MyPRs, prjMRs.MyPRs...)
		mrs.RequestedPRs = append(mrs.RequestedPRs, prjMRs.RequestedPRs...)
		if len(prjMRs.Message) > 0 {
			mrs.Message = fmt.Sprintf("%s\n%s", mrs.Message, prjMRs.Message)
		}
	}

	return mrs, nil
}

func getProjectMRs(project model.GLProject, org model.GitlabOrg) model.SimplePullRequests {
	mrs := model.SimplePullRequests{
		Message:      "",
		MyPRs:        []model.SimplePullRequest{},
		RequestedPRs: []model.SimplePullRequest{},
	}

	mrUrl := fmt.Sprintf("%s?state=opened", project.Links.MergeRequestsLink)
	mrBody, err := getRequestBearerAuth(mrUrl, org.PrivateToken)
	if err != nil {
		fmt.Println(mrUrl, err)
		mrs.Message = err.Error()
		return mrs
	}

	mrList := []model.MergeRequest{}
	err = json.Unmarshal(mrBody, &mrList)
	if err != nil {
		mrs.Message = fmt.Sprintf("could not unmarshall merge request bytes: %s", err.Error())
		fmt.Println(mrs.Message)
		return mrs
	}

	for _, mr := range mrList {
		path := fmt.Sprintf("/projects/%d/merge_requests/%d/approval_state", project.Id, mr.Iid)
		approvalStateUrl := fmt.Sprintf("%s%s", org.BaseUrl, path)
		apprStateBody, err := getRequestBearerAuth(approvalStateUrl, org.PrivateToken)
		if err != nil {
			mrs.Message = fmt.Sprint(approvalStateUrl, err)
			fmt.Println(mrs.Message)
			return mrs
		}

		approvalState := model.GLApprovalState{}
		err = json.Unmarshal(apprStateBody, &approvalState)
		if err != nil {
			mrs.Message = fmt.Sprintf("could not unmarshall approval state bytes: %s", err)
			fmt.Println(mrs.Message)
			return mrs
		}

		reviewers := buildReviewerMap(approvalState)

		simplePR := model.SimplePullRequest{
			RepositoryName: project.Name,
			Number:         mr.Iid,
			User:           mr.Author.Username,
			Title:          mr.Title,
			Reviews:        reviewers,
			SHA:            mr.SHA[:7],
			WebURL:         mr.WebURL,
		}

		if mr.Author.Username == org.Username {
			mrs.MyPRs = append(mrs.MyPRs, simplePR)
		} else {
			mrs.RequestedPRs = append(mrs.RequestedPRs, simplePR)
		}
	}

	return mrs
}

func buildReviewerMap(approvalState model.GLApprovalState) []model.PullReview {
	reviewers := make(map[model.GLUser]string, 0)

	for _, rule := range approvalState.Rules {
		for _, approver := range rule.ApprovedBy {
			reviewers[approver] = "Approved"
		}
	}

	reviews := make([]model.PullReview, len(reviewers))
	i := 0
	for k, v := range reviewers {
		r := model.PullReview{
			User: model.GithubUser{
				Login: k.Username,
			},
			State: v,
		}
		reviews[i] = r
		i += 1
	}

	return reviews
}
