package service

import (
	"gbserver/models"
	"hash/fnv"
	"math/rand"
	"regexp"
	"slices"
	"strconv"
	"strings"
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
	// 	  "URL": "https://api.gbserver.com/Repos/admin/Hello-World/commits/c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
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
	//	{Ref: "refs/heads/gbbranch", NodeID: "MDM6UmVmcmVmcy9oZWFkcy9mZWF0dXJlQQ==", URL: "https://api.gbserver.com/Repos/gbUser/gbRepo/git/refs/heads/gbbranch",
	//
	//	"Object": {
	//			  "type": "Commit",
	//			  "SHA": "aa218f56b14c9653891f9e74264a383fa43fefbd",
	//			  "URL": "https://api.gbserver.com/Repos/octocat/Hello-World/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
	//			},
}

type baseHeadPRResponse struct {
	Ref  string
	SHA  string
	User OwnerInfo
	Repo string
}

type PRResponse struct {
	URL          string
	ID           string
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
	ListPRs(owner, repoName string) ([]PRResponse, error)                                       // get /Repos/{owner}/{Repo}/pulls
	CreatePR(owner, repoName string, cPRReq *PRRequest) (PRResponse, error)                     // post /Repos/{owner}/{Repo}/pulls
	UpdatePR(owner, repoName string, pull_number int, prRequest *PRRequest) (PRResponse, error) //patch /Repos/{owner}/{Repo}/pulls/{pull_number} State - closed

}

type GbService struct {
	GbStoreInstance *models.GbStore
}

func hasher(data string) string {
	h := fnv.New64a()
	h.Write([]byte(data))
	return strconv.FormatUint(h.Sum64(), 10)
}
func (g *GbService) ListRepos(orgName, ownerName string) ([]RepoResponse, error) {
	var outputResp []RepoResponse
	if _, exists := g.GbStoreInstance.Orgs[orgName]; !exists {
		return outputResp, ErrOrgNotFound
	}
	if _, exists := g.GbStoreInstance.Users[orgName+"/"+ownerName]; !exists {
		return outputResp, ErrOwnerNotFound
	}

	g.GbStoreInstance.MU.RLock()
	repoList := g.GbStoreInstance.Users[orgName+"/"+ownerName].Repos
	repoOwner := g.GbStoreInstance.Users[orgName+"/"+ownerName]

	g.GbStoreInstance.MU.RUnlock()
	ownerInfo := OwnerInfo{Login: repoOwner.LoginName, ID: repoOwner.ID, NodeID: repoOwner.NodeID, UserType: repoOwner.UserType}

	for _, repoInfo := range repoList {
		g.GbStoreInstance.MU.RLock()
		repoDetails := g.GbStoreInstance.Repos[orgName+"/"+ownerName+"/"+repoInfo]
		g.GbStoreInstance.MU.RUnlock()
		repoResponse := RepoResponse{ID: repoDetails.ID, Name: repoDetails.Name, Node_ID: repoDetails.Node_ID, Description: repoDetails.Description, OwnerInfo: ownerInfo}
		outputResp = append(outputResp, repoResponse)
	}
	return outputResp, nil
}

