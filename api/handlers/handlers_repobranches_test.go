// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"testing"

	hu "github.com/swinslow/peridot-api/test/handlerutils"
)

// ===== GET /repos/2/branches =====

func TestCanGetRepoBranchesSubHandler(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repos/2/branches", ``, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoBranchesSubHandler), "/repos/{id}/branches")
	hu.ConfirmOKResponse(t, rec)

	// should be returned in alphabetical order
	wanted := `{"branches": ["alpha", "beta", "master"]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetRepoBranchesSubHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repos/2/branches", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoBranchesSubHandler), "/repos/{id}/branches")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/repos/2/branches", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoBranchesSubHandler), "/repos/{id}/branches")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== POST /repos/2/branches =====

func TestCanPostRepoBranchesSubHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/repos/2/branches", `{"branch": "new-branch"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoBranchesSubHandler), "/repos/{id}/branches")
	hu.ConfirmCreatedResponse(t, rec)

	wanted := `{"branch": "new-branch"}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	rbs, err := env.db.GetAllRepoBranchesForRepoID(2)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(rbs) != 4 {
		t.Errorf("expected %d, got %d", 4, len(rbs))
	}
}

func TestCannotPostRepoBranchesSubHandlerAsOtherUser(t *testing.T) {
	// as commenter
	rec, req, env := setupTestEnv(t, "POST", "/repos/2/branches", `{"branch": "new-branch"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoBranchesSubHandler), "/repos/{id}/branches")
	hu.ConfirmAccessDenied(t, rec)

	// as viewer
	rec, req, env = setupTestEnv(t, "POST", "/repos/2/branches", `{"branch": "new-branch"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoBranchesSubHandler), "/repos/{id}/branches")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCannotPostRepoBranchesSubHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/repos/2/branches", `{"branch": "new-branch"}`, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoBranchesSubHandler), "/repos/{id}/branches")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "POST", "/repos/2/branches", `{"branch": "new-branch"}`, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoBranchesSubHandler), "/repos/{id}/branches")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}
