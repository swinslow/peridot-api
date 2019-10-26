// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"testing"

	hu "github.com/swinslow/peridot-api/test/handlerutils"
	"github.com/swinslow/peridot-db/pkg/datastore"
)

// ===== POST /repopulls/3/jobs =====

func TestCanPostJobsHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/jobs", `{"name":"job1", "is_active":true, "address":"https://example.com/jobs", "port":8090, "is_codereader":true, "is_spdxreader":true, "is_codewriter":false, "is_spdxwriter":false}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.jobsHandler), "/jobs")
	hu.ConfirmCreatedResponse(t, rec)

	wanted := `{"id": 7}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	jobs, err := env.db.GetAllJobs()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(jobs) != 7 {
		t.Errorf("expected %d, got %d", 7, len(jobs))
	}
	newJob, err := env.db.GetJobByID(7)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	// =========== FIXME STOPPED HERE MID-TEST
	wantedJob := &datastore.Job{ID: 7, Name: "agent1", IsActive: true, Address: "https://example.com/agents", Port: 8090, IsCodeReader: true, IsSpdxReader: true, IsCodeWriter: false, IsSpdxWriter: false}
	if newAgent.ID != wantedJob.ID || newAgent.Name != wantedJob.Name || newAgent.Address != wantedJob.Address {
		t.Errorf("expected %#v, got %#v", wantedJob, newAgent)
	}
}

func TestCannotPostAgentsHandlerAsOtherUser(t *testing.T) {
	// as commenter
	rec, req, env := setupTestEnv(t, "POST", "/agents", `{"name":"agent1", "is_active":true, "address":"https://example.com/agents", "port":8090, "is_codereader":true, "is_spdxreader":true, "is_codewriter":false, "is_spdxwriter":false}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsHandler), "/agents")
	hu.ConfirmAccessDenied(t, rec)

	// as viewer
	rec, req, env = setupTestEnv(t, "POST", "/agents", `{"name":"agent1", "is_active":true, "address":"https://example.com/agents", "port":8090, "is_codereader":true, "is_spdxreader":true, "is_codewriter":false, "is_spdxwriter":false}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsHandler), "/agents")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCannotPostAgentHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/agents", `{"name":"agent1", "is_active":true, "address":"https://example.com/agents", "port":8090, "is_codereader":true, "is_spdxreader":true, "is_codewriter":false, "is_spdxwriter":false}`, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsHandler), "/agents")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "POST", "/agents", `{"name":"agent1", "is_active":true, "address":"https://example.com/agents", "port":8090, "is_codereader":true, "is_spdxreader":true, "is_codewriter":false, "is_spdxwriter":false}`, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsHandler), "/agents")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}