func generateCustomID(IDType string) string {
	var randomChar string
	var IDLen int
	if IDType == "NODEID" {
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

func (g *GbService) CreateRepo(orgName, ownerName string, RepoRequest *CreateRepoRequest) (RepoResponse, error) {

	g.GbStoreInstance.MU.RLock()
	if _, exists := g.GbStoreInstance.Orgs[orgName]; !exists {
		g.GbStoreInstance.MU.RUnlock()
		return RepoResponse{}, ErrOrgNotFound
	}
	if _, exists := g.GbStoreInstance.Users[orgName+"/"+ownerName]; !exists {
		g.GbStoreInstance.MU.RUnlock()
		return RepoResponse{}, ErrOwnerNotFound
	}

	repoList := g.GbStoreInstance.Users[orgName+"/"+ownerName].Repos
	if slices.Contains(repoList, RepoRequest.Name) {
		g.GbStoreInstance.MU.RUnlock()
		return RepoResponse{}, ErrRepoAlreadyExists
	}

	repoID := g.GbStoreInstance.Orgs[orgName].ReposCount + 1

	g.GbStoreInstance.MU.RUnlock()
	nodeID := generateCustomID("NODEID")
	g.GbStoreInstance.MU.Lock()

	g.GbStoreInstance.Repos[orgName+"/"+ownerName+"/"+RepoRequest.Name] = &models.Repository{ID: repoID, Name: RepoRequest.Name, Node_ID: nodeID,
		Description: RepoRequest.Description,
		OrgName:     orgName, UserName: ownerName, Branches: []string{}}

	g.GbStoreInstance.Users[orgName+"/"+ownerName].Repos = append(g.GbStoreInstance.Users[orgName+"/"+ownerName].Repos, RepoRequest.Name)
	g.GbStoreInstance.Orgs[orgName].Repos = append(g.GbStoreInstance.Orgs[orgName].Repos, RepoRequest.Name)
	g.GbStoreInstance.Orgs[orgName].ReposCount = repoID
	g.GbStoreInstance.MU.Unlock()
	g.GbStoreInstance.MU.RLock()

	resp := RepoResponse{ID: repoID, Name: RepoRequest.Name, Node_ID: nodeID, Description: RepoRequest.Description,
		OwnerInfo: OwnerInfo{Login: g.GbStoreInstance.Users[orgName+"/"+ownerName].LoginName,
			ID:       g.GbStoreInstance.Users[orgName+"/"+ownerName].ID,
			NodeID:   g.GbStoreInstance.Users[orgName+"/"+ownerName].NodeID,
			UserType: g.GbStoreInstance.Users[orgName+"/"+ownerName].UserType,
		}}
	g.GbStoreInstance.MU.RUnlock()

	return resp, nil
}

func removeElementByValue(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			slice = slices.Delete(slice, i, i+1)
		}
	}
	return slice

}
func (g *GbService) DeleteRepo(orgName, owner, repoName string) (bool, error) {

	err := g.validateOrgOwnerRepo(orgName, owner, repoName)
	if err != nil {
		return false, err
	}
	g.GbStoreInstance.MU.Lock()
	delete(g.GbStoreInstance.Repos, orgName+"/"+owner+"/"+repoName)
	g.GbStoreInstance.Users[orgName+"/"+owner].Repos = removeElementByValue(g.GbStoreInstance.Users[orgName+"/"+owner].Repos, repoName)
	g.GbStoreInstance.Orgs[orgName].Repos = removeElementByValue(g.GbStoreInstance.Orgs[orgName].Repos, repoName)
	g.GbStoreInstance.MU.Unlock()

	return true, nil
}

func (g *GbService) ListBranches(orgName, owner, repoName string) ([]ListBranchresponse, error) {

	var outputResp []ListBranchresponse
	err := g.validateOrgOwnerRepo(orgName, owner, repoName)
	if err != nil {
		return outputResp, err
	}
	g.GbStoreInstance.MU.RLock()
	branchList := g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].Branches
	g.GbStoreInstance.MU.RUnlock()
	if len(branchList) == 0 {
		return outputResp, ErrBranchesNotFound
	}
	for _, branchName := range branchList {
		g.GbStoreInstance.MU.RLock()

		branchData := g.GbStoreInstance.Branches[orgName+"/"+owner+"/"+repoName+"/"+branchName]
		g.GbStoreInstance.MU.RUnlock()
		commitDetails := CommitDetails{SHA: branchData.CommitInfo.SHA, URL: branchData.CommitInfo.URL}

		branchresponse := ListBranchresponse{Name: branchName, Commit: commitDetails, Protected: branchData.Protected}
		outputResp = append(outputResp, branchresponse)
	}

	return outputResp, nil
}

func (g *GbService) CreateBranch(orgName, owner, repoName string, cbreq *CreateBranchRequest) (CreateBranchResponse, error) {

	var createBranchResp CreateBranchResponse

	err := g.validateOrgOwnerRepo(orgName, owner, repoName)
	if err != nil {
		return createBranchResp, err
	}
	re := regexp.MustCompile(`refs/heads/([^/]+)$`)
	matches := re.FindStringSubmatch(cbreq.Ref)
	var branch string
	if len(matches) >= 2 {
		branch = matches[1]
		//	fmt.Println(branch)
	} else {
		//	fmt.Println("Invalid branch name found in refs")
		return CreateBranchResponse{}, ErrInvalidBranchName
	}
	g.GbStoreInstance.MU.RLock()
	branchList := g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].Branches
	branchID := len(g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].Branches) + 1
	g.GbStoreInstance.MU.RUnlock()
	if slices.Contains(branchList, branch) {
		return createBranchResp, ErrBranchesAlreadyExists
	}

	nodeID := generateCustomID("NODEID")
	url := "https://api.gbserver.com/repos/" + owner + "/" + repoName + "/git/commits/" + cbreq.SHA
	commit := models.CommitDetails{SHA: cbreq.SHA, URL: url}
	fullBranchName := orgName + "/" + owner + "/" + repoName + "/" + branch

	g.GbStoreInstance.MU.Lock()
	g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].Branches = append(g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].Branches, branch)

	g.GbStoreInstance.Branches[fullBranchName] = &models.Branch{
		ID:         branchID,
		RepoName:   repoName,
		Name:       branch,
		NodeID:     nodeID,
		URL:        url,
		Protected:  true,
		CommitInfo: commit,
	}
	g.GbStoreInstance.MU.Unlock()
	createBranchResp = CreateBranchResponse{Ref: cbreq.Ref, NodeID: nodeID, URL: url,
		Object: CreateBranchObjectResponse{Type: "commit", SHA: cbreq.SHA, URL: url}}
	return createBranchResp, nil
}

