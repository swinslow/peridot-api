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
