package service

import (
	"fmt"
	"testing"
)

var repoRequest = &CreateRepoRequest{
	"test",
	"sample test repo",
}

var orgname = "gb"

var gbService = Gbserver{}

func TestCreateRepo(t *testing.T) {

	fmt.Println(gbService.CreateRepo(orgname, repoRequest))
}

func TestListRepos(t *testing.T) {
	fmt.Println(gbService.ListRepos(orgname))
}

func TestDeleteRepo(t *testing.T) {
	fmt.Println(gbService.DeleteRepo("gbuser", "gbrepo"))
}

func TestListBranches(t *testing.T) {
	fmt.Println(gbService.ListBranches("gbuser", "gbrepo"))
	fmt.Println(gbService.ListBranches("testuser", "testrepo"))
}

func TestCreateBranch(t *testing.T) {
	cbReq := CreateBranchRequest{Ref: "refs/heads/featureA", SHA: "aa218f56b14c9653891f9e74264a383fa43fefbd"}
	fmt.Println(gbService.CreateBranch("gbuser", "gbrepo", &cbReq))
	//fmt.Println(gbService.ListBranches("testuser", "testrepo"))
}

func TestDeleteBranch(t *testing.T) {
	fmt.Println(gbService.DeleteBranch("gbuser", "gbrepo", "refs/heads/gbbranch"))
	//fmt.Println(gbService.ListBranches("testuser", "testrepo"))
}

func TestListPRs(t *testing.T) {
	fmt.Println(gbService.ListPRs("gbuser", "gbrepo"))

}

func TestUpdatePR(t *testing.T) {
	closePRRequest := PRRequest{State: "closed"}
	fmt.Println(gbService.UpdatePR("gbuser", "gbrepo", 1347, &closePRRequest))
}

func TestCreatePR(t *testing.T) {
	createPRRequest := PRRequest{Title: "FeatureB", Body: "Add this B feature", Head: "gbbranch:new-feature", Base: "master", State: "open"}
	fmt.Println(gbService.CreatePR("gbuser", "gbrepo", &createPRRequest))
	fmt.Println(gbService.ListPRs("gbuser", "gbrepo"))
}