func (g *GbService) DeleteBranch(orgName, owner, repoName, branch string) (bool, error) {

	err := g.validateOrgOwnerRepo(orgName, owner, repoName)
	fullBranchName := orgName + "/" + owner + "/" + repoName + "/" + branch
	if err != nil {
		return false, err
	}
	g.GbStoreInstance.MU.RLock()
	if _, branchExists := g.GbStoreInstance.Branches[fullBranchName]; !branchExists {
		g.GbStoreInstance.MU.RUnlock()
		return false, ErrBranchesNotFound
	}
	g.GbStoreInstance.MU.RUnlock()
	g.GbStoreInstance.MU.Lock()
	if g.GbStoreInstance.Branches[fullBranchName].PullRequestID != "" {

		delete(g.GbStoreInstance.PullRequests, g.GbStoreInstance.Branches[fullBranchName].PullRequestID)
		g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].PrIDs = removeElementByValue(g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].PrIDs, g.GbStoreInstance.Branches[fullBranchName].PullRequestID)
	}
	delete(g.GbStoreInstance.Branches, fullBranchName)

	g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].Branches = removeElementByValue(g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].Branches, branch)
	g.GbStoreInstance.MU.Unlock()
	//fmt.Println("After delete branch", g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName])
	return true, nil
}

func (g *GbService) validateOrgOwnerRepo(orgName, owner, repoName string) error {
	g.GbStoreInstance.MU.RLock()
	defer g.GbStoreInstance.MU.RUnlock()
	if _, exists := g.GbStoreInstance.Orgs[orgName]; !exists {
		return ErrOrgNotFound
	}
	if _, exists := g.GbStoreInstance.Users[orgName+"/"+owner]; !exists {
		return ErrOwnerNotFound
	}
	if _, exists := g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName]; !exists {
		return ErrRepoNotFound
	}

	return nil
}

func (g *GbService) ListPRs(orgName, owner, repoName string) ([]PRResponse, error) {

	var listPRresponse []PRResponse

	err := g.validateOrgOwnerRepo(orgName, owner, repoName)
	if err != nil {
		return listPRresponse, err
	}
	g.GbStoreInstance.MU.RLock()
	//branchList := g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].Branches
	prsList := g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].PrIDs
	g.GbStoreInstance.MU.RUnlock()
	for _, prID := range prsList {
		//	prID := g.GbStoreInstance.Branches[branchName].PullRequestID

		//	fmt.Println("PRS...", g.GbStoreInstance.PullRequests[prID])
		g.GbStoreInstance.MU.RLock()
		prDetails := g.GbStoreInstance.PullRequests[prID]
		g.GbStoreInstance.MU.RUnlock()

		featureBranch := strings.Split(prDetails.FromBranch, ":")
		//featureBranchUser := featureBranch[0]
		featureBranchName := featureBranch[1]
		fullfeatureBranchName := orgName + "/" + owner + "/" + repoName + "/" + featureBranchName
		fullbaseBranchName := orgName + "/" + owner + "/" + repoName + "/" + prDetails.ToBranch

		g.GbStoreInstance.MU.RLock()
		prOwner := OwnerInfo{Login: g.GbStoreInstance.Users[orgName+"/"+owner].LoginName,
			NodeID: g.GbStoreInstance.Users[orgName+"/"+owner].NodeID, ID: g.GbStoreInstance.Users[orgName+"/"+owner].ID,
			UserType: g.GbStoreInstance.Users[orgName+"/"+owner].UserType}

		prHeadResp := baseHeadPRResponse{
			Ref: g.GbStoreInstance.Branches[fullfeatureBranchName].Name, SHA: g.GbStoreInstance.Branches[fullfeatureBranchName].CommitInfo.SHA,
			User: prOwner, Repo: repoName}
		prBaseResp := baseHeadPRResponse{
			Ref: g.GbStoreInstance.Branches[fullbaseBranchName].Name, SHA: g.GbStoreInstance.Branches[fullbaseBranchName].CommitInfo.SHA,
			User: prOwner, Repo: repoName}
		g.GbStoreInstance.MU.RUnlock()

		prResp := PRResponse{
			URL: prDetails.URL, ID: prDetails.ID, NodeID: prDetails.NodeID,
			Title: prDetails.Title, Body: prDetails.Body, State: prDetails.State,
			User: prOwner, Commits: prDetails.Commits, Additions: prDetails.Additions, Deletions: prDetails.Deletions,
			ChangedFiles: prDetails.ChangedFiles,
			Head:         prHeadResp,
			Base:         prBaseResp,
		}
		listPRresponse = append(listPRresponse, prResp)

	}
	return listPRresponse, nil
}

