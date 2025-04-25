package service

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type OwnerInfo struct {
	Login    string `json:"login"`
	ID       int    `json:"id"`
	NodeID   string `json:"nodeId"`
	UserType string `json:"Type"`
}

type RepoResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Node_ID     string `json:"node_id"`
	Description string `json:"description"`
	OwnerInfo   OwnerInfo
}

type CreateRepoRequest struct {
	Name        string
	Description string
	//{"Name":"Hello-World","Description":"This is your first Repository",
	//"homepage":"https://github.com","private":false,"has_issues":true,"has_projects":true,"has_wiki":true}'
}

type CommitDetails struct {
	SHA string
	URL string
}

type ListBranchresponse struct {
	Name      string
	Commit    CommitDetails
	Protected bool
	// {
	// 	"Name": "master",
	// 	"Commit": {
	// 	  "SHA": "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc",
	// 	  "URL": "https://api.github.com/Repos/admin/Hello-World/commits/c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	// 	},
	// 	"Protected": true,
}

type CreateBranchRequest struct {
	Ref string
	SHA string
	//{"Ref":"refs/heads/featureA","SHA":"aa218f56b14c9653891f9e74264a383fa43fefbd"}'
}

type CreateBranchObjectResponse struct {
	Type string
	SHA  string
	URL  string
}

type CreateBranchResponse struct {
	Ref    string
	NodeID string
	URL    string
	Object CreateBranchObjectResponse
	//	{Ref: "refs/heads/gbbranch", NodeID: "MDM6UmVmcmVmcy9oZWFkcy9mZWF0dXJlQQ==", URL: "https://api.github.com/Repos/gbUser/gbRepo/git/refs/heads/gbbranch",
	//
	//	"Object": {
	//			  "type": "Commit",
	//			  "SHA": "aa218f56b14c9653891f9e74264a383fa43fefbd",
	//			  "URL": "https://api.github.com/Repos/octocat/Hello-World/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
	//			},
}

type baseHeadPRResponse struct {
	Label string
	Ref   string
	SHA   string
	User  OwnerInfo
	Repo  RepoResponse
}

type PRResponse struct {
	URL          string
	ID           int
	NodeID       string
	Title        string
	Body         string
	State        string
	User         OwnerInfo
	Commits      int
	Additions    int
	Deletions    int
	ChangedFiles int
	Head         baseHeadPRResponse
	Base         baseHeadPRResponse
}

type PRRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
	State string `json:"state"`
	// '{"Title":"Amazing new feature",
	// "Body":"Please pull these awesome changes in!","head":"admin:new-feature","base":"master"}'

	//close PR
	//'{"Title":"new Title","Body":"updated Body","State":"open","base":"master"}'
}

type Gbservice interface {
	ListRepos(orgName string) []RepoResponse                                         //get  /orgs/{org}/Repos
	CreateRepo(orgName string, RepoRequest *CreateRepoRequest) (RepoResponse, error) // post   /orgs/{org}/Repos
	DeleteRepo(owner, repoName string) (bool, error)                                 //delete /Repos/{owner}/{Repo}

	// branches
	ListBranches(owner, repoName string) ([]ListBranchresponse, error)                             // get /Repos/{owner}/{Repo}/branches
	CreateBranch(owner, repoName string, cbReq *CreateBranchRequest) (CreateBranchResponse, error) // post /Repos/{owner}/{Repo}/git/refs
	DeleteBranch(owner, repoName, Ref string) (bool, error)                                        //delete /Repos/{owner}/{Repo}/git/refs/{Ref}

	// Pull request
	ListPRs(owner, repoName string) ([]PRResponse, error)                                      // get /Repos/{owner}/{Repo}/pulls
	CreatePR(owner, repoName string, cPRReq *PRRequest) (PRResponse, error)                    // post /Repos/{owner}/{Repo}/pulls
	ClosePR(owner, repoName string, pull_number int, prRequest *PRRequest) (PRResponse, error) //patch /Repos/{owner}/{Repo}/pulls/{pull_number} State - closed

}

type Gbserver struct {
	muRW sync.RWMutex
}

func (g *Gbserver) ListRepos(orgName string) []RepoResponse {
	owners := getOrgMembers(orgName)
	var Repos []RepoResponse
	for _, repoData := range gbServerRepoData {
		for _, owner := range owners {
			if owner.Login == repoData.OwnerInfo.Login { // == owner && data.OwnerInfo.ID == ID {
				//fmt.Println(repoData)
				Repos = append(Repos, repoData)
			}
		}
	}
	return Repos
}

func validateOrg(orgName string) bool {
	for org := range gbServerOrg {
		if orgName == org {
			return true
		}
	}
	return false
}

func generateID() int {
	low := 10000
	hi := 99999
	return low + rand.Intn(hi-low)
}

