package server

import (
	"fmt"
	"gbserver/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/didip/tollbooth/v8"
	"github.com/didip/tollbooth/v8/limiter"
	"github.com/gorilla/mux"
)

var BaseURL = " "
var ServerPort = ":9090"

func TollboothMiddleware(limiter *limiter.Limiter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, r.URL.Path)
			httpError := tollbooth.LimitByRequest(limiter, w, r)
			//fmt.Println(httpError)
			if httpError != nil {
				// If rate limit exceeded
				w.Header().Add("Content-Type", limiter.GetMessageContentType())
				w.WriteHeader(httpError.StatusCode)
				w.Write([]byte(httpError.Message))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func StartServer() {

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	l := log.New(os.Stdout, "gbServer ", log.LstdFlags)
	gbH := handlers.NewGitRepo(l)

	limit := tollbooth.NewLimiter(1, nil)
	limit.SetIPLookup(limiter.IPLookup{
		Name:           "RemoteAddr", // "X-Real-IP"
		IndexFromRight: 0,
	})
	//	limit.SetMethods([]string{"GET"})
	limit.SetMessage("Reached maximum request limit.")

	router := mux.NewRouter()
	router.Use(TollboothMiddleware(limit))

	apiRouter := router.PathPrefix("/").Subrouter()
	//get  /orgs/{org}/repos
	apiRouter.Path("/orgs/{org}/repos").Methods(http.MethodGet).HandlerFunc(gbH.ListRepoHandler)

	//post   /orgs/{org}/repos
	apiRouter.Path("/orgs/{org}/repos").Methods(http.MethodPost).HandlerFunc(gbH.CreateRepoHandler)

	//delete /Repos/{owner}/{Repo}
	apiRouter.Path("/repos/{owner}/{repo}").Methods(http.MethodDelete).HandlerFunc(gbH.DeleteRepoHandler)

	// get /Repos/{owner}/{Repo}/branches
	apiRouter.Path("/repos/{owner}/{repo}/branches").Methods(http.MethodGet).HandlerFunc(gbH.ListBranchesHandler)

	// post /Repos/{owner}/{Repo}/git/Refs
	apiRouter.Path("/repos/{owner}/{repo}/git/refs").Methods(http.MethodPost).HandlerFunc(gbH.CreateBranchHandler)

	//delete /Repos/{owner}/{Repo}/git/Refs/{Ref}
	apiRouter.Path("/repos/{owner}/{repo}/git/refs/{ref}").Methods(http.MethodDelete).HandlerFunc(gbH.DeleteBranchHandler)

	// get /repos/{owner}/{repo}/pulls
	apiRouter.Path("/repos/{owner}/{repo}/pulls").Methods(http.MethodGet).HandlerFunc(gbH.ListPRHandler)

	// post /repos/{owner}/{Repo}/pulls
	apiRouter.Path("/repos/{owner}/{repo}/pulls").Methods(http.MethodPost).HandlerFunc(gbH.CreatePRHandler)

	//patch /repos/{owner}/{repo}/pulls/{pull_number} State - closed
	apiRouter.Path("/repos/{owner}/{repo}/pulls/{pull_number}").Methods(http.MethodPatch).HandlerFunc(gbH.UpdatePRHandler)

	go func() {
		log.Println("Starting GB server on ..", ServerPort)
		log.Fatal(http.ListenAndServe(ServerPort, router))

	}()
	<-sigChan
	fmt.Println("Stopping GB server..")
}
