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

// ===== GET /agents/3 =====

func TestCanGetAgentsOneHandlerAsViewer(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/agents/3", "", "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsOneHandler), "/agents/{id}")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"agent": {"id": 3, "name":"broken-agent", "is_active":false, "address":"example.com", "port":9003, "is_codereader":true, "is_spdxreader":false, "is_codewriter":true, "is_spdxwriter":true}}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetAgentsOneHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/agents/3", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsOneHandler), "/agents/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/agents/3", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsOneHandler), "/agents/{id}")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== PUT /agents/3 =====

func TestCanPutAgentsOneHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/agents/3", `{"is_active":true, "address":"agentHost", "port":8089, "is_codereader":true, "is_spdxreader":true, "is_codewriter":false, "is_spdxwriter":false}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsOneHandler), "/agents/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	agent, err := env.db.GetAgentByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedAgent := &datastore.Agent{ID: 3, Name: "broken-agent", IsActive: true, Address: "agentHost", Port: 8089, IsCodeReader: true, IsSpdxReader: true, IsCodeWriter: false, IsSpdxWriter: false}
	if agent.ID != wantedAgent.ID || agent.Name != wantedAgent.Name || agent.IsActive != wantedAgent.IsActive || agent.Address != wantedAgent.Address || agent.IsCodeReader != wantedAgent.IsCodeReader || agent.IsSpdxReader != wantedAgent.IsSpdxReader || agent.IsCodeWriter != wantedAgent.IsCodeWriter || agent.IsSpdxWriter != wantedAgent.IsSpdxWriter {
		t.Errorf("expected %#v, got %#v", wantedAgent, agent)
	}
}

func TestCanPutAgentsOneHandlerAsOperatorWithJustIsActive(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/agents/3", `{"is_active":true}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsOneHandler), "/agents/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	agent, err := env.db.GetAgentByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedAgent := &datastore.Agent{ID: 3, Name: "broken-agent", IsActive: true, Address: "example.com", Port: 9003, IsCodeReader: true, IsSpdxReader: false, IsCodeWriter: true, IsSpdxWriter: true}
	if agent.ID != wantedAgent.ID || agent.Name != wantedAgent.Name || agent.IsActive != wantedAgent.IsActive || agent.Address != wantedAgent.Address || agent.IsCodeReader != wantedAgent.IsCodeReader || agent.IsSpdxReader != wantedAgent.IsSpdxReader || agent.IsCodeWriter != wantedAgent.IsCodeWriter || agent.IsSpdxWriter != wantedAgent.IsSpdxWriter {
		t.Errorf("expected %#v, got %#v", wantedAgent, agent)
	}
}

func TestCanPutAgentsOneHandlerAsOperatorWithJustIsActiveFalse(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/agents/2", `{"is_active":false}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsOneHandler), "/agents/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	agent, err := env.db.GetAgentByID(2)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedAgent := &datastore.Agent{ID: 2, Name: "attributer", IsActive: false, Address: "localhost", Port: 9002, IsCodeReader: false, IsSpdxReader: true, IsCodeWriter: true, IsSpdxWriter: false}
	if agent.ID != wantedAgent.ID || agent.Name != wantedAgent.Name || agent.IsActive != wantedAgent.IsActive || agent.Address != wantedAgent.Address || agent.IsCodeReader != wantedAgent.IsCodeReader || agent.IsSpdxReader != wantedAgent.IsSpdxReader || agent.IsCodeWriter != wantedAgent.IsCodeWriter || agent.IsSpdxWriter != wantedAgent.IsSpdxWriter {
		t.Errorf("expected %#v, got %#v", wantedAgent, agent)
	}
}

func TestCanPutAgentsOneHandlerAsOperatorWithJustAddressAndPort(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/agents/3", `{"address":"localhost", "port":8089}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsOneHandler), "/agents/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	agent, err := env.db.GetAgentByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedAgent := &datastore.Agent{ID: 3, Name: "broken-agent", IsActive: false, Address: "localhost", Port: 8089, IsCodeReader: true, IsSpdxReader: false, IsCodeWriter: true, IsSpdxWriter: true}
	if agent.ID != wantedAgent.ID || agent.Name != wantedAgent.Name || agent.IsActive != wantedAgent.IsActive || agent.Address != wantedAgent.Address || agent.IsCodeReader != wantedAgent.IsCodeReader || agent.IsSpdxReader != wantedAgent.IsSpdxReader || agent.IsCodeWriter != wantedAgent.IsCodeWriter || agent.IsSpdxWriter != wantedAgent.IsSpdxWriter {
		t.Errorf("expected %#v, got %#v", wantedAgent, agent)
	}
}

