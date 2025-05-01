package service

import (
	"fmt"
	"gbserver/models"

	"testing"

	"github.com/stretchr/testify/assert"
)

// var repoRequest = &CreateRepoRequest{
// 	"test",
// 	"sample test repo",
// }

// var orgname = "gborg"
// var owner = "gbuser"

var gbService = GbService{GbStoreInstance: models.NewGbStore()}

func TestListRepo(t *testing.T) {
	type input struct {
		orgName string
		owner   string
	}
	tests := []struct {
		name  string
		input input
		//got   got
		wantResp []RepoResponse
		wantErr  error
	}{
		{
			name: "Test invalid user",
			input: input{
				orgName: "gborg",
				owner:   "gbInvalidUser",
			},
			wantErr: ErrOwnerNotFound,
		},
		{
			name: "Test invalid Org",
			input: input{
				orgName: "gbInvalidorg",
				owner:   "gbUser",
			},

			wantErr: ErrOrgNotFound,
		},
		{
			name: "Test repo results for valid org & user",
			input: input{
				orgName: "gborg",
				owner:   "gbuser",
			},

			wantResp: []RepoResponse{
				{
					ID: 1, Name: "gbrepo", Node_ID: "MDEwOlJlcG9zaXRvcnkxMjk2MjY5", Description: "gbuser repo",
					OwnerInfo: OwnerInfo{Login: "gbuser", ID: 1, NodeID: "MDQ6VXNlcjE=", UserType: "User"},
				},
			},
		},
	}
	for _, tt := range tests {
		resp, err := gbService.ListRepos(tt.input.orgName, tt.input.owner)
		if err != nil {
			assert.Equal(t, tt.wantErr.Error(), err.Error())
		} else {
			assert.Equal(t, tt.wantResp, resp)
		}
	}
}

func TestCreateRepo(t *testing.T) {
	type input struct {
		orgName string
		owner   string
		repoReq CreateRepoRequest
	}
	tests := []struct {
		name  string
		input input
		//got   got
		wantResp RepoResponse
		wantErr  error
	}{
		{
			name: "Test invalid user",
			input: input{
				orgName: "gborg",
				owner:   "gbInvalidUser",
			},
			wantErr: ErrOwnerNotFound,
		},
		{
			name: "Test invalid Org",
			input: input{
				orgName: "gbInvalidorg",
				owner:   "gbUser",
			},

			wantErr: ErrOrgNotFound,
		},
		{
			name: "Existing repo name",
			input: input{
				orgName: "gborg",
				owner:   "gbuser",
				repoReq: CreateRepoRequest{Name: "gbrepo", Description: "Test repo request"},
			},

			wantErr: ErrRepoAlreadyExists,
		},
		{
			name: "Test create repo results for valid org & user",
			input: input{
				orgName: "gborg",
				owner:   "gbuser",
				repoReq: CreateRepoRequest{Name: "testrepo", Description: "Test repo request"},
			},

			wantResp: RepoResponse{
				ID: 1, Name: "testrepo", Description: "Test repo request",
				OwnerInfo: OwnerInfo{"gbuser", 1, "MDQ6VXNlcjE=", "User"},
			},
		},
	}
	for _, tt := range tests {
		resp, err := gbService.CreateRepo(tt.input.orgName, tt.input.owner, &tt.input.repoReq)
		if err != nil {
			assert.Equal(t, tt.wantErr.Error(), err.Error())
		} else {
			assert.Equal(t, tt.wantResp.Name, resp.Name)
		}
	}
}

func TestDeleteRepo(t *testing.T) {
	type input struct {
		orgName  string
		owner    string
		repoName string
	}
	tests := []struct {
		name  string
		input input
		//got   got
		wantResp bool
		wantErr  error
	}{
		{
			name: "Test invalid user",
			input: input{
				orgName:  "gborg",
				owner:    "gbInvalidUser",
				repoName: "gbrepo",
			},
			wantErr: ErrOwnerNotFound,
		},
		{
			name: "Test invalid Org",
			input: input{
				orgName:  "gbInvalidorg",
				owner:    "gbUser",
				repoName: "gbrepo",
			},

			wantErr: ErrOrgNotFound,
		},
		{
			name: "Test invalid Repo",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "Invalidrepo",
			},

			wantErr: ErrRepoNotFound,
		},
		{
			name: "Test delete valid repo for valid org & user",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "gbrepo",
			},
			wantResp: true,
		},
	}
	for _, tt := range tests {
		resp, err := gbService.DeleteRepo(tt.input.orgName, tt.input.owner, tt.input.repoName)
		if err != nil {
			assert.Equal(t, tt.wantErr.Error(), err.Error())
		} else {
			assert.Equal(t, tt.wantResp, resp)
		}
	}
}