func generateCustomID(IDType string) string {
	var randomChar string
	var IDLen int
	if IDType == "NodeID" {
		randomChar = "aAbBcC1d2e3D4fe5gE6hF7ij8Gk9HlI0JmKnX=o=LWYMNpZqOrPsVQtuvRwSxTyUz"
		IDLen = 23
	}
	if IDType == "SHA" {
		randomChar = "abcdefghijklmnopqrstuvwxyz0123456789"
		IDLen = 40
	}
	var NodeID string
	length := len(randomChar)
	for k := range IDLen {
		NodeID += string(randomChar[rand.Intn(length-k)])
	}
	return NodeID
}

func getOrgMembers(orgName string) []OwnerInfo {
	for org, members := range gbServerOrg {
		if orgName == org {
			return members
		}
	}
	return nil
}

func getIDFromOwner(owner string) (int, error) {

	for _, ownerDetails := range gbServerOwners {
		if ownerDetails.Login == owner {
			return ownerDetails.ID, nil
		}
	}
	return 0, ErrOwnerNotFound
}

func (g *Gbserver) CreateRepo(orgName string, RepoRequest *CreateRepoRequest) (RepoResponse, error) {

	if !validateOrg(orgName) {
		return RepoResponse{}, ErrOrgNotFound
	}
	var err error
	repoList := g.ListRepos(orgName)
	for _, repo := range repoList {
		if repo.Name == RepoRequest.Name {
			err = errors.New("repo already exist")
			return RepoResponse{}, err
		}
	}
	//owner, ID := getUserNameAndID(orgName)
	var owner = gbServerOrg[orgName][0].Login
	var ID = gbServerOrg[orgName][0].ID
	resp := RepoResponse{ID: generateID(),
		Name:        RepoRequest.Name,
		Node_ID:     generateCustomID("NodeID"),
		Description: RepoRequest.Description,
		OwnerInfo:   OwnerInfo{Login: owner, ID: ID, NodeID: "MDQ6VXNlcjE=", UserType: "User"}}
	g.muRW.Lock()
	gbServerRepoData = append(gbServerRepoData, resp)
	g.muRW.Unlock()
	return resp, err
}

func (g *Gbserver) DeleteRepo(owner, repoName string) (bool, error) {

	ID, err := getIDFromOwner(owner)
	if err != nil {
		fmt.Println("Error occurred.", err)
		return false, err
	}
	resp := make([]RepoResponse, len(gbServerRepoData))
	copy(resp, gbServerRepoData)
	for index, data := range resp {
		if data.OwnerInfo.Login == owner && data.OwnerInfo.ID == ID && data.Name == repoName {
			gbServerRepoData = append(gbServerRepoData[:index], gbServerRepoData[index+1:]...)
			fmt.Println(gbServerRepoData)
			return true, nil
		}
	}
	return false, ErrRepoNotFound
}

func (g *Gbserver) validateOwnerAndRepo(owner, repoName string) error {
	_, err := getIDFromOwner(owner)
	if err != nil {
		return ErrOwnerNotFound
	}
	repoList := g.ListRepos(orgName)
	var repoFound bool
	for _, repo := range repoList {
		if repo.Name == repoName {
			repoFound = true
		}
	}
	if !repoFound {
		return ErrRepoNotFound
	}
	return nil
}

func (g *Gbserver) ListBranches(owner, repoName string) ([]ListBranchresponse, error) {
	var found bool
	var output []ListBranchresponse

	for _, data := range gbSeverBranchListData {
		input := data.Commit.URL
		re := regexp.MustCompile(`repos/([^/]+)/([^/]+)/commits`)
		matches := re.FindStringSubmatch(input)

		if len(matches) >= 3 {
			User := matches[1]
			Repo := matches[2]
			//	fmt.Println("Organization & Repo name..", User, Repo)

			err := g.validateOwnerAndRepo(owner, repoName)
			if err != nil {
				return output, err
			}

			if User == owner && Repo == repoName {
				output = append(output, data)
				found = true
			}
		}
	}
	if found {
		return output, nil
	}
	return output, ErrBranchesNotFound
}