func TestCanPutAgentsOneHandlerAsOperatorWithJustAbilities(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/agents/3", `{"is_codereader":true, "is_spdxreader":true, "is_codewriter":false, "is_spdxwriter":false}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsOneHandler), "/agents/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	agent, err := env.db.GetAgentByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedAgent := &datastore.Agent{ID: 3, Name: "broken-agent", IsActive: false, Address: "example.com", Port: 9003, IsCodeReader: true, IsSpdxReader: true, IsCodeWriter: false, IsSpdxWriter: false}
	if agent.ID != wantedAgent.ID || agent.Name != wantedAgent.Name || agent.IsActive != wantedAgent.IsActive || agent.Address != wantedAgent.Address || agent.IsCodeReader != wantedAgent.IsCodeReader || agent.IsSpdxReader != wantedAgent.IsSpdxReader || agent.IsCodeWriter != wantedAgent.IsCodeWriter || agent.IsSpdxWriter != wantedAgent.IsSpdxWriter {
		t.Errorf("expected %#v, got %#v", wantedAgent, agent)
	}
}

func TestCannotPutAgentsOneHandlerAsCommenter(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/agents/3", `{"is_active":true, "address":"agentHost", "port":8089, "is_codereader":true, "is_spdxreader":true, "is_codewriter":false, "is_spdxwriter":false}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsOneHandler), "/agents/{id}")
	hu.ConfirmAccessDenied(t, rec)

	// and verify state of database now
	agent, err := env.db.GetAgentByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedAgent := &datastore.Agent{ID: 3, Name: "broken-agent", IsActive: false, Address: "example.com", Port: 9003, IsCodeReader: true, IsSpdxReader: false, IsCodeWriter: true, IsSpdxWriter: true}
	if agent.ID != wantedAgent.ID || agent.Name != wantedAgent.Name || agent.IsActive != wantedAgent.IsActive || agent.Address != wantedAgent.Address || agent.IsCodeReader != wantedAgent.IsCodeReader || agent.IsSpdxReader != wantedAgent.IsSpdxReader || agent.IsCodeWriter != wantedAgent.IsCodeWriter || agent.IsSpdxWriter != wantedAgent.IsSpdxWriter {
		t.Errorf("expected %#v, got %#v", wantedAgent, agent)
	}
}

// ===== DELETE /agents/3 =====

func TestCanDeleteAgentsOneHandlerAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "DELETE", "/agents/3", ``, "admin")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsOneHandler), "/agents/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	agents, err := env.db.GetAllAgents()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(agents) != 5 {
		t.Errorf("expected %d, got %d", 5, len(agents))
	}
	agent, err := env.db.GetAgentByID(3)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil and %#v", agent)
	}
}

func TestCannotDeleteAgentsOneHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "DELETE", "/agents/3", ``, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.agentsOneHandler), "/agents/{id}")
	hu.ConfirmAccessDenied(t, rec)

	// and verify state of database has not changed
	agents, err := env.db.GetAllAgents()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(agents) != 6 {
		t.Errorf("expected %d, got %d", 6, len(agents))
	}
	agent, err := env.db.GetAgentByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedAgent := &datastore.Agent{ID: 3, Name: "broken-agent", IsActive: false, Address: "example.com", Port: 9003, IsCodeReader: true, IsSpdxReader: false, IsCodeWriter: true, IsSpdxWriter: true}
	if agent.ID != wantedAgent.ID || agent.IsActive != wantedAgent.IsActive || agent.Name != wantedAgent.Name || agent.Address != wantedAgent.Address || agent.Port != wantedAgent.Port || agent.IsCodeReader != wantedAgent.IsCodeReader || agent.IsSpdxReader != wantedAgent.IsSpdxReader || agent.IsCodeWriter != wantedAgent.IsCodeWriter || agent.IsSpdxWriter != wantedAgent.IsSpdxWriter {
		t.Errorf("expected %#v, got %#v", wantedAgent, agent)
	}
}
