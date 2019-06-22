// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"testing"

	"github.com/swinslow/peridot-api/internal/datastore"
	hu "github.com/swinslow/peridot-api/test/handlerutils"
)

// ===== GET /projects =====

func TestCanGetProjectsHandler(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/projects", ``, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsHandler), "/projects")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"projects": [{"id": 1, "name": "prj1", "fullname": "project 1"}, {"id": 2, "name": "prj2", "fullname": "project 2"}, {"id": 3, "name": "prj3", "fullname": "project 3"}]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetProjectsHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/projects", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsHandler), "/projects")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/projects", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsHandler), "/projects")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== POST /projects =====

func TestCanPostProjectsHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/projects", `{"name": "prj4", "fullname": "project 4"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsHandler), "/projects")
	hu.ConfirmCreatedResponse(t, rec)

	wanted := `{"success": true, "id": 4}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	projects, err := env.db.GetAllProjects()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(projects) != 4 {
		t.Errorf("expected %d, got %d", 4, len(projects))
	}
	newProject, err := env.db.GetProjectByID(4)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	wantedProject := &datastore.Project{ID: 4, Name: "prj4", Fullname: "project 4"}
	if newProject.ID != wantedProject.ID || newProject.Name != wantedProject.Name || newProject.Fullname != wantedProject.Fullname {
		t.Errorf("expected %#v, got %#v", wantedProject, newProject)
	}
}

func TestCannotPostProjectsHandlerAsOtherUser(t *testing.T) {
	// as commenter
	rec, req, env := setupTestEnv(t, "POST", "/projects", `{"name": "prj4", "fullname": "project 4"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsHandler), "/projects")
	hu.ConfirmAccessDenied(t, rec)

	// as viewer
	rec, req, env = setupTestEnv(t, "POST", "/projects", `{"name": "prj4", "fullname": "project 4"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsHandler), "/projects")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCannotPostProjectsHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/projects", `{"name": "prj4", "fullname": "project 4"}`, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsHandler), "/projects")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "POST", "/projects", `{"name": "prj4", "fullname": "project 4"}`, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsHandler), "/projects")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== GET /projects/3 =====

func TestCanGetProjectsOneHandlerAsViewer(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/projects/3", "", "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsOneHandler), "/projects/{id}")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"success": true, "project": {"id": 3, "name": "prj3", "fullname": "project 3"}}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetProjectsOneHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/projects/3", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsOneHandler), "/projects/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/projects/3", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsOneHandler), "/projects/{id}")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== PUT /projects/3 =====

func TestCanPutProjectsOneHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/projects/3", `{"name": "new-name", "fullname": "new-fullname"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsOneHandler), "/projects/{id}")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"success": true}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	p, err := env.db.GetProjectByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedProject := &datastore.Project{ID: 3, Name: "new-name", Fullname: "new-fullname"}
	if p.ID != wantedProject.ID || p.Name != wantedProject.Name || p.Fullname != wantedProject.Fullname {
		t.Errorf("expected %#v, got %#v", wantedProject, p)
	}
}

func TestCanPutProjectsOneHandlerAsOperatorWithJustName(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/projects/3", `{"name": "new-name"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsOneHandler), "/projects/{id}")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"success": true}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	p, err := env.db.GetProjectByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedProject := &datastore.Project{ID: 3, Name: "new-name", Fullname: "project 3"}
	if p.ID != wantedProject.ID || p.Name != wantedProject.Name || p.Fullname != wantedProject.Fullname {
		t.Errorf("expected %#v, got %#v", wantedProject, p)
	}
}

func TestCanPutProjectsOneHandlerAsOperatorWithJustFullname(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/projects/3", `{"fullname": "new-fullname"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsOneHandler), "/projects/{id}")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"success": true}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	p, err := env.db.GetProjectByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedProject := &datastore.Project{ID: 3, Name: "prj3", Fullname: "new-fullname"}
	if p.ID != wantedProject.ID || p.Name != wantedProject.Name || p.Fullname != wantedProject.Fullname {
		t.Errorf("expected %#v, got %#v", wantedProject, p)
	}
}

func TestCannotPutProjectsOneHandlerAsViewer(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/projects/3", `{"name": "new-name", "fullname": "new-fullname"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.projectsOneHandler), "/projects/{id}")
	hu.ConfirmAccessDenied(t, rec)

	// and verify state of database now
	p, err := env.db.GetProjectByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedProject := &datastore.Project{ID: 3, Name: "prj3", Fullname: "project 3"}
	if p.ID != wantedProject.ID || p.Name != wantedProject.Name || p.Fullname != wantedProject.Fullname {
		t.Errorf("expected %#v, got %#v", wantedProject, p)
	}
}
