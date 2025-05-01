package server

import (
	"context"
	"fmt"
	"gbserver/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didip/tollbooth/v8"
	"github.com/didip/tollbooth/v8/limiter"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var BaseURL = " "
var ServerPort = ":9090"

var ReqLimit float64 = 10

type contextKey string

const requestIDKey = contextKey("requestID")

func uuidMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := uuid.New().String()
		w.Header().Set("X-Request-ID", uuid)
		ctx := context.WithValue(r.Context(), requestIDKey, uuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loggingMiddleware prints the request ID and the request details.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			uuid = "unknown"
		}

		start := time.Now()
		log.Printf("Started %s %s (Request ID: %s)", r.Method, r.RequestURI, uuid)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v (Request ID: %s)", r.RequestURI, time.Since(start), uuid)
	})
}

func TollboothMiddleware(limiter *limiter.Limiter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, r.URL.Path)
			httpError := tollbooth.LimitByRequest(limiter, w, r)

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

	limit := tollbooth.NewLimiter(ReqLimit, nil)
	limit.SetIPLookup(limiter.IPLookup{
		Name:           "RemoteAddr",
		IndexFromRight: 0,
	})

	limit.SetMessage("Reached maximum request limit.")

	router := mux.NewRouter()
	//	router.Use(uuidMiddleware)
	//	router.Use(loggingMiddleware)
	router.Use(TollboothMiddleware(limit))

	apiRouter := router.PathPrefix("/").Subrouter()
	//get  /orgs/{org}/{owner}/repos
	apiRouter.Path("/orgs/{org}/{owner}/repos").Methods(http.MethodGet).HandlerFunc(gbH.ListRepoHandler)

	// //post   /orgs/{org}/{owner}/repos
	apiRouter.Path("/orgs/{org}/{owner}/repos").Methods(http.MethodPost).HandlerFunc(gbH.CreateRepoHandler)

	// //delete /Repos/{org}/{owner}/{Repo}
	apiRouter.Path("/repos/{org}/{owner}/{repo}").Methods(http.MethodDelete).HandlerFunc(gbH.DeleteRepoHandler)

	// // get /Repos/{org}/{owner}/{Repo}/branches
	apiRouter.Path("/repos/{org}/{owner}/{repo}/branches").Methods(http.MethodGet).HandlerFunc(gbH.ListBranchesHandler)

	// // post /Repos/{org}/{owner}/{Repo}/git/Refs
	apiRouter.Path("/repos/{org}/{owner}/{repo}/git/refs").Methods(http.MethodPost).HandlerFunc(gbH.CreateBranchHandler)

	// //delete /Repos/{org}/{owner}/{Repo}/git/Refs/{Ref}
	apiRouter.Path("/repos/{org}/{owner}/{repo}/git/refs/{ref}").Methods(http.MethodDelete).HandlerFunc(gbH.DeleteBranchHandler)

	// // get /repos/{org}/{owner}/{repo}/pulls
	apiRouter.Path("/repos/{org}/{owner}/{repo}/pulls").Methods(http.MethodGet).HandlerFunc(gbH.ListPRHandler)

	// // post /repos/{org}/{owner}/{Repo}/pulls
	apiRouter.Path("/repos/{org}/{owner}/{repo}/pulls").Methods(http.MethodPost).HandlerFunc(gbH.CreatePRHandler)

	// //patch /repos/{org}/{owner}/{repo}/pulls/{pull_number} State - closed
	apiRouter.Path("/repos/{org}/{owner}/{repo}/pulls/{pull_number}").Methods(http.MethodPatch).HandlerFunc(gbH.UpdatePRHandler)

	go func() {
		log.Println("Starting GB server on ..", ServerPort)
		log.Fatal(http.ListenAndServe(ServerPort, router))

	}()
	<-sigChan
	fmt.Println("Stopping GB server..")
}