// // patch /Repos/{owner}/{Repo}/pulls/{pull_number} State - closed
func (g *GbService) UpdatePR(orgName, owner, repoName, pull_number string, prRequest *PRRequest) (PRResponse, error) {
	//'{"Title":"new Title","Body":"updated Body","State":"open","base":"master"}'
	//closePRRequest := PRRequest{State: "closed"}
	var closedPR PRResponse

	err := g.validateOrgOwnerRepo(orgName, owner, repoName)
	if err != nil {
		return closedPR, err
	}
	//pullNumberStr := strconv.Itoa(pull_number)
	//fullPRID := orgName + "/" + owner + "/" + repoName + pullNumberStr
	g.GbStoreInstance.MU.RLock()
	if !slices.Contains(g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].PrIDs, pull_number) {
		g.GbStoreInstance.MU.RUnlock()
		return closedPR, ErrPRNotFound
	}
	prDetails := g.GbStoreInstance.PullRequests[pull_number]
	if g.GbStoreInstance.PullRequests[pull_number].State == "closed" {
		g.GbStoreInstance.MU.RUnlock()
		return closedPR, ErrPRAlreadyClosed
	}
	g.GbStoreInstance.MU.RUnlock()
	branchInfo := prDetails.FromBranch
	branchName := strings.Split(branchInfo, ":")[1]
	fullBranchName := orgName + "/" + owner + "/" + repoName + "/" + branchName
	g.GbStoreInstance.MU.Lock()
	g.GbStoreInstance.Branches[fullBranchName].PullRequestID = ""
	//	g.GbStoreInstance.PullRequests[pull_number].Title = prRequest.Title
	//	g.GbStoreInstance.PullRequests[pull_number].Body = prRequest.Body
	g.GbStoreInstance.PullRequests[pull_number].State = prRequest.State
	closedPRData := g.GbStoreInstance.PullRequests[pull_number]
	g.GbStoreInstance.MU.Unlock()

	featureBranch := strings.Split(closedPRData.FromBranch, ":")

	featureBranchName := featureBranch[1]
	fullFeatureBranchName := orgName + "/" + owner + "/" + repoName + "/" + featureBranchName
	fullBaseBranchName := orgName + "/" + owner + "/" + repoName + "/" + closedPRData.ToBranch

	g.GbStoreInstance.MU.RLock()
	ownerDetails := OwnerInfo{
		Login:    g.GbStoreInstance.Users[orgName+"/"+owner].LoginName,
		ID:       g.GbStoreInstance.Users[orgName+"/"+owner].ID,
		UserType: g.GbStoreInstance.Users[orgName+"/"+owner].UserType,
		NodeID:   g.GbStoreInstance.Users[orgName+"/"+owner].NodeID,
	}

	prHeadResp := baseHeadPRResponse{

		Ref:  closedPRData.FromBranch,
		SHA:  g.GbStoreInstance.Branches[fullFeatureBranchName].CommitInfo.SHA,
		User: ownerDetails,
		Repo: repoName,
	}

	prBaseResp := baseHeadPRResponse{

		Ref:  closedPRData.ToBranch,
		SHA:  g.GbStoreInstance.Branches[fullBaseBranchName].CommitInfo.SHA,
		User: ownerDetails,
		Repo: repoName,
	}
	g.GbStoreInstance.MU.RUnlock()

	closedPR = PRResponse{
		URL:          closedPRData.URL,
		ID:           closedPRData.ID,
		NodeID:       closedPRData.NodeID,
		Title:        closedPRData.Title,
		Body:         closedPRData.Body,
		State:        closedPRData.State,
		User:         ownerDetails,
		Commits:      closedPRData.Commits,
		Additions:    closedPRData.Additions,
		Deletions:    closedPRData.Deletions,
		ChangedFiles: closedPRData.ChangedFiles,
		Head:         prHeadResp,
		Base:         prBaseResp,
	}

	return closedPR, nil
}

