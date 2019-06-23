// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"testing"

	"github.com/swinslow/peridot-api/internal/datastore"
	hu "github.com/swinslow/peridot-api/test/handlerutils"
)

// ===== GET /repos =====

func TestCanGetReposHandler(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repos", ``, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposHandler), "/repos")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"repos": [{"id": 1, "subproject_id": 2, "name": "repo1", "address": "https://example.com/repo1.git"},{"id": 2, "subproject_id": 4, "name": "repo2", "address": "https://example.com/repo2.git"},{"id": 3, "subproject_id": 4, "name": "repo3", "address": "https://example.com/repo3.git"},{"id": 4, "subproject_id": 4, "name": "repo4", "address": "https://example.com/repo4.git"}]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetReposHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repos", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposHandler), "/repos")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/repos", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposHandler), "/repos")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== POST /repos =====

func TestCanPostReposHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/repos", `{"subproject_id": 2, "name": "repo5", "address": "https://example.com/newrepo5.git"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposHandler), "/repos")
	hu.ConfirmCreatedResponse(t, rec)

	wanted := `{"id": 5}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	repos, err := env.db.GetAllRepos()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(repos) != 5 {
		t.Errorf("expected %d, got %d", 5, len(repos))
	}
	newRepo, err := env.db.GetRepoByID(5)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	wantedRepo := &datastore.Repo{ID: 5, SubprojectID: 2, Name: "repo5", Address: "https://example.com/newrepo5.git"}
	if newRepo.ID != wantedRepo.ID || newRepo.Name != wantedRepo.Name || newRepo.Address != wantedRepo.Address {
		t.Errorf("expected %#v, got %#v", wantedRepo, newRepo)
	}
}

func TestCannotPostReposHandlerAsOtherUser(t *testing.T) {
	// as commenter
	rec, req, env := setupTestEnv(t, "POST", "/repos", `{"subproject_id": 2, "name": "repo5", "address": "https://example.com/newrepo5.git"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposHandler), "/repos")
	hu.ConfirmAccessDenied(t, rec)

	// as viewer
	rec, req, env = setupTestEnv(t, "POST", "/repos", `{"subproject_id": 2, "name": "repo5", "address": "https://example.com/newrepo5.git"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposHandler), "/repos")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCannotPostReposHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/repos", `{"subproject_id": 2, "name": "repo5", "address": "https://example.com/newrepo5.git"}`, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposHandler), "/repos")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "POST", "/repos", `{"subproject_id": 2, "name": "repo5", "address": "https://example.com/newrepo5.git"}`, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposHandler), "/repos")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== GET /subprojects/4/repos =====

func TestCanGetReposSubHandler(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/subprojects/4/repos", ``, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposSubHandler), "/subprojects/{id}/repos")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"repos": [{"id": 2, "subproject_id": 4, "name": "repo2", "address": "https://example.com/repo2.git"},{"id": 3, "subproject_id": 4, "name": "repo3", "address": "https://example.com/repo3.git"},{"id": 4, "subproject_id": 4, "name": "repo4", "address": "https://example.com/repo4.git"}]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetReposSubHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/subprojects/4/repos", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposSubHandler), "/subprojects/{id}/repos")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/subprojects/4/repos", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposSubHandler), "/subprojects/{id}/repos")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== POST /subprojects/4/repos =====

func TestCanPostReposSubHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/subprojects/4/repos", `{"name": "repo5", "address": "https://example.com/newrepo5.git"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposSubHandler), "/subprojects/{id}/repos")
	hu.ConfirmCreatedResponse(t, rec)

	wanted := `{"id": 5}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	repos, err := env.db.GetAllRepos()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(repos) != 5 {
		t.Errorf("expected %d, got %d", 5, len(repos))
	}
	newRepo, err := env.db.GetRepoByID(5)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	wantedRepo := &datastore.Repo{ID: 5, SubprojectID: 4, Name: "repo5", Address: "https://example.com/newrepo5.git"}
	if newRepo.ID != wantedRepo.ID || newRepo.Name != wantedRepo.Name || newRepo.Address != wantedRepo.Address {
		t.Errorf("expected %#v, got %#v", wantedRepo, newRepo)
	}
}

func TestCannotPostReposSubHandlerAsOtherUser(t *testing.T) {
	// as commenter
	rec, req, env := setupTestEnv(t, "POST", "/subprojects/4/repos", `{"name": "repo5", "address": "https://example.com/newrepo5.git"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposSubHandler), "/subprojects/{id}/repos")
	hu.ConfirmAccessDenied(t, rec)

	// as viewer
	rec, req, env = setupTestEnv(t, "POST", "/subprojects/4/repos", `{"name": "repo5", "address": "https://example.com/newrepo5.git"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposSubHandler), "/subprojects/{id}/repos")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCannotPostReposSubHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/subprojects/4/repos", `{"name": "repo5", "address": "https://example.com/newrepo5.git"}`, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposSubHandler), "/subprojects/{id}/repos")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "POST", "/subprojects/4/repos", `{"name": "repo5", "address": "https://example.com/newrepo5.git"}`, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposSubHandler), "/subprojects/{id}/repos")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== GET /repos/3 =====

func TestCanGetReposOneHandlerAsViewer(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repos/3", "", "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposOneHandler), "/repos/{id}")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"repo": {"id": 3, "subproject_id": 4, "name": "repo3", "address": "https://example.com/repo3.git"}}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetReposOneHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repos/3", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposOneHandler), "/repos/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/repos/3", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposOneHandler), "/repos/{id}")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== PUT /repos/3 =====

func TestCanPutReposOneHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/repos/3", `{"name": "new-name", "address": "https://example.com/new-name.git"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposOneHandler), "/repos/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	repo, err := env.db.GetRepoByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedRepo := &datastore.Repo{ID: 3, SubprojectID: 4, Name: "new-name", Address: "https://example.com/new-name.git"}
	if repo.ID != wantedRepo.ID || repo.SubprojectID != wantedRepo.SubprojectID || repo.Name != wantedRepo.Name || repo.Address != wantedRepo.Address {
		t.Errorf("expected %#v, got %#v", wantedRepo, repo)
	}
}

