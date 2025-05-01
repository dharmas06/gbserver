package handlers

import (
	"fmt"
	"gbserver/service"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var l = log.New(os.Stdout, "gbTestServer ", log.LstdFlags)
var gitTestRepo = NewGitRepo(l)

func TestListRepoHandler(t *testing.T) {
	type got struct {
		req *http.Request
	}
	type want struct {
		statusCode int
		//	response   []service.RepoResponse
		wantErr error
	}
	tests := []struct {
		name string
		got  got
		want want
	}{
		{
			name: "Test invalid user",
			got: got{
				req: generateRequest(t, "InvalidUser"),
			},
			want: want{
				statusCode: 404,
				wantErr:    service.ErrOwnerNotFound,
			},
		},
		{
			name: "Test repo results for valid org & user",
			got: got{
				req: generateRequest(t, "ValidData"),
			},
			want: want{
				statusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		resp := httptest.NewRecorder()
		gitTestRepo.ListRepoHandler(resp, tt.got.req)

		assert.Equal(t, tt.want.statusCode, resp.Code)
		fmt.Println(resp.Body, resp.Code)

	}
}

func generateRequest(_ *testing.T, testType string) *http.Request {

	var httpReq *http.Request
	if testType == "InvalidUser" {
		httpReq = httptest.NewRequest(http.MethodGet, "http://localhost:9090/orgs/{org}/{owner}/repos", nil)

		httpReq = mux.SetURLVars(httpReq, map[string]string{
			"org":   "gborg",
			"owner": "testuser1",
		})

		return httpReq
	}
	if testType == "ValidData" {
		httpReq = httptest.NewRequest(http.MethodGet, "http://localhost:9090/orgs/{org}/{owner}/repos", nil)

		httpReq = mux.SetURLVars(httpReq, map[string]string{
			"org":   "gborg",
			"owner": "gbuser",
		})

		return httpReq
	}
	return httpReq

}
