// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"testing"

	hu "github.com/swinslow/peridot-api/test/handlerutils"
)

func TestCanClearDBAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/admin/db", `{"command": "resetDB"}`, "admin")
	http.HandlerFunc(env.adminDBHandler).ServeHTTP(rec, req)
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"success": true}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	users, err := env.db.GetAllUsers()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected %d, got %d", 1, len(users))
	}
}

func TestAdminDBRequiresJSON(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/admin/db", `command: oops`, "admin")
	http.HandlerFunc(env.adminDBHandler).ServeHTTP(rec, req)
	hu.ConfirmBadRequestResponse(t, rec)

	wanted := `{"success": false, "error": "Invalid JSON request"}`
	hu.CheckResponse(t, rec, wanted)
}

func TestAdminDBRequiresCommand(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/admin/db", `{}`, "admin")
	http.HandlerFunc(env.adminDBHandler).ServeHTTP(rec, req)
	hu.ConfirmBadRequestResponse(t, rec)

	wanted := `{"success": false, "error": "No command specified"}`
	hu.CheckResponse(t, rec, wanted)
}

func TestAdminDBRequiresKnownCommand(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/admin/db", `{"command": "oops"}`, "admin")
	http.HandlerFunc(env.adminDBHandler).ServeHTTP(rec, req)
	hu.ConfirmBadRequestResponse(t, rec)

	wanted := `{"success": false, "error": "Unknown command 'oops'"}`
	hu.CheckResponse(t, rec, wanted)
}
func TestCannotClearDBUnlessAdmin(t *testing.T) {
	// as operator
	rec, req, env := setupTestEnv(t, "POST", "/admin/db", `{"command": "resetDB"}`, "operator")
	http.HandlerFunc(env.adminDBHandler).ServeHTTP(rec, req)
	hu.ConfirmAccessDenied(t, rec)

	// as commenter
	rec, req, env = setupTestEnv(t, "POST", "/admin/db", `{"command": "resetDB"}`, "commenter")
	http.HandlerFunc(env.adminDBHandler).ServeHTTP(rec, req)
	hu.ConfirmAccessDenied(t, rec)

	// as viewer
	rec, req, env = setupTestEnv(t, "POST", "/admin/db", `{"command": "resetDB"}`, "viewer")
	http.HandlerFunc(env.adminDBHandler).ServeHTTP(rec, req)
	hu.ConfirmAccessDenied(t, rec)
}
