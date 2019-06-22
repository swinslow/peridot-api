// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"testing"

	"github.com/swinslow/peridot-api/internal/datastore"
	hu "github.com/swinslow/peridot-api/test/handlerutils"
)

// ===== GET /users =====

func TestCanGetUsersHandlerAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users", "", "admin")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)

	// expect full user data b/c we're an admin
	wanted := `{"users": [{"id": 1, "name": "Admin", "github": "admin", "access": "admin"}, {"id": 2, "name": "Operator", "github": "operator", "access": "operator"}, {"id": 3, "name": "Commenter", "github": "commenter", "access": "commenter"}, {"id": 4, "name": "Viewer", "github": "viewer", "access": "viewer"}, {"id": 10, "name": "Disabled", "github": "disabled", "access": "disabled"}]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCanGetUsersHandlerAsOtherUsers(t *testing.T) {
	// should be same return for all
	wanted := `{"users": [{"id": 1, "github": "admin"}, {"id": 2, "github": "operator"}, {"id": 3, "github": "commenter"}, {"id": 4, "github": "viewer"}, {"id": 10, "github": "disabled"}]}`

	// as operator
	rec, req, env := setupTestEnv(t, "GET", "/users", "", "operator")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)

	// as commenter
	rec, req, env = setupTestEnv(t, "GET", "/users", "", "commenter")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)

	// as viewer
	rec, req, env = setupTestEnv(t, "GET", "/users", "", "viewer")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetUsersHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users", "", "disabled")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/users", "", "invalid")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== POST /users =====

func TestCanPostUsersHandlerAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/users", `{"name": "Steve", "github": "swinslow", "access": "operator"}`, "admin")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"success": true, "id": 11}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	users, err := env.db.GetAllUsers()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(users) != 6 {
		t.Errorf("expected %d, got %d", 6, len(users))
	}
	newUser, err := env.db.GetUserByID(11)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	wantedUser := &datastore.User{ID: 11, Name: "Steve", Github: "swinslow", AccessLevel: datastore.AccessOperator}
	if newUser.ID != wantedUser.ID || newUser.Name != wantedUser.Name || newUser.Github != wantedUser.Github || newUser.AccessLevel != wantedUser.AccessLevel {
		t.Errorf("expected %#v, got %#v", wantedUser, newUser)
	}
}

func TestCannotPostUsersHandlerAsOtherUser(t *testing.T) {
	// as operator
	rec, req, env := setupTestEnv(t, "POST", "/users", `{"name": "Steve", "github": "swinslow", "access": "operator"}`, "operator")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmAccessDenied(t, rec)

	// as commenter
	rec, req, env = setupTestEnv(t, "POST", "/users", `{"name": "Steve", "github": "swinslow", "access": "commenter"}`, "operator")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmAccessDenied(t, rec)

	// as viewer
	rec, req, env = setupTestEnv(t, "POST", "/users", `{"name": "Steve", "github": "swinslow", "access": "viewer"}`, "operator")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmAccessDenied(t, rec)
}

// ===== GET /users/3 =====

func TestCanGetUsersOneHandlerAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users/3", "", "admin")
	http.HandlerFunc(env.usersOneHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)

	// expect full user data b/c we're an admin
	wanted := `{"success": true, "user": {"id": 3, "name": "Commenter", "github": "commenter", "access": "commenter"}}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCanGetUsersOneHandlerAsOtherUsers(t *testing.T) {
	// should be same return for all EXCEPT self
	wanted := `{"success": true, "user": {"id": 3, "github": "commenter"}}`

	// as operator
	rec, req, env := setupTestEnv(t, "GET", "/users/3", "", "operator")
	http.HandlerFunc(env.usersOneHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)

	// as viewer
	rec, req, env = setupTestEnv(t, "GET", "/users/3", "", "viewer")
	http.HandlerFunc(env.usersOneHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)

	// and for commenter getting somebody else's data
	wanted = `{"success": true, "user": {"id": 1, "github": "admin"}}`
	rec, req, env = setupTestEnv(t, "GET", "/users/1", "", "commenter")
	http.HandlerFunc(env.usersOneHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)
}

func TestCanGetUsersOneHandlerAsSelf(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users/3", "", "commenter")
	http.HandlerFunc(env.usersOneHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)

	// expect full user data b/c we're getting our own data
	wanted := `{"success": true, "user": {"id": 3, "name": "Commenter", "github": "commenter", "access": "commenter"}}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetUsersOneHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users/3", "", "disabled")
	http.HandlerFunc(env.usersOneHandler).ServeHTTP(rec, req)
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/users/3", "", "invalid")
	http.HandlerFunc(env.usersOneHandler).ServeHTTP(rec, req)
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}
