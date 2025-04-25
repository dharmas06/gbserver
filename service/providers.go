package service

var gbServerOrg = map[string][]OwnerInfo{
	"gb": []OwnerInfo{
		{Login: "gbuser", ID: 1, NodeID: "MDQ6VXNlcjE=", UserType: "User"},
	},
	"testorg": []OwnerInfo{
		{Login: "testuser", ID: 2, NodeID: "NHUjIDklcjE=", UserType: "User"},
	},
}

var gbServerOwners = []OwnerInfo{
	{Login: "gbuser", ID: 1, NodeID: "MDQ6VXNlcjE=", UserType: "User"},
	{Login: "testuser", ID: 2, NodeID: "NHUjIDklcjE=", UserType: "User"},
}

var gbServerRepoData = []RepoResponse{
	{ID: 12121, Node_ID: "MDEwOlJlcG9zaXRvcnkxMjk2MjY5", Name: "gbrepo", Description: "This is a gb repo",
		OwnerInfo: OwnerInfo{Login: "gbuser", ID: 1, NodeID: "MDQ6VXNlcjE=", UserType: "User"},
	},
	{ID: 12345, Node_ID: "Nhdg72hdko9zaXRvcnkxMjk2MjY5", Name: "gbtestrepo", Description: "This is a gb repo 2",
		OwnerInfo: OwnerInfo{Login: "gbuser", ID: 1, NodeID: "MDQ6VXNlcjE=", UserType: "User"},
	},
	{ID: 13285, Node_ID: "ODHjnHyjDFj7hsRvcnkxMjk2MjY5", Name: "testrepo", Description: "This is a test repo",
		OwnerInfo: OwnerInfo{Login: "testuser", ID: 2, NodeID: "NHUjIDklcjE=", UserType: "User"},
	},
	{ID: 32516, Node_ID: "HDNKC8989jkncdmkHJMKMjk2MjY5", Name: "testrepo1", Description: "This is a test repo 1",
		OwnerInfo: OwnerInfo{Login: "testuser", ID: 2, NodeID: "NHUjIDklcjE=", UserType: "User"},
	},
}

var gbServerBranchData = []CreateBranchResponse{
	{Ref: "refs/heads/master", NodeID: "NOSKDK8SDJSDHSD92KDkcy9mZWF0dXJlQQ==",
		URL: "https://api.github.com/repos/gbuser/gbrepo/git/refs/heads/master",
		Object: CreateBranchObjectResponse{Type: "commit", SHA: "bchdjsd9jdowjd29ejiwd8y3hd3a383fa43fefbd",
			URL: "https://api.github.com/repos/gbuser/gbrepo/git/commits/bchdjsd9jdowjd29ejiwd8y3hd3a383fa43fefbd"}},
	{Ref: "refs/heads/gbbranch", NodeID: "MDM6UmVmcmVmcy9oZWFkcy9mZWF0dXJlQQ==",
		URL: "https://api.github.com/repos/gbuser/gbrepo/git/refs/heads/gbbranch",
		Object: CreateBranchObjectResponse{Type: "commit", SHA: "aa218f56b14c9653891f9e74264a383fa43fefbd",
			URL: "https://api.github.com/repos/gbuser/gbrepo/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"}},
	{Ref: "refs/heads/featureX", NodeID: "FMTUGRm4m4mEy3RZ33D3D3DD3F0dXJlQQ==", URL: "https://api.github.com/repos/testuser/testrepo/git/refs/heads/featureX"},
	{Ref: "refs/heads/featureY", NodeID: "PDMFUmRmGmTmGyBSZ4F5cy9mZWF0dXJlQQ==", URL: "https://api.github.com/repos/testuser/testrepo/git/refs/heads/featureY"},

	// {
	// 	"Ref": "refs/heads/featureA",
	// 	"Node_ID": "MDM6UmVmcmVmcy9oZWFkcy9mZWF0dXJlQQ==",
	// 	"URL": "https://api.github.com/repos/octocat/Hello-World/git/refs/heads/featureA",
	// 	"object": {
	// 	  "type": "commit",
	// 	  "SHA": "aa218f56b14c9653891f9e74264a383fa43fefbd",
	// 	  "URL": "https://api.github.com/repos/octocat/Hello-World/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
	// 	}
	//   }
}

