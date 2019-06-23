// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"testing"

	"github.com/swinslow/peridot-api/internal/datastore"
	hu "github.com/swinslow/peridot-api/test/handlerutils"
)

// ===== GET /subprojects =====

func TestCanGetSubprojectsHandler(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/subprojects", ``, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.subprojectsHandler), "/subprojects")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"subprojects": [{"id": 1, "project_id": 3, "name": "subprj1", "fullname": "subproject 1"}, {"id": 2, "project_id": 1, "name": "subprj2", "fullname": "subproject 2"}, {"id": 3, "project_id": 1, "name": "subprj3", "fullname": "subproject 3"}, {"id": 4, "project_id": 1, "name": "subprj4", "fullname": "subproject 4"}]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetSubprojectsHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/subprojects", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.subprojectsHandler), "/subprojects")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/subprojects", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.subprojectsHandler), "/subprojects")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== POST /subprojects =====

func TestCanPostSubprojectsHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/subprojects", `{"project_id": 2, "name": "subprj5", "fullname": "subproject 5"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.subprojectsHandler), "/subprojects")
	hu.ConfirmCreatedResponse(t, rec)

	wanted := `{"success": true, "id": 5}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	subprojects, err := env.db.GetAllSubprojects()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(subprojects) != 5 {
		t.Errorf("expected %d, got %d", 5, len(subprojects))
	}
	newSubproject, err := env.db.GetSubprojectByID(5)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	wantedSubproject := &datastore.Subproject{ID: 5, ProjectID: 2, Name: "subprj5", Fullname: "subproject 5"}
	if newSubproject.ID != wantedSubproject.ID || newSubproject.Name != wantedSubproject.Name || newSubproject.Fullname != wantedSubproject.Fullname {
		t.Errorf("expected %#v, got %#v", wantedSubproject, newSubproject)
	}
}

func TestCannotPostSubprojectsHandlerAsOtherUser(t *testing.T) {
	// as commenter
	rec, req, env := setupTestEnv(t, "POST", "/subprojects", `{"project_id": 2, "name": "subprj5", "fullname": "subproject 5"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.subprojectsHandler), "/subprojects")
	hu.ConfirmAccessDenied(t, rec)

	// as viewer
	rec, req, env = setupTestEnv(t, "POST", "/subprojects", `{"project_id": 2, "name": "subprj5", "fullname": "subproject 5"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.subprojectsHandler), "/subprojects")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCannotPostSubprojectsHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/subprojects", `{"project_id": 2, "name": "subprj5", "fullname": "subproject 5"}`, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.subprojectsHandler), "/subprojects")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "POST", "/subprojects", `{"project_id": 2, "name": "subprj5", "fullname": "subproject 5"}`, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.subprojectsHandler), "/subprojects")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}