// // post /Repos/{owner}/{Repo}/pulls
func (g *GbService) CreatePR(orgName, owner, repoName string, cPRReq *PRRequest) (PRResponse, error) {
	//'{"Title":"Amazing new feature",
	//"Body":"Please pull these awesome changes in!","head":"octocat:new-feature","base":"master"}'
	var createPRresponse PRResponse

	err := g.validateOrgOwnerRepo(orgName, owner, repoName)
	if err != nil {
		return createPRresponse, err
	}

	featureBranch := strings.Split(cPRReq.Head, ":")
	//featureBranchUser := featureBranch[0]
	featureBranchName := featureBranch[1]
	fullFeatureBranchName := orgName + "/" + owner + "/" + repoName + "/" + featureBranchName
	fullBaseBranchName := orgName + "/" + owner + "/" + repoName + "/" + cPRReq.Base

	g.GbStoreInstance.MU.RLock()
	if _, branchExists := g.GbStoreInstance.Branches[fullFeatureBranchName]; !branchExists {
		g.GbStoreInstance.MU.RUnlock()
		return createPRresponse, ErrBranchesNotFound
	}
	if _, branchExists := g.GbStoreInstance.Branches[fullBaseBranchName]; !branchExists {
		g.GbStoreInstance.MU.RUnlock()
		return createPRresponse, ErrBranchesNotFound
	}

	if g.GbStoreInstance.Branches[fullFeatureBranchName].PullRequestID != "" {
		g.GbStoreInstance.MU.RUnlock()
		return createPRresponse, ErrPRAlreadyExists
	}

	prCount := g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].TotalPRs + 1
	g.GbStoreInstance.MU.RUnlock()

	fullPRName := orgName + "/" + owner + "/" + repoName + "/" + strconv.Itoa(prCount)
	prID := hasher(fullPRName)

	g.GbStoreInstance.MU.Lock()
	g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].TotalPRs = prCount

	g.GbStoreInstance.Branches[fullFeatureBranchName].PullRequestID = prID

	g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].PrIDs = append(g.GbStoreInstance.Repos[orgName+"/"+owner+"/"+repoName].PrIDs, prID)

	g.GbStoreInstance.MU.Unlock()
	nodeId := generateCustomID("NODEID")
	//prIDAsString := strconv.Itoa(prID)

	commits := rand.Intn(50)
	additions := rand.Intn(50)
	deletions := rand.Intn(50)
	changedFiles := rand.Intn(50)

	url := "https://api.gbserver.com/repos/" + owner + "/" + repoName + "/" + prID

	g.GbStoreInstance.MU.RLock()
	g.GbStoreInstance.PullRequests[prID] = &models.PullRequest{
		NodeID:       nodeId,
		URL:          url,
		ID:           prID,
		RepoName:     repoName,
		FromBranch:   cPRReq.Head,
		ToBranch:     cPRReq.Base,
		AuthorID:     g.GbStoreInstance.Users[orgName+"/"+owner].ID,
		State:        "open",
		Title:        cPRReq.Title,
		Body:         cPRReq.Body,
		Commits:      commits,
		Additions:    additions,
		Deletions:    deletions,
		ChangedFiles: changedFiles,
	}

	ownerDetails := OwnerInfo{
		Login:    g.GbStoreInstance.Users[orgName+"/"+owner].LoginName,
		ID:       g.GbStoreInstance.Users[orgName+"/"+owner].ID,
		UserType: g.GbStoreInstance.Users[orgName+"/"+owner].UserType,
		NodeID:   g.GbStoreInstance.Users[orgName+"/"+owner].NodeID,
	}

	prHeadResp := baseHeadPRResponse{

		Ref:  cPRReq.Base,
		SHA:  g.GbStoreInstance.Branches[fullBaseBranchName].CommitInfo.SHA,
		User: ownerDetails,
		Repo: repoName,
	}

	prBaseResp := baseHeadPRResponse{

		Ref:  featureBranchName,
		SHA:  g.GbStoreInstance.Branches[fullFeatureBranchName].CommitInfo.SHA,
		User: ownerDetails,
		Repo: repoName,
	}
	g.GbStoreInstance.MU.RUnlock()
	//fmt.Println(g.GbStoreInstance.PullRequests)
	createPRresponse = PRResponse{
		URL:          url,
		ID:           prID,
		NodeID:       nodeId,
		Title:        cPRReq.Title,
		Body:         cPRReq.Body,
		State:        "open",
		User:         ownerDetails,
		Commits:      commits,
		Additions:    additions,
		Deletions:    deletions,
		ChangedFiles: changedFiles,
		Head:         prHeadResp,
		Base:         prBaseResp,
	}

	return createPRresponse, nil

}
