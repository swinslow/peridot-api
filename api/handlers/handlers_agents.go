// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/swinslow/peridot-db/pkg/datastore"
)

// ========== HANDLER for /agents

func (env *Env) agentsHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// we only take GET or POST requests
	switch r.Method {
	case "GET":
		env.agentsGetHelper(w, r)
	case "POST":
		env.agentsPostHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) agentsGetHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	// must be at least viewer
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access; get agents from database
	agents, err := env.db.GetAllAgents()
	if err != nil {
		fmt.Fprintf(w, `{"error": "Database retrieval error"}`)
		return
	}

	// create map so we return a JSON object
	agentsMap := map[string][]*datastore.Agent{}
	agentsMap["agents"] = agents
	js, err := json.Marshal(agentsMap)
	if err != nil {
		fmt.Fprintf(w, `{"error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) agentsPostHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	// must be at least operator
	user := extractUser(w, r, datastore.AccessOperator)
	if user == nil {
		return
	}

	// sufficient access; parse JSON request
	js := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&js)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Invalid JSON request"}`)
		return
	}

	// and extract data
	name, ok := js["name"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'name'"}`)
		return
	}
	isActive, ok := js["is_active"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'is_active'"}`)
		return
	}
	address, ok := js["address"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'address'"}`)
		return
	}
	portf, ok := js["port"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'port'"}`)
		return
	}
	isCodeReader, ok := js["is_codereader"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'is_codereader'"}`)
		return
	}
	isSpdxReader, ok := js["is_spdxreader"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'is_spdxreader'"}`)
		return
	}
	isCodeWriter, ok := js["is_codewriter"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'is_codewriter'"}`)
		return
	}
	isSpdxWriter, ok := js["is_spdxwriter"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'is_spdxwriter'"}`)
		return
	}

	// convert port float64 to int
	port := int(portf.(float64))

	// add the new agent
	newID, err := env.db.AddAgent(name.(string), isActive.(bool), address.(string), port, isCodeReader.(bool), isSpdxReader.(bool), isCodeWriter.(bool), isSpdxWriter.(bool))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Unable to create agent"}`)
		return
	}

	// success!
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id": %d}`, newID)
}

// ========== HANDLER for /agents/{id}

func (env *Env) agentsOneHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// check valid request types
	switch r.Method {
	case "GET":
		env.agentsOneGetHelper(w, r)
	case "PUT":
		env.agentsOnePutHelper(w, r)
	case "DELETE":
		env.agentsOneDeleteHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) agentsOneGetHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	agentID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid ID"}`)
		return
	}

	// get agent from database
	agent, err := env.db.GetAgentByID(agentID)
	if err != nil {
		fmt.Fprintf(w, `{"error": "Database retrieval error"}`)
		return
	}

	// create map so we return a JSON object
	jsData := struct {
		Agent *datastore.Agent `json:"agent"`
	}{Agent: agent}
	js, err := json.Marshal(jsData)
	if err != nil {
		fmt.Fprintf(w, `{"error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) agentsOnePutHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessOperator)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	agentID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid ID"}`)
		return
	}

	// get existing agent from database
	agent, err := env.db.GetAgentByID(agentID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "Unknown agent ID"}`)
		return
	}

	// parse JSON request
	js := map[string]interface{}{}
	err = json.NewDecoder(r.Body).Decode(&js)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Invalid JSON request"}`)
		return
	}

	// and extract data; if absent, use existing data
	flagStatus := false
	flagAbilities := false

	// updateable status vars
	newIsActive, ok := js["is_active"].(bool)
	if ok {
		flagStatus = true
	} else {
		newIsActive = agent.IsActive
	}
	newAddress, ok := js["address"].(string)
	if ok {
		flagStatus = true
	} else {
		newAddress = agent.Address
	}
	var newPort int
	newPortTmp, ok := js["port"].(float64)
	if ok {
		flagStatus = true
		newPort = int(newPortTmp)
	} else {
		newPort = agent.Port
	}

	// updateable ability vars
	newIsCodeReader, ok := js["is_codereader"].(bool)
	if ok {
		flagAbilities = true
	} else {
		newIsCodeReader = agent.IsCodeReader
	}
	newIsSpdxReader, ok := js["is_spdxreader"].(bool)
	if ok {
		flagAbilities = true
	} else {
		newIsSpdxReader = agent.IsSpdxReader
	}
	newIsCodeWriter, ok := js["is_codewriter"].(bool)
	if ok {
		flagAbilities = true
	} else {
		newIsCodeWriter = agent.IsCodeWriter
	}
	newIsSpdxWriter, ok := js["is_spdxwriter"].(bool)
	if ok {
		flagAbilities = true
	} else {
		newIsSpdxWriter = agent.IsSpdxWriter
	}

	if !flagStatus && !flagAbilities {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "No updateable values found in request"}`)
		return
	}

	// modify the repo data where applicable
	if flagStatus {
		err = env.db.UpdateAgentStatus(agentID, newIsActive, newAddress, newPort)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "Unable to update repo"}`)
			return
		}
	}
	if flagAbilities {
		err = env.db.UpdateAgentAbilities(agentID, newIsCodeReader, newIsSpdxReader, newIsCodeWriter, newIsSpdxWriter)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "Unable to update repo"}`)
			return
		}
	}

	// success!
	w.WriteHeader(http.StatusNoContent)
}

func (env *Env) agentsOneDeleteHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessAdmin)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	agentID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid ID"}`)
		return
	}

	// delete the agent
	err = env.db.DeleteAgent(agentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Unable to delete agent"}`)
		return
	}

	// success!
	w.WriteHeader(http.StatusNoContent)
}
