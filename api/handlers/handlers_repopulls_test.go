// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"testing"

	"github.com/swinslow/peridot-db/pkg/datastore"
	hu "github.com/swinslow/peridot-api/test/handlerutils"
)

// ===== GET /repos/2/branches/alpha =====

func TestCanGetRepoPullsSubHandler(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repos/2/branches/master", ``, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoPullsSubHandler), "/repos/{id}/branches/{branch}")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"pulls": [
		{"id":1,"repo_id":2,"branch":"master","started_at":"0001-01-01T00:00:00Z","finished_at":"0001-01-01T00:00:00Z","status":"stopped","health":"error","commit":"abcdef012345abcdef012345abcdef0123451234","tag":"v1.1","spdx_id":""},
		{"id":2,"repo_id":2,"branch":"master","started_at":"0001-01-01T00:00:00Z","finished_at":"0001-01-01T00:00:00Z","status":"stopped","health":"ok","commit":"abcdef012345abcdef012345abcdef0123455678","tag":"v1.2","spdx_id":""}
	]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetRepoPullsSubHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repos/2/branches/master", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoPullsSubHandler), "/repos/{id}/branches/{branch}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/repos/2/branches/master", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoPullsSubHandler), "/repos/{id}/branches/{branch}")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== POST /repos/2/branches/alpha =====

func TestCanPostRepoPullsSubHandlerWithCommitAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/repos/2/branches/alpha", `{"commit": "123490ab56123490ab56123490ab56123490ab56"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoPullsSubHandler), "/repos/{id}/branches/{branch}")
	hu.ConfirmCreatedResponse(t, rec)

	wanted := `{"id": 5}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	rps, err := env.db.GetAllRepoPullsForRepoBranch(2, "alpha")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(rps) != 1 {
		t.Errorf("expected %d, got %d", 1, len(rps))
	}
	newRepoPull, err := env.db.GetRepoPullByID(5)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	wantedRepoPull := &datastore.RepoPull{ID: 5, RepoID: 2, Branch: "alpha", Status: datastore.StatusStartup, Health: datastore.HealthOK, Commit: "123490ab56123490ab56123490ab56123490ab56"}
	if newRepoPull.ID != wantedRepoPull.ID || newRepoPull.RepoID != wantedRepoPull.RepoID || newRepoPull.Branch != wantedRepoPull.Branch || newRepoPull.Status != wantedRepoPull.Status || newRepoPull.Health != wantedRepoPull.Health || newRepoPull.Commit != wantedRepoPull.Commit {
		t.Errorf("expected %#v, got %#v", wantedRepoPull, newRepoPull)
	}
}

// ===== GET /repopulls/3 =====

func TestCanGetRepoPullsOneHandlerAsViewer(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repopulls/3", "", "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoPullsOneHandler), "/repopulls/{id}")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"repopull": {"id":3,"repo_id":4,"branch":"dev","started_at":"0001-01-01T00:00:00Z","finished_at":"0001-01-01T00:00:00Z","status":"running","health":"degraded","commit":"abcdef012345abcdef012345abcdef01234590ab","spdx_id":""}}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetRepoPullsOneHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repopulls/3", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoPullsOneHandler), "/repopulls/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/repopulls/3", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoPullsOneHandler), "/repopulls/{id}")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== DELETE /repopulls/3 =====

func TestCanDeleteRepoPullsOneHandlerAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "DELETE", "/repopulls/2", ``, "admin")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoPullsOneHandler), "/repopulls/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	rps, err := env.db.GetAllRepoPullsForRepoBranch(2, "master")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(rps) != 1 {
		t.Errorf("expected %d, got %d", 1, len(rps))
	}
	rp, err := env.db.GetRepoPullByID(2)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil and %#v", rp)
	}
}

func TestCannotDeleteRepoPullsOneHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "DELETE", "/repopulls/2", ``, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.repoPullsOneHandler), "/repopulls/{id}")
	hu.ConfirmAccessDenied(t, rec)

	// and verify state of database has not changed
	rps, err := env.db.GetAllRepoPullsForRepoBranch(2, "master")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(rps) != 2 {
		t.Errorf("expected %d, got %d", 2, len(rps))
	}
	rp, err := env.db.GetRepoPullByID(2)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	wantedRepoPull := &datastore.RepoPull{ID: 2, RepoID: 2, Branch: "master", Status: datastore.StatusStopped, Health: datastore.HealthOK, Commit: "abcdef012345abcdef012345abcdef0123455678", Tag: "v1.2"}
	if rp.ID != wantedRepoPull.ID || rp.RepoID != wantedRepoPull.RepoID || rp.Branch != wantedRepoPull.Branch || rp.Status != wantedRepoPull.Status || rp.Health != wantedRepoPull.Health || rp.Commit != wantedRepoPull.Commit || rp.Tag != wantedRepoPull.Tag {
		t.Errorf("expected %#v, got %#v", wantedRepoPull, rp)
	}

}