func (g *Gbserver) CreateBranch(owner, repoName string, cbreq *CreateBranchRequest) (CreateBranchResponse, error) {
	err := g.validateOwnerAndRepo(owner, repoName)
	if err != nil {
		return CreateBranchResponse{}, err
	}
	RefData := cbreq.Ref
	re := regexp.MustCompile(`refs/heads/([^/]+)$`)
	matches := re.FindStringSubmatch(RefData)
	var branch string
	if len(matches) >= 2 {
		branch = matches[1]
		fmt.Println(branch)
	} else {
		fmt.Println("Invalid branch name found in refs")
		return CreateBranchResponse{}, errors.New("invalid branch name. Specify as refs/head/<branch>")
	}
	branchList, _ := g.ListBranches(owner, repoName)
	for _, branchData := range branchList {
		if branchData.Name == branch {
			return CreateBranchResponse{}, ErrBranchesAlreadyExists
		}
	}

	branchInput := CreateBranchResponse{Ref: "refs/heads/" + branch,
		NodeID: generateCustomID("NodeID"), URL: "https://api.github.com/repos/" + owner + "/" + repoName + "/git/" + cbreq.Ref,
		Object: CreateBranchObjectResponse{Type: "Commit", SHA: cbreq.SHA,
			URL: "https://api.github.com/repos/" + owner + "/" + repoName + "/git/commits/" + cbreq.SHA},
	}

	listbranchInput := ListBranchresponse{Name: branch, Commit: CommitDetails{SHA: cbreq.SHA, URL: "https://api.github.com/repos/" + owner + "/" + repoName + "/commits/" + cbreq.SHA}, Protected: true}
	g.muRW.Lock()
	defer g.muRW.Unlock()
	gbServerBranchData = append(gbServerBranchData, branchInput)
	gbSeverBranchListData = append(gbSeverBranchListData, listbranchInput)

	return branchInput, nil

}

func (g *Gbserver) DeleteBranch(owner, repoName, ref string) (bool, error) {
	branch := ref
	err := g.validateOwnerAndRepo(owner, repoName)
	if err != nil {
		return false, err
	}
	branchList, _ := g.ListBranches(owner, repoName)
	var repoFound bool
	for _, branchData := range branchList {
		if branchData.Name == branch {
			repoFound = true
		}
	}
	if !repoFound {
		fmt.Println("branch not found")
		return false, errors.New("branch not found")
	}
	var done, done1 bool
	var branchIndex, branchListIndex int
	temp := make([]CreateBranchResponse, len(gbServerBranchData))
	copy(temp, gbServerBranchData)
	for index, data := range temp {
		if data.Ref == "refs/heads/"+branch && strings.Contains(data.URL, "repos/"+owner+"/"+repoName+"/git/refs/heads/"+branch) {
			branchIndex = index
			done = true
		}
	}
	if done {
		copyData1 := make([]ListBranchresponse, len(gbSeverBranchListData))
		copy(copyData1, gbSeverBranchListData)
		for index, data := range copyData1 {
			if data.Name == branch {
				branchListIndex = index
				done1 = true
			}
		}
	}
	if done && done1 {
		gbServerBranchData = slices.Delete(gbServerBranchData, branchIndex, branchIndex+1)
		gbSeverBranchListData = slices.Delete(gbSeverBranchListData, branchListIndex, branchListIndex+1)
		return true, nil
	}
	return false, ErrBranchesNotFound
}

func (g *Gbserver) ListPRs(owner, repoName string) ([]PRResponse, error) {
	err := g.validateOwnerAndRepo(owner, repoName)
	if err != nil {
		return []PRResponse{}, err
	}
	var listPRresponse []PRResponse
	//	fmt.Println("Total Prs..", gbServerListPR)
	for _, data := range gbServerListPR {
		if data.User.Login == owner {
			re := regexp.MustCompile(`repos/([^/]+)/([^/]+)/pulls`)
			matches := re.FindStringSubmatch(data.URL)
			if len(matches) >= 3 {
				User := matches[1]
				Repo := matches[2]
				_, err := getIDFromOwner(owner)
				if err != nil {
					return listPRresponse, ErrOwnerNotFound
				}

				repoList := g.ListRepos(orgName)
				var repoFound bool
				for _, repo := range repoList {
					if repo.Name == repoName {
						repoFound = true
					}
				}
				if !repoFound {
					return listPRresponse, ErrRepoNotFound
				}

				if User == owner && Repo == repoName { // && data.State == "open" {
					listPRresponse = append(listPRresponse, data)
				}
			}
		}
	}
	return listPRresponse, nil
}

// patch /Repos/{owner}/{Repo}/pulls/{pull_number} State - closed
func (g *Gbserver) UpdatePR(owner, repoName string, pull_number int, prRequest *PRRequest) (PRResponse, error) {
	//'{"Title":"new Title","Body":"updated Body","State":"open","base":"master"}'
	//closePRRequest := PRRequest{State: "closed"}
	var closedPR PRResponse
	err := g.validateOwnerAndRepo(owner, repoName)
	if err != nil {
		return closedPR, err
	}
	for index, data := range gbServerListPR {
		re := regexp.MustCompile(`repos/([^/]+)/([^/]+)/pulls/([^/]+)`)
		matches := re.FindStringSubmatch(data.URL)
		if len(matches) >= 3 {
			User := matches[1]
			Repo := matches[2]
			pullNumber, err := strconv.Atoi(matches[3])
			if err != nil {
				fmt.Println("Error occured while converting the pull number")
			}
			if pullNumber == pull_number {
				if User == owner && Repo == repoName && data.State == "open" {
					g.muRW.Lock()
					defer g.muRW.Unlock()
					gbServerListPR[index].State = prRequest.State
					closedPR = data
					//fmt.Println("Updated data..", data)
					return closedPR, nil
				}
			}
		}
	}

	//	fmt.Println("*****No PRs found.", pull_number)
	return closedPR, ErrPRNotFound
}