func TestCanPutReposOneHandlerAsOperatorWithJustName(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/repos/3", `{"name": "new-name"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposOneHandler), "/repos/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	repo, err := env.db.GetRepoByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedRepo := &datastore.Repo{ID: 3, SubprojectID: 4, Name: "new-name", Address: "https://example.com/repo3.git"}
	if repo.ID != wantedRepo.ID || repo.SubprojectID != wantedRepo.SubprojectID || repo.Name != wantedRepo.Name || repo.Address != wantedRepo.Address {
		t.Errorf("expected %#v, got %#v", wantedRepo, repo)
	}
}

func TestCanPutReposOneHandlerAsOperatorWithJustFullname(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/repos/3", `{"address": "https://example.com/newRepo3.git"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposOneHandler), "/repos/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	repo, err := env.db.GetRepoByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedRepo := &datastore.Repo{ID: 3, SubprojectID: 4, Name: "repo3", Address: "https://example.com/newRepo3.git"}
	if repo.ID != wantedRepo.ID || repo.SubprojectID != wantedRepo.SubprojectID || repo.Name != wantedRepo.Name || repo.Address != wantedRepo.Address {
		t.Errorf("expected %#v, got %#v", wantedRepo, repo)
	}
}

func TestCannotPutReposOneHandlerAsCommenter(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/repos/3", `{"name": "new-name", "address": "https://example.com/new-name.git"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposOneHandler), "/repos/{id}")
	hu.ConfirmAccessDenied(t, rec)

	// and verify state of database now
	repo, err := env.db.GetRepoByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedRepo := &datastore.Repo{ID: 3, SubprojectID: 4, Name: "repo3", Address: "https://example.com/repo3.git"}
	if repo.ID != wantedRepo.ID || repo.SubprojectID != wantedRepo.SubprojectID || repo.Name != wantedRepo.Name || repo.Address != wantedRepo.Address {
		t.Errorf("expected %#v, got %#v", wantedRepo, repo)
	}
}

// ===== DELETE /repos/3 =====

func TestCanDeleteReposOneHandlerAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "DELETE", "/repos/3", ``, "admin")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposOneHandler), "/repos/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	repos, err := env.db.GetAllRepos()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(repos) != 3 {
		t.Errorf("expected %d, got %d", 3, len(repos))
	}
	repo, err := env.db.GetRepoByID(3)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil and %#v", repo)
	}
}

func TestCannotDeleteReposOneHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "DELETE", "/repos/3", ``, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.reposOneHandler), "/repos/{id}")
	hu.ConfirmAccessDenied(t, rec)

	// and verify state of database has not changed
	repos, err := env.db.GetAllRepos()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(repos) != 4 {
		t.Errorf("expected %d, got %d", 4, len(repos))
	}
	repo, err := env.db.GetRepoByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedRepo := &datastore.Repo{ID: 3, SubprojectID: 4, Name: "repo3", Address: "https://example.com/repo3.git"}
	if repo.ID != wantedRepo.ID || repo.SubprojectID != wantedRepo.SubprojectID || repo.Name != wantedRepo.Name || repo.Address != wantedRepo.Address {
		t.Errorf("expected %#v, got %#v", wantedRepo, repo)
	}
}
