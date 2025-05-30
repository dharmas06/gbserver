package handlers

import (
	"encoding/json"
	"fmt"
	"gbserver/models"
	"gbserver/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type GitRepo struct {
	l         *log.Logger
	gbService service.GbService
}

func NewGitRepo(l *log.Logger) *GitRepo {
	return &GitRepo{l, service.GbService{GbStoreInstance: models.NewGbStore()}}
}

func (g *GitRepo) ListRepoHandler(rw http.ResponseWriter, r *http.Request) {

	g.l.Println("Processing Get request..List Repo handler")
	vars := mux.Vars(r)
	orgName := vars["org"]
	ownerName := vars["owner"]
	//	g.l.Println("Organization & owner name..", orgName, ownerName)

	repoList, err := g.gbService.ListRepos(orgName, ownerName)
	if err != nil {
		if err == service.ErrOwnerNotFound || err == service.ErrOrgNotFound {
			g.l.Println("Error occurred while fetching the repo list.", err)
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		g.l.Println("Error occurred while fetching the repo list.", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	//g.l.Println("Retrieved Repo list.", repoList)
	g.l.Println("Retrieved Repo list.")
	rw.Header().Set("Content-Type", "Application/json")

	err = json.NewEncoder(rw).Encode(repoList)
	if err != nil {
		g.l.Println("Error occured while decoding the Git repo list", err)
		http.Error(rw, "Error occured while decoding the output", http.StatusInternalServerError)
		return
	}
}

func (g *GitRepo) CreateRepoHandler(rw http.ResponseWriter, r *http.Request) {
	g.l.Println("Processing POST request..Create Repo handler")
	vars := mux.Vars(r)
	orgName := vars["org"]
	ownerName := vars["owner"]
	//g.l.Println("Organization name..", orgName, ownerName)
	var createRepoReq service.CreateRepoRequest
	//g.l.Println("receieved models..", r.Body)
	err := json.NewDecoder(r.Body).Decode(&createRepoReq)
	if err != nil {
		g.l.Println("Error occurred while decoding the request data", err)
		http.Error(rw, "Error occurred while decoding the request data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	repoStatus, err := g.gbService.CreateRepo(orgName, ownerName, &createRepoReq)
	if err != nil {
		if err == service.ErrOrgNotFound || err == service.ErrOwnerNotFound {
			g.l.Println("Error occurred.", err)
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		g.l.Println("Error occurred.", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	g.l.Println("Repository got created.")
	rw.Header().Set("Content-Type", "Application/json")
	rw.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(rw).Encode(repoStatus)
	if err != nil {
		g.l.Println("Error occured while decoding the Git repo list", err)
		http.Error(rw, "Error occured while decoding the output", http.StatusInternalServerError)
		return
	}
}

func (g *GitRepo) DeleteRepoHandler(rw http.ResponseWriter, r *http.Request) {
	g.l.Println("Processing Delete Request..")
	vars := mux.Vars(r)
	orgName := vars["org"]
	ownerName := vars["owner"]
	repoName := vars["repo"]
	//	g.l.Println("Organization & Repo name..", orgName, ownerName, repoName)

	status, err := g.gbService.DeleteRepo(orgName, ownerName, repoName)
	if err != nil {
		g.l.Println("Error occurred while deleting the repo.", err)
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	if status {
		rw.WriteHeader(http.StatusNoContent)
		fmt.Fprintln(rw, "Repository got deleted")
		g.l.Println("Repository got deleted.")
	}
}

func (g *GitRepo) ListBranchesHandler(rw http.ResponseWriter, r *http.Request) {
	g.l.Println("Processing Get branch Request..")
	vars := mux.Vars(r)
	orgName := vars["org"]
	ownerName := vars["owner"]
	repoName := vars["repo"]
	//	g.l.Println("Owner & Repo name..", ownerName, repoName)

	branchList, err := g.gbService.ListBranches(orgName, ownerName, repoName)
	if err != nil {
		switch err {
		case service.ErrOwnerNotFound, service.ErrRepoNotFound, service.ErrBranchesNotFound, service.ErrOrgNotFound:
			g.l.Println("Error occurred while fetching the branch list.", err)
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		default:
			g.l.Println("Error occurred while fetching the branch list.", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	//g.l.Println("Retrieved Branch list..", branchList)
	g.l.Println("Retrieved Branch list.")
	rw.Header().Set("Content-Type", "Application/json")
	err = json.NewEncoder(rw).Encode(branchList)
	if err != nil {
		g.l.Println("Error occured while decoding the Git branch list.", err)
		http.Error(rw, "Error occured while decoding the output.", http.StatusInternalServerError)
		return
	}
}

func (g *GitRepo) CreateBranchHandler(rw http.ResponseWriter, r *http.Request) {
	g.l.Println("Processing Create branch Request..")
	vars := mux.Vars(r)
	orgName := vars["org"]
	ownerName := vars["owner"]
	repoName := vars["repo"]
	//	g.l.Println("Organization & Repo name..", orgName, ownerName, repoName)
	var cbreq service.CreateBranchRequest
	err := json.NewDecoder(r.Body).Decode(&cbreq)
	if err != nil {
		g.l.Println("Error occurred while decoding the request data", err)
		http.Error(rw, "Error occurred while decoding the request data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cbResp, err := g.gbService.CreateBranch(orgName, ownerName, repoName, &cbreq)
	if err != nil {
		if err == service.ErrBranchesAlreadyExists || err == service.ErrInvalidBranchName {
			g.l.Println("Error occurred while decoding the request data", err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		g.l.Println("Error occurred while decoding the request data", err)
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	g.l.Println("Received reponse for create branch", cbResp)
	rw.Header().Set("Content-Type", "Application/json")
	err = json.NewEncoder(rw).Encode(cbResp)
	if err != nil {
		g.l.Println("Error occurred while decoding the request data", err)
		http.Error(rw, "Error occurred while decoding the request data", http.StatusBadRequest)
		return
	}
}

func (g *GitRepo) DeleteBranchHandler(rw http.ResponseWriter, r *http.Request) {
	g.l.Println("Processing Delete branch Request..")
	vars := mux.Vars(r)
	orgName := vars["org"]
	ownerName := vars["owner"]
	repoName := vars["repo"]
	refName := vars["ref"]
	//	g.l.Println("Organization, Repo & branch name..", orgName, ownerName, repoName, refName)

	resp, err := g.gbService.DeleteBranch(orgName, ownerName, repoName, refName)
	if err != nil {
		g.l.Println("Error occurred while decoding the request data", err)
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	if resp {
		rw.WriteHeader(http.StatusNoContent)
		fmt.Fprintln(rw, "Repo got deleted")
		g.l.Println("Repo got deleted")
	} else {
		g.l.Println("Error occurred while deleting the branch.", repoName, refName)
		http.Error(rw, "Error occurred while deleting the repo.", http.StatusBadRequest)
		return
	}
}

// get /repos/{org}/{owner}/{Repo}/pulls
func (g *GitRepo) ListPRHandler(rw http.ResponseWriter, r *http.Request) {
	g.l.Println("Processing list PR Request..")
	vars := mux.Vars(r)
	orgName := vars["org"]
	ownerName := vars["owner"]
	repoName := vars["repo"]
	//	g.l.Println("Organization, Repo name..", ownerName, repoName)

	listPRs, err := g.gbService.ListPRs(orgName, ownerName, repoName)

	if err != nil {
		g.l.Println("Error occurred while fetching PRs list.", err)
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	//g.l.Println("Retrieved PRs list..", listPRs)
	g.l.Println("Retrieved PRs list.")
	rw.Header().Set("Content-Type", "Application/json")
	err = json.NewEncoder(rw).Encode(listPRs)
	if err != nil {
		g.l.Println("Error occured while decoding the Git branch list.", err)
		http.Error(rw, "Error occured while decoding the output.", http.StatusInternalServerError)
		return
	}
}

func (g *GitRepo) CreatePRHandler(rw http.ResponseWriter, r *http.Request) {
	g.l.Println("Processing Create PR Request..")
	vars := mux.Vars(r)
	orgName := vars["org"]
	ownerName := vars["owner"]
	repoName := vars["repo"]
	//	g.l.Println("Organization & Repo name..", ownerName, repoName)
	var prReq service.PRRequest
	err := json.NewDecoder(r.Body).Decode(&prReq)
	if err != nil {
		g.l.Println("Error occurred while decoding the request data", err)
		http.Error(rw, "Error occurred while decoding the request data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	prResp, err := g.gbService.CreatePR(orgName, ownerName, repoName, &prReq)
	if err != nil {
		switch err {
		case service.ErrOwnerNotFound, service.ErrRepoNotFound, service.ErrBranchesNotFound, service.ErrOrgNotFound:
			g.l.Println("Error occurred while fetching the branch list.", err)
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		default:
			g.l.Println("Error occurred while fetching the branch list.", err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}

	g.l.Println("Received reponse for create PR.", prResp)
	rw.Header().Set("Content-Type", "Application/json")
	err = json.NewEncoder(rw).Encode(prResp)
	if err != nil {
		g.l.Println("Error occurred while decoding the request data", err)
		http.Error(rw, "Error occurred while decoding the request data", http.StatusBadRequest)
		return
	}
}
