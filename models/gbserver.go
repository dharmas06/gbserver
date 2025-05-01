package models

import "sync"

type User struct {
	ID        int      `json:"id"`
	LoginName string   `json:"name"`
	OrgID     int      `json:"org_id"`
	NodeID    string   `json:"nodeId"`
	UserType  string   `json:"type"`
	Repos     []string `json:"repos"`
}

type Organization struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Users      []string `json:"users"`
	Repos      []string `json:"repos"`
	ReposCount int
}

type Repository struct {
	ID          int      `json:"id"`
	Node_ID     string   `json:"node_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	OrgName     string   `json:"org_name"`
	UserName    string   `json:"user_name"`
	Branches    []string `json:"branches"`
	TotalPRs    int
	PrIDs       []string
}

type CommitDetails struct {
	SHA string
	URL string
}

type Branch struct {
	ID            int    `json:"id"`
	RepoName      string `json:"repo_id"`
	Name          string `json:"name"`
	NodeID        string `json:"nodeID"`
	URL           string `json:"url"`
	Protected     bool   `json:"protected"`
	CommitInfo    CommitDetails
	PullRequestID string // should be random characters encoded characters of orgname + owner+reponame + prid
}

type PullRequest struct {
	NodeID       string `json:"nodeID"`
	URL          string `json:"url"`
	ID           string `json:"id"`
	RepoName     string `json:"repo_name"`
	FromBranch   string `json:"from_branch"`
	ToBranch     string `json:"to_branch"`
	AuthorID     int    `json:"author_id"`
	State        string `json:"status"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	Commits      int    `json:"commits"`
	Additions    int    `json:"additions"`
	Deletions    int    `json:"deletions"`
	ChangedFiles int    `json:"changed_files"`
}

type GbStore struct {
	MU           sync.RWMutex
	Users        map[string]*User
	Orgs         map[string]*Organization
	Repos        map[string]*Repository
	Branches     map[string]*Branch
	PullRequests map[string]*PullRequest
}

func NewGbStore() *GbStore {
	gbStore := &GbStore{
		Users:        make(map[string]*User),
		Orgs:         make(map[string]*Organization),
		Repos:        make(map[string]*Repository),
		Branches:     make(map[string]*Branch),
		PullRequests: make(map[string]*PullRequest),
	}

	gbStore.Users["gborg/gbuser"] = &User{ID: 1, LoginName: "gbuser", OrgID: 1, NodeID: "MDQ6VXNlcjE=", UserType: "User", Repos: []string{"gbrepo"}}
	gbStore.Orgs["gborg"] = &Organization{ID: 1, Name: "gborg", Users: []string{"gbuser"}, Repos: []string{"gbrepo"}, ReposCount: 1}
	gbStore.Repos["gborg/gbuser/gbrepo"] = &Repository{ID: 1, Name: "gbrepo", Node_ID: "MDEwOlJlcG9zaXRvcnkxMjk2MjY5", Description: "gbuser repo",
		OrgName: "gborg", UserName: "gbuser", Branches: []string{"master", "gbbranch"}, TotalPRs: 1, PrIDs: []string{"1534407926273468195"}}
	gbStore.Branches["gborg/gbuser/gbrepo/gbbranch"] = &Branch{ID: 1, RepoName: "gbrepo", Name: "gbbranch", NodeID: "NOSKDK8SDJSDHSD92KDkcy9mZWF0dXJlQQ==", URL: "https://api.gbserver.com/repos/gbuser/gbrepo/git/refs/heads/gbbranch",
		CommitInfo:    CommitDetails{SHA: "bchdjsd9jdowjd29ejiwd8y3hd3a383fa43fefbd", URL: "https://api.gbserver.com/repos/gbuser/gbrepo/git/commits/bchdjsd9jdowjd29ejiwd8y3hd3a383fa43fefbd"},
		PullRequestID: "1534407926273468195",
	}
	gbStore.Branches["gborg/gbuser/gbrepo/master"] = &Branch{ID: 2, RepoName: "gbrepo", Name: "master", NodeID: "MDM6UmVmcmVmcy9oZWFkcy9mZWF0dXJlQQ==", URL: "https://api.gbserver.com/repos/gbuser/gbrepo/git/refs/heads/master",
		CommitInfo: CommitDetails{SHA: "aa218f56b14c9653891f9e74264a383fa43fefbd", URL: "https://api.gbserver.com/repos/gbuser/gbrepo/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"}}
	gbStore.PullRequests["1534407926273468195"] = &PullRequest{ID: "1534407926273468195", NodeID: "MDExOlB1bGxSZXF1ZXN0MQ==",
		URL:      "https://api.gbserver.com/repos/gbuser/gbrepo/pulls/1",
		RepoName: "gbrepo", FromBranch: "gbuser:gbbranch",
		ToBranch: "master", AuthorID: 1, State: "open", Commits: 10,
		Title:        "Amazing new feature",
		Body:         "Please pull these awesome changes in!",
		Additions:    100,
		Deletions:    7,
		ChangedFiles: 23,
	}

	return gbStore
}
