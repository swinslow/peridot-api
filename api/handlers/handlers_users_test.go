// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"testing"

	hu "github.com/swinslow/peridot-api/test/handlerutils"
)

func TestCanGetUsersHandlerAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users", nil, "admin")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)

	// expect full user data b/c we're an admin
	wanted := `{"users": [{"id": 1, "name": "Admin", "github": "admin", "access": "admin"}, {"id": 2, "name": "Operator", "github": "operator", "access": "operator"}, {"id": 3, "name": "Commenter", "github": "commenter", "access": "commenter"}, {"id": 4, "name": "Viewer", "github": "viewer", "access": "viewer"}, {"id": 10, "name": "Disabled", "github": "disabled", "access": "disabled"}]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCanGetUsersHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users", nil, "operator")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)

	// expect less data available here since user is just an operator, not admin
	wanted := `{"users": [{"id": 1, "github": "admin"}, {"id": 2, "github": "operator"}, {"id": 3, "github": "commenter"}, {"id": 4, "github": "viewer"}, {"id": 10, "github": "disabled"}]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCanGetUsersHandlerAsCommenter(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users", nil, "commenter")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)

	// expect less data available here since user is just an operator, not admin
	wanted := `{"users": [{"id": 1, "github": "admin"}, {"id": 2, "github": "operator"}, {"id": 3, "github": "commenter"}, {"id": 4, "github": "viewer"}, {"id": 10, "github": "disabled"}]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCanGetUsersHandlerAsViewer(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users", nil, "viewer")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)

	// expect less data available here since user is just an operator, not admin
	wanted := `{"users": [{"id": 1, "github": "admin"}, {"id": 2, "github": "operator"}, {"id": 3, "github": "commenter"}, {"id": 4, "github": "viewer"}, {"id": 10, "github": "disabled"}]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetUsersHandlerAsDisabledUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users", nil, "disabled")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmDisabledAuth(t, rec)
}

func TestCannotGetUsersHandlerAsInvalidUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users", nil, "invalid")
	http.HandlerFunc(env.usersHandler).ServeHTTP(rec, req)
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}