func TestCreateBranch(t *testing.T) {
	type input struct {
		orgName   string
		owner     string
		repoName  string
		branchReq CreateBranchRequest
	}
	tests := []struct {
		name  string
		input input
		//got   got
		wantResp CreateBranchResponse
		wantErr  error
	}{
		{
			name: "Test invalid user",
			input: input{
				orgName:  "gborg",
				owner:    "gbInvalidUser",
				repoName: "gbrepo",
			},
			wantErr: ErrOwnerNotFound,
		},
		{
			name: "Test invalid Org",
			input: input{
				orgName:  "gbInvalidorg",
				owner:    "gbUser",
				repoName: "gbrepo",
			},

			wantErr: ErrOrgNotFound,
		},
		{
			name: "Test invalid Repo",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "gbrepo1",
			},

			wantErr: ErrRepoNotFound,
		},
		{
			name: "Test create branch results for valid org, user & Invalid reponame",
			input: input{
				orgName:   "gborg",
				owner:     "gbuser",
				repoName:  "testrepo",
				branchReq: CreateBranchRequest{Ref: "refs/heads", SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd"},
			},

			wantErr: ErrInvalidBranchName,
		},
		{
			name: "Test create branch results for valid org, user & repo",
			input: input{
				orgName:   "gborg",
				owner:     "gbuser",
				repoName:  "testrepo",
				branchReq: CreateBranchRequest{Ref: "refs/heads/featureCD", SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd"},
			},

			wantResp: CreateBranchResponse{
				Ref:    "refs/heads/featureCD",
				NodeID: "XOgXav=s=AF86WNi9I2C=MY",
				URL:    "https://api.gbserver.com/repos/gbuser/gbrepo/git/commits/csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd",
				Object: CreateBranchObjectResponse{
					Type: "commit",
					SHA:  "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd",
					URL:  "https://api.gbserver.com/repos/gbuser/gbrepo/git/commits/csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd",
				},
			},
		},
		// {
		// 	name: "Existing branch name",
		// 	input: input{
		// 		orgName:   "gborg",
		// 		owner:     "gbuser",
		// 		repoName:  "testrepo",
		// 		branchReq: CreateBranchRequest{Ref: "refs/heads/gbbranch", SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd"},
		// 	},

		// 	wantErr: ErrBranchesAlreadyExists,
		// },
	}
	for _, tt := range tests {
		resp1, _ := gbService.CreateRepo(tt.input.orgName, tt.input.owner, &CreateRepoRequest{Name: "testrepo", Description: "Test repo request"})
		resp, err := gbService.CreateBranch(tt.input.orgName, tt.input.owner, tt.input.repoName, &tt.input.branchReq)
		if err != nil {
			assert.Equal(t, tt.wantErr.Error(), err.Error())
		} else {
			fmt.Println(resp1, resp)
			assert.Equal(t, tt.wantResp.Object.SHA, resp.Object.SHA)
		}
	}
}

func TestListBranches(t *testing.T) {
	type input struct {
		orgName  string
		owner    string
		repoName string
	}
	tests := []struct {
		name  string
		input input
		//got   got
		wantResp []ListBranchresponse
		wantErr  error
	}{
		{
			name: "Test invalid user",
			input: input{
				orgName: "gborg",
				owner:   "gbInvalidUser",
			},
			wantErr: ErrOwnerNotFound,
		},
		{
			name: "Test invalid Org",
			input: input{
				orgName: "gbInvalidorg",
				owner:   "gbUser",
			},

			wantErr: ErrOrgNotFound,
		},
		{
			name: "Test invalid Repo",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "Invalidrepo",
			},

			wantErr: ErrRepoNotFound,
		},
		{
			name: "Test branch results for valid org, user & repo",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "testrepo",
			},

			wantResp: []ListBranchresponse{
				{
					Name: "featureCD",
					Commit: CommitDetails{
						SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd",
						URL: "https://api.gbserver.com/repos/gbuser/testrepo/git/commits/csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd",
					},
					Protected: false,
				},
			},
		},
	}

	for _, tt := range tests {
		resp1, _ := gbService.CreateRepo(tt.input.orgName, tt.input.owner, &CreateRepoRequest{Name: "testrepo", Description: "Test repo request"})
		resp2, _ := gbService.CreateBranch(tt.input.orgName, tt.input.owner, tt.input.repoName, &CreateBranchRequest{Ref: "refs/heads/featureCD", SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd"})
		resp, err := gbService.ListBranches(tt.input.orgName, tt.input.owner, tt.input.repoName)
		fmt.Println(tt.name, resp, err)
		if err != nil {
			assert.Equal(t, tt.wantErr.Error(), err.Error())
		} else {
			fmt.Println(resp1, resp2, resp)
			assert.Equal(t, tt.wantResp[0].Name, resp[0].Name)
		}
	}
}

func TestDeleteBranch(t *testing.T) {
	type input struct {
		orgName    string
		owner      string
		repoName   string
		branchName string
	}
	tests := []struct {
		name  string
		input input
		//got   got
		wantResp bool
		wantErr  error
	}{
		{
			name: "Test invalid user",
			input: input{
				orgName:  "gborg",
				owner:    "gbInvalidUser",
				repoName: "testrepo",
			},
			wantErr: ErrOwnerNotFound,
		},
		{
			name: "Test invalid Org",
			input: input{
				orgName:  "gbInvalidorg",
				owner:    "gbUser",
				repoName: "testrepo",
			},

			wantErr: ErrOrgNotFound,
		},
		{
			name: "Test invalid Repo",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "Invalidrepo",
			},

			wantErr: ErrRepoNotFound,
		},
		{
			name: "Test invalid Repo",
			input: input{
				orgName:    "gborg",
				owner:      "gbuser",
				repoName:   "testrepo",
				branchName: "InvalidBranchName",
			},
			wantErr: ErrBranchesNotFound,
		},
		{
			name: "Test delete branch results for valid org, user & repo",
			input: input{
				orgName:    "gborg",
				owner:      "gbuser",
				repoName:   "testrepo",
				branchName: "featureCD",
			},

			wantResp: true,
		},
	}

	for _, tt := range tests {
		resp1, _ := gbService.CreateRepo(tt.input.orgName, tt.input.owner, &CreateRepoRequest{Name: "testrepo", Description: "Test repo request"})
		resp2, _ := gbService.CreateBranch(tt.input.orgName, tt.input.owner, tt.input.repoName, &CreateBranchRequest{Ref: "refs/heads/featureCD", SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd"})
		resp, err := gbService.DeleteBranch(tt.input.orgName, tt.input.owner, tt.input.repoName, tt.input.branchName)
		fmt.Println(tt.name, resp, err)
		if err != nil {
			assert.Equal(t, tt.wantErr.Error(), err.Error())
		} else {
			fmt.Println(resp1, resp2, resp)
			assert.Equal(t, tt.wantResp, resp)
		}
	}
}

func TestCreatePRs(t *testing.T) {
	type input struct {
		orgName  string
		owner    string
		repoName string
		prReq    PRRequest
	}
	tests := []struct {
		name  string
		input input
		//got   got
		wantResp PRResponse
		wantErr  error
	}{
		{
			name: "Test invalid user",
			input: input{
				orgName:  "gborg",
				owner:    "gbInvalidUser",
				repoName: "gbrepo",
			},
			wantErr: ErrOwnerNotFound,
		},
		{
			name: "Test invalid Org",
			input: input{
				orgName:  "gbInvalidorg",
				owner:    "gbUser",
				repoName: "gbrepo",
			},

			wantErr: ErrOrgNotFound,
		},
		{
			name: "Test invalid Repo",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "gbrepo1",
			},

			wantErr: ErrRepoNotFound,
		},
		{
			name: "Test create PR for valid org, user & Invalid reponame",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "testrepo",
				prReq:    PRRequest{Title: "Amazing new feature", Body: "Please pull these awesome changes in!", Head: "gbuser:featureD", Base: "master"},
			},

			wantErr: ErrBranchesNotFound,
		},
		{
			name: "Test create branch results for valid org, user & repo",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "testrepo",
				prReq:    PRRequest{Title: "Amazing new feature", Body: "Please pull these awesome changes in!", Head: "gbuser:featureCD", Base: "master"},
			},

			wantResp: PRResponse{
				URL:    "https://api.gbserver.com/repos/gbuser/gbrepo/1534409025785096406",
				ID:     "1534409025785096406",
				NodeID: "w5PCfNJBg=pJfWjYn6eceB0",
				Title:  "Amazing new feature",
				Body:   "Please pull these awesome changes in!",
				State:  "open",
				User: OwnerInfo{
					Login:    "gbuser",
					ID:       1,
					NodeID:   "MDQ6VXNlcjE=",
					UserType: "User",
				},
				Commits:      0,
				Additions:    17,
				Deletions:    24,
				ChangedFiles: 41,
				Head: baseHeadPRResponse{
					Ref: "master",
					SHA: "aa218f56b14c9653891f9e74264a383fa43fefbd",
					User: OwnerInfo{
						Login:    "gbuser",
						ID:       1,
						NodeID:   "MDQ6VXNlcjE=",
						UserType: "User",
					},
					Repo: "gbrepo",
				},
				Base: baseHeadPRResponse{
					Ref: "featureCD",
					SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd",
					User: OwnerInfo{
						Login:    "gbuser",
						ID:       1,
						NodeID:   "MDQ6VXNlcjE=",
						UserType: "User",
					},
					Repo: "gbrepo",
				},
			},
		},
	}
	for _, tt := range tests {
		resp1, _ := gbService.CreateRepo(tt.input.orgName, tt.input.owner, &CreateRepoRequest{Name: "testrepo", Description: "Test repo request"})
		resp2, _ := gbService.CreateBranch(tt.input.orgName, tt.input.owner, tt.input.repoName, &CreateBranchRequest{Ref: "refs/heads/featureCD", SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd"})
		resp3, _ := gbService.CreateBranch(tt.input.orgName, tt.input.owner, tt.input.repoName, &CreateBranchRequest{Ref: "refs/heads/master", SHA: "abcgsd2esdf56b14c9653891f9e74264a383fa43fefbd"})
		resp, err := gbService.CreatePR(tt.input.orgName, tt.input.owner, tt.input.repoName, &tt.input.prReq)
		fmt.Println(tt.name, "..", resp, err)
		if err != nil {

			assert.Equal(t, tt.wantErr.Error(), err.Error())
		} else {
			fmt.Println(resp1, resp2, resp3)
			assert.Equal(t, tt.wantResp.Title, resp.Title)
		}
	}
}

func TestListPRs(t *testing.T) {
	type input struct {
		orgName  string
		owner    string
		repoName string
	}
	tests := []struct {
		name  string
		input input
		//got   got
		wantResp []PRResponse
		wantErr  error
	}{
		{
			name: "Test invalid user",
			input: input{
				orgName:  "gborg",
				owner:    "gbInvalidUser",
				repoName: "gbrepo",
			},
			wantErr: ErrOwnerNotFound,
		},
		{
			name: "Test invalid Org",
			input: input{
				orgName:  "gbInvalidorg",
				owner:    "gbUser",
				repoName: "gbrepo",
			},

			wantErr: ErrOrgNotFound,
		},
		{
			name: "Test invalid Repo",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "gbrepo1",
			},

			wantErr: ErrRepoNotFound,
		},
		{
			name: "Test list pr with valid org, user & repo",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "testrepo",
			},

			wantResp: []PRResponse{
				{
					URL:    "https://api.gbserver.com/repos/gbuser/gbrepo/1534409025785096406",
					ID:     "1534409025785096406",
					NodeID: "w5PCfNJBg=pJfWjYn6eceB0",
					Title:  "Amazing new feature",
					Body:   "Please pull these awesome changes in!",
					State:  "open",
					User: OwnerInfo{
						Login:    "gbuser",
						ID:       1,
						NodeID:   "MDQ6VXNlcjE=",
						UserType: "User",
					},
					Commits:      0,
					Additions:    17,
					Deletions:    24,
					ChangedFiles: 41,
					Head: baseHeadPRResponse{
						Ref: "master",
						SHA: "aa218f56b14c9653891f9e74264a383fa43fefbd",
						User: OwnerInfo{
							Login:    "gbuser",
							ID:       1,
							NodeID:   "MDQ6VXNlcjE=",
							UserType: "User",
						},
						Repo: "gbrepo",
					},
					Base: baseHeadPRResponse{
						Ref: "featureCD",
						SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd",
						User: OwnerInfo{
							Login:    "gbuser",
							ID:       1,
							NodeID:   "MDQ6VXNlcjE=",
							UserType: "User",
						},
						Repo: "gbrepo",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		resp1, _ := gbService.CreateRepo(tt.input.orgName, tt.input.owner, &CreateRepoRequest{Name: "testrepo", Description: "Test repo request"})
		resp2, _ := gbService.CreateBranch(tt.input.orgName, tt.input.owner, tt.input.repoName, &CreateBranchRequest{Ref: "refs/heads/featureCD", SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd"})
		resp3, _ := gbService.CreateBranch(tt.input.orgName, tt.input.owner, tt.input.repoName, &CreateBranchRequest{Ref: "refs/heads/master", SHA: "abcgsd2esdf56b14c9653891f9e74264a383fa43fefbd"})
		resp4, _ := gbService.CreatePR(tt.input.orgName, tt.input.owner, tt.input.repoName, &PRRequest{Title: "Amazing new feature", Body: "Please pull these awesome changes in!", Head: "gbuser:featureCD", Base: "master"})
		resp, err := gbService.ListPRs(tt.input.orgName, tt.input.owner, tt.input.repoName)
		fmt.Println(tt.name, "..", resp, err)
		if err != nil {

			assert.Equal(t, tt.wantErr.Error(), err.Error())
		} else {
			fmt.Println(resp1, resp2, resp3, resp4)
			assert.Equal(t, tt.wantResp[0].Title, resp[0].Title)
		}
	}
}

func TestApprovePRs(t *testing.T) {
	type input struct {
		orgName  string
		owner    string
		repoName string
	}
	tests := []struct {
		name  string
		input input
		//got   got
		wantResp PRResponse
		wantErr  error
	}{
		{
			name: "Test approve pr with valid org, user & repo",
			input: input{
				orgName:  "gborg",
				owner:    "gbuser",
				repoName: "testrepo",
			},

			wantResp: PRResponse{
				URL:    "https://api.gbserver.com/repos/gbuser/gbrepo/1534409025785096406",
				ID:     "1534409025785096406",
				NodeID: "w5PCfNJBg=pJfWjYn6eceB0",
				Title:  "Amazing new feature",
				Body:   "Please pull these awesome changes in!",
				State:  "approved",
				User: OwnerInfo{
					Login:    "gbuser",
					ID:       1,
					NodeID:   "MDQ6VXNlcjE=",
					UserType: "User",
				},
				Commits:      0,
				Additions:    17,
				Deletions:    24,
				ChangedFiles: 41,
				Head: baseHeadPRResponse{
					Ref: "master",
					SHA: "aa218f56b14c9653891f9e74264a383fa43fefbd",
					User: OwnerInfo{
						Login:    "gbuser",
						ID:       1,
						NodeID:   "MDQ6VXNlcjE=",
						UserType: "User",
					},
					Repo: "gbrepo",
				},
				Base: baseHeadPRResponse{
					Ref: "featureCD",
					SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd",
					User: OwnerInfo{
						Login:    "gbuser",
						ID:       1,
						NodeID:   "MDQ6VXNlcjE=",
						UserType: "User",
					},
					Repo: "gbrepo",
				},
			},
		},
	}
	for _, tt := range tests {
		resp1, _ := gbService.CreateRepo(tt.input.orgName, tt.input.owner, &CreateRepoRequest{Name: "testrepo", Description: "Test repo request"})
		resp2, _ := gbService.CreateBranch(tt.input.orgName, tt.input.owner, tt.input.repoName, &CreateBranchRequest{Ref: "refs/heads/featureCD", SHA: "csdsdsdsdf56b14c9653891f9e74264a383fa43fefbd"})
		resp3, _ := gbService.CreateBranch(tt.input.orgName, tt.input.owner, tt.input.repoName, &CreateBranchRequest{Ref: "refs/heads/master", SHA: "abcgsd2esdf56b14c9653891f9e74264a383fa43fefbd"})
		resp4, _ := gbService.CreatePR(tt.input.orgName, tt.input.owner, tt.input.repoName, &PRRequest{Title: "Amazing new feature", Body: "Please pull these awesome changes in!", Head: "gbuser:featureCD", Base: "master"})
		fmt.Println(resp1, resp2, resp3, resp4)
		pullNumber := resp4.ID
		updatePRReq := PRRequest{State: "approved"}
		resp, err := gbService.ApprovePR(tt.input.orgName, tt.input.owner, tt.input.repoName, pullNumber, &updatePRReq)
		fmt.Println(tt.name, "..", resp, err)
		if err != nil {
			assert.Equal(t, tt.wantErr.Error(), err.Error())
		} else {
			assert.Equal(t, tt.wantResp.State, resp.State)
		}
	}
}