func getOwnerInfoByOwner(Name string) OwnerInfo {
	for _, data := range gbServerOwners {
		if data.Login == Name {
			return data
		}
	}
	return OwnerInfo{}
}

func getRepoDetailsByRepoOwner(owner string, Name string) RepoResponse {
	for _, data := range gbServerRepoData {
		if data.OwnerInfo.Login == owner && data.Name == Name {
			return data
		}
	}
	return RepoResponse{}
}

func getSHAFromBranchDataByBranch(branch string) string {
	for _, data := range gbSeverBranchListData {
		if data.Name == branch {
			return data.Commit.SHA
		}
	}
	return ""

}

func (g *Gbserver) validateRepoOwnerOrgAccess(owner, headOwnerName string) error {
	var orgMembers []string
	for _, members := range gbServerOrg {
		for _, owners := range members {
			orgMembers = append(orgMembers, owners.Login)
		}
		if slices.Contains(orgMembers, owner) && slices.Contains(orgMembers, headOwnerName) {
			return nil
		}
	}
	return ErrOwnerNotInSameOrg
}

func (g *Gbserver) validateBranches(owner, repoName, branch string) error {
	branchList, err := g.ListBranches(owner, repoName)
	if err != nil {
		fmt.Println("Error occurred while validating branch access", err)
		return err
	}
	for _, branchData := range branchList {
		if branchData.Name == branch {
			return nil
		}
	}
	fmt.Println("branch not found")
	return ErrBranchesNotFound
}

// post /Repos/{owner}/{Repo}/pulls
func (g *Gbserver) CreatePR(owner, repoName string, cPRReq *PRRequest) (PRResponse, error) {
	//'{"Title":"Amazing new feature",
	//"Body":"Please pull these awesome changes in!","head":"octocat:new-feature","base":"master"}'
	UserDetails := getOwnerInfoByOwner(owner)
	RepoDetails := getRepoDetailsByRepoOwner(owner, repoName)

	// To build and store into list PR
	RefHeadBranchData := strings.Split(cPRReq.Head, ":")

	headOwner := RefHeadBranchData[0]
	err := g.validateRepoOwnerOrgAccess(owner, headOwner)
	if err != nil {
		return PRResponse{}, err
	}
	err = g.validateBranches(owner, repoName, RefHeadBranchData[1])
	if err != nil {
		return PRResponse{}, err
	}
	err = g.validateBranches(owner, repoName, cPRReq.Base)
	if err != nil {
		return PRResponse{}, err
	}
	SHAHeadData := getSHAFromBranchDataByBranch(RefHeadBranchData[1])
	SHABaseData := getSHAFromBranchDataByBranch(cPRReq.Base)

	prList, err := g.ListPRs(owner, repoName)
	if err != nil {
		//	fmt.Println("Error occurred.", err)
		return PRResponse{}, err
	}
	for _, pr := range prList {
		if pr.Head.Label == cPRReq.Head && pr.State == "open" {
			//	fmt.Println("PR already exists.Error occurred.", err)
			return PRResponse{}, ErrPRAlreadyExists
		}
	}
	prData := PRResponse{ID: len(gbServerListPR) + 1, NodeID: generateCustomID("NodeID"), URL: "https://api.github.com/repos/" + owner + "/" + repoName + "/pulls/" + strconv.Itoa(len(gbServerListPR)+1),
		Title: cPRReq.Title, Body: cPRReq.Body, State: "open", User: UserDetails, Commits: 10, Additions: 23, Deletions: 5, ChangedFiles: 8,
		Head: baseHeadPRResponse{Label: cPRReq.Head, Ref: RefHeadBranchData[1], SHA: SHAHeadData, User: UserDetails, Repo: RepoDetails},
		Base: baseHeadPRResponse{Label: cPRReq.Base, Ref: cPRReq.Base, SHA: SHABaseData, User: UserDetails, Repo: RepoDetails},
	}
	//	fmt.Println("before cr creation", gbServerListPR)
	g.muRW.Lock()
	gbServerListPR = append(gbServerListPR, prData)
	//	fmt.Println("After cr creation", gbServerListPR)
	g.muRW.Unlock()
	return prData, nil
}