var gbSeverBranchListData = []ListBranchresponse{
	{Name: "master", Commit: CommitDetails{SHA: "bchdjsd9jdowjd29ejiwd8y3hd3a383fa43fefbd", URL: "https://api.github.com/repos/gbuser/gbrepo/commits/bchdjsd9jdowjd29ejiwd8y3hd3a383fa43fefbd"}, Protected: true},
	{Name: "gbbranch", Commit: CommitDetails{SHA: "aa218f56b14c9653891f9e74264a383fa43fefbd", URL: "https://api.github.com/repos/gbuser/gbrepo/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"}, Protected: true},
	{Name: "devbranch", Commit: CommitDetails{SHA: "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc", URL: "https://api.github.com/repos/testuser/testrepo/commits/c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"}, Protected: true},
	{Name: "devbranch2", Commit: CommitDetails{SHA: "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc", URL: "https://api.github.com/repos/testuser/testrepo/commits/c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"}, Protected: true},
	// type commitDetails struct {
	// 	SHA string
	// 	URL string
	// }

	// type ListBranchresponse struct {
	// 	Name      string
	// 	commit    commitDetails
	// 	Protected bool
	// 	// {
	// 	"Name": "master",
	// 	"commit": {
	// 	  "SHA": "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc",
	// 	  "URL": "https://api.github.com/repos/admin/Hello-World/commits/c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	// 	},
	// 	"Protected": true,
}

var gbServerListPR = []PRResponse{
	{URL: "https://api.github.com/repos/gbuser/gbrepo/pulls/1",
		ID:           1,
		NodeID:       "MDExOlB1bGxSZXF1ZXN0MQ==",
		Title:        "Amazing new feature",
		Body:         "Please pull these awesome changes in!",
		State:        "open",
		User:         OwnerInfo{Login: "gbuser", ID: 1, NodeID: "MDQ6VXNlcjE=", UserType: "User"},
		Commits:      10,
		Additions:    100,
		Deletions:    7,
		ChangedFiles: 23,
		// head : branch has implemented changes
		Head: baseHeadPRResponse{
			Label: "gbuser:featuredAbranch",
			Ref:   "featuredAbranch",
			SHA:   "defh7rjk9sdjsdk9j2dksdl4264a383fa43fefbd",
			User:  OwnerInfo{Login: "gbuser", ID: 1, NodeID: "MDQ6VXNlcjE=", UserType: "User"},
			Repo: RepoResponse{
				ID: 12121, Node_ID: "MDEwOlJlcG9zaXRvcnkxMjk2MjY5", Name: "gbrepo", Description: "This is a gb repo",
				OwnerInfo: OwnerInfo{Login: "gbuser", ID: 1, NodeID: "MDQ6VXNlcjE=", UserType: "User"},
			},
		},
		// base: where changes need to be added.
		Base: baseHeadPRResponse{
			Label: "gbuser:gbbranch",
			Ref:   "gbbranch",
			SHA:   "aa218f56b14c9653891f9e74264a383fa43fefbd",
			User:  OwnerInfo{Login: "gbuser", ID: 1, NodeID: "MDQ6VXNlcjE=", UserType: "User"},
			Repo: RepoResponse{ID: 12121, Node_ID: "MDEwOlJlcG9zaXRvcnkxMjk2MjY5", Name: "gbrepo", Description: "This is a gb repo",
				OwnerInfo: OwnerInfo{Login: "gbuser", ID: 1, NodeID: "MDQ6VXNlcjE=", UserType: "User"}},
		},
	},
}
