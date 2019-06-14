// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/swinslow/peridot-api/internal/datastore"
)

// getTestEnv creates the Env object used for the handlers
// unit test suite. It is not exported and should NEVER be
// called by production code.
func getTestEnv() *Env {
	db := createMockDB()

	oauthConf := &oauth2.Config{
		ClientID:     "abcdef0123abcdef4567",
		ClientSecret: "abcdef0123abcdef4567abcdef8901abcdef2345",
		Scopes:       []string{"user:email"},
		Endpoint:     githuboauth.Endpoint,
	}

	env := &Env{
		db:           db,
		jwtSecretKey: "keyForTesting",
		oauthConf:    oauthConf,
		oauthState:   "nonRandomStateString",
	}
	return env
}


// loginWithTestUser loads the mock user with the given
// github name, adds it to the request context, and
// returns the request object. If ghUsername is "invalid", it
// will be set with an invalid user ID (ID: 0, AccessLevel:
// AccessDisabled).
func loginWithTestUser(t *testing.T, req *http.Request, env *Env, ghUsername string) *http.Request {
	var user *datastore.User
	var err error

	if ghUsername == "invalid" {
		user = &datastore.User{ID: 0, AccessLevel: datastore.AccessDisabled}
	} else {
		user, err = env.db.GetUserByGithub(ghUsername)
		if err != nil {
			t.Fatalf("error getting mock user by github name: %v", err)
			return nil
		}
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey(0), user)
	return req.WithContext(ctx)
}

// setupTestEnv sets up the test infrastructure, including
// the HTTP request and response, the Env object and the
// logged-in user context. If ghUsername is "invalid", it
// will be set with an invalid user ID (ID: 0, AccessLevel:
// AccessDisabled).
func setupTestEnv(t *testing.T, method string, endpoint string, body io.Reader, ghUsername string) (*httptest.ResponseRecorder, *http.Request, *Env) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	env := getTestEnv()
	req = loginWithTestUser(t, req, env, ghUsername)

	return rec, req, env
}
