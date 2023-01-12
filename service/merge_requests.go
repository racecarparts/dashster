package service

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/goccy/go-json"

	"github.com/racecarparts/dashster/model"
)

const groupsUrlPath = "/groups"

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

	teamMembers, err := getTeamMembers(org)
	if err != nil {
		mrs.Message = fmt.Sprintf("could not get team members: %s", err.Error())
		fmt.Println(mrs.Message)
		return mrs, err
	}

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
			mrs := getProjectMRs(prj, org, teamMembers)
			// TODO: this should be a channel
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

	mrs.MyPRs = sortPullReqs(mrs.MyPRs)
	mrs.RequestedPRs = sortPullReqs(mrs.RequestedPRs)

	return mrs, nil
}

func getProjectMRs(project model.GLProject, org model.GitlabOrg, teamMembers map[string]model.GLUser) model.SimplePullRequests {
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
		if org.FilterMRsByGroup {
			if _, ok := teamMembers[mr.Author.Username]; !ok {
				continue
			}
		}

		notesPath := fmt.Sprintf("/projects/%d/merge_requests/%d/notes", project.Id, mr.Iid)
		notesUrl := fmt.Sprintf("%s%s", org.BaseUrl, notesPath)
		notesBody, err := getRequestBearerAuth(notesUrl, org.PrivateToken)
		if err != nil {
			mrs.Message = fmt.Sprint(notesUrl, err)
			fmt.Println(mrs.Message)
			return mrs
		}

		notes := []model.MergeRequestNote{}
		err = json.Unmarshal(notesBody, &notes)
		if err != nil {
			mrs.Message = fmt.Sprintf("could not unmarshall MR notes bytes: %s", err)
			fmt.Println(mrs.Message)
			return mrs
		}

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

		reviewers := buildReviewerMap(approvalState, notes)

		simplePR := model.SimplePullRequest{
			RepositoryName: project.Name,
			Number:         mr.Iid,
			User:           mr.Author.Username,
			Title:          mr.Title,
			Reviews:        reviewers,
			SHA:            mr.SHA[:7],
			WebURL:         mr.WebURL,
			UpdatedAt:      mr.UpdatedAt.In(time.Local).Format("02 Jan 03:04PM"),
			UpdatedAtTime:  mr.UpdatedAt,
		}

		if mr.Author.Username == org.Username {
			mrs.MyPRs = append(mrs.MyPRs, simplePR)
		} else {
			mrs.RequestedPRs = append(mrs.RequestedPRs, simplePR)
		}
	}

	return mrs
}

func getTeamMembers(org model.GitlabOrg) (map[string]model.GLUser, error) {
	if !org.FilterMRsByGroup {
		return map[string]model.GLUser{}, nil
	}

	groups, err := getConfigGroups(org)
	if err != nil {
		return map[string]model.GLUser{}, err
	}

	teamMembers := map[string]model.GLUser{}
	for _, group := range groups {
		membersUrl := fmt.Sprintf("%s%s/%d/members", org.BaseUrl, groupsUrlPath, group.Id)
		membersBody, err := getRequestBearerAuth(membersUrl, org.PrivateToken)
		if err != nil {
			return map[string]model.GLUser{}, err
		}

		members := []model.GLUser{}
		err = json.Unmarshal(membersBody, &members)
		if err != nil {
			return map[string]model.GLUser{}, err
		}

		for _, member := range members {
			teamMembers[member.Username] = member
		}
	}

	return teamMembers, nil
}

func getConfigGroups(org model.GitlabOrg) ([]model.GLGroup, error) {
	if len(org.GroupNames) == 0 {
		return []model.GLGroup{}, nil
	}

	groupsUrl := fmt.Sprintf("%s%s", org.BaseUrl, groupsUrlPath)
	groupsBody, err := getRequestBearerAuth(groupsUrl, org.PrivateToken)
	if err != nil {
		return []model.GLGroup{}, err
	}

	allGroups := []model.GLGroup{}
	err = json.Unmarshal(groupsBody, &allGroups)
	if err != nil {
		return []model.GLGroup{}, err
	}

	foundGroups := []model.GLGroup{}
	for _, groupName := range org.GroupNames {
		for _, group := range allGroups {
			if group.Name == groupName {
				foundGroups = append(foundGroups, group)
				break
			}
		}
	}

	return foundGroups, nil
}

func sortPullReqs(prs []model.SimplePullRequest) []model.SimplePullRequest {
	// sort.Slice(prs, func(i, j int) bool {
	// 	sortedByRepoName := false
	// 	sortedByMRNumber := false

	// 	sortedByRepoName = prs[i].RepositoryName < prs[j].RepositoryName
	// 	if prs[i].RepositoryName == prs[j].RepositoryName {
	// 		sortedByMRNumber = prs[i].Number < prs[j].Number
	// 		return sortedByMRNumber
	// 	}
	// 	return sortedByRepoName
	// })

	sort.Slice(prs, func(i, j int) bool {
		return prs[i].UpdatedAtTime.After(prs[j].UpdatedAtTime)
	})

	return prs
}

func buildReviewerMap(approvalState model.GLApprovalState, notes []model.MergeRequestNote) []model.PullReview {
	reviewers := make(map[model.GLUser]string, 0)
	commenters := make(map[model.GLUser]string, 0)

	for _, rule := range approvalState.Rules {
		for _, approver := range rule.ApprovedBy {
			reviewers[approver] = "✅"
		}
	}

	for _, comment := range notes {
		commenters[comment.Author] = "💬"
	}

	reviews := make([]model.PullReview, len(reviewers)+len(commenters))
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

	i = len(reviewers)
	for k, v := range commenters {
		r := model.PullReview{
			User: model.GithubUser{
				Login: k.Username,
			},
			State: v,
		}
		reviews[i] = r
		i += 1
	}

	sort.Slice(reviews[:], func(i, j int) bool {
		return reviews[i].State < reviews[j].State
	})

	return reviews
}
