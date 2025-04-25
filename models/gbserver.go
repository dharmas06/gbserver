package models

// import "slices"

// type GitRepoData struct {
// 	RepoName     string   `json:"repo"`
// 	BranchNames  []string `json:"branch"`
// 	PullRequests []PullRequest
// }

// type PullRequest struct {
// 	SourceBranch string `json:"source-branch"`
// 	TargetBranch string `json:"target-branch"`
// 	Title        string `json:"title"`
// 	Description  string `json:"description"`
// }

// var gitData = []*GitRepoData{
// 	{RepoName: "barerepo"},
// 	{RepoName: "devrepo", BranchNames: []string{"dev1"}},
// }

// func AddGitRepo(g *GitRepoData) bool {
// 	repoList := ListGitRepo()
// 	if !slices.Contains(repoList, g.RepoName) {
// 		gitData = append(gitData, g)
// 		return true
// 	}
// 	return false

// }

// func ListGitRepo() []string {
// 	var repoList []string
// 	for _, data := range gitData {
// 		repoList = append(repoList, data.RepoName)
// 	}
// 	return repoList

// }
