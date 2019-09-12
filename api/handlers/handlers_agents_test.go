// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"testing"

	hu "github.com/swinslow/peridot-api/test/handlerutils"
	"github.com/swinslow/peridot-db/pkg/datastore"
)

// ===== GET /agents =====

func TestCanGetAgentsHandler(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/agents", ``, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsHandler), "/agents")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"agents": [
		{"id": 1, "name":"idsearcher", "is_active":true, "address":"localhost", "port":9001, "is_codereader":true, "is_spdxreader":false, "is_codewriter":false, "is_spdxwriter":true},
		{"id": 2, "name":"attributer", "is_active":true, "address":"localhost", "port":9002, "is_codereader":false, "is_spdxreader":true, "is_codewriter":true, "is_spdxwriter":false},
		{"id": 3, "name":"broken-agent", "is_active":false, "address":"example.com", "port":9003, "is_codereader":true, "is_spdxreader":false, "is_codewriter":true, "is_spdxwriter":true},
		{"id": 4, "name":"getter-github", "is_active":true, "address":"localhost", "port":9004, "is_codereader":false, "is_spdxreader":false, "is_codewriter":true, "is_spdxwriter":false},
		{"id": 5, "name":"analyze-godeps", "is_active":true, "address":"localhost", "port":9005, "is_codereader":true, "is_spdxreader":true, "is_codewriter":true, "is_spdxwriter":true},
		{"id": 6, "name":"decider", "is_active":true, "address":"localhost", "port":9006, "is_codereader":false, "is_spdxreader":true, "is_codewriter":false, "is_spdxwriter":true}
	]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetAgentsHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/agents", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsHandler), "/agents")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/agents", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsHandler), "/agents")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== POST /agents =====

func TestCanPostAgentsHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/agents", `{"name":"agent1", "is_active":true, "address":"https://example.com/agents", "port":8090, "is_codereader":true, "is_spdxreader":true, "is_codewriter":false, "is_spdxwriter":false}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsHandler), "/agents")
	hu.ConfirmCreatedResponse(t, rec)

	wanted := `{"id": 7}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	agents, err := env.db.GetAllAgents()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(agents) != 7 {
		t.Errorf("expected %d, got %d", 7, len(agents))
	}
	newAgent, err := env.db.GetAgentByID(7)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	wantedAgent := &datastore.Agent{ID: 7, Name: "agent1", IsActive: true, Address: "https://example.com/agents", Port: 8090, IsCodeReader: true, IsSpdxReader: true, IsCodeWriter: false, IsSpdxWriter: false}
	if newAgent.ID != wantedAgent.ID || newAgent.Name != wantedAgent.Name || newAgent.Address != wantedAgent.Address {
		t.Errorf("expected %#v, got %#v", wantedAgent, newAgent)
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
