package model

import "time"

type GLProject struct {
	Id    int            `json:"id"`
	Name  string         `json:"name"`
	Links GLProjectLinks `json:"_links"`
}

type GLProjectLinks struct {
	MergeRequestsLink string `json:"merge_requests"`
}

type MergeRequest struct {
	Id        int       `json:"id`
	Iid       int       `json:"iid"`
	Author    GLUser    `json:"author"`
	Title     string    `json:"title"`
	Draft     bool      `json:"draft"`
	SHA       string    `json:"sha"`
	Reviewers []GLUser  `json:"reviewers"`
	WebURL    string    `json:"web_url"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GLUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type GLGroup struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GLApprovalState struct {
	Rules []GLApprovalRule `json:"rules"`
}

type GLApprovalRule struct {
	Id         int      `json:"id"`
	Approved   bool     `json:"approved"`
	Users      []GLUser `json:"users"`
	ApprovedBy []GLUser `json:"approved_by"`
}

type MergeRequestNote struct {
	Id     int    `json:"id"`
	Author GLUser `json:"author"`
}
