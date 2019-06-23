// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/swinslow/peridot-api/internal/datastore"
)

// ========== HANDLER for /repos

func (env *Env) reposHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// we only take GET or POST requests
	switch r.Method {
	case "GET":
		env.reposGetHelper(w, r)
	case "POST":
		env.reposPostHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) reposGetHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	// must be at least viewer
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access; get repos from database
	repos, err := env.db.GetAllRepos()
	if err != nil {
		fmt.Fprintf(w, `{"error": "Database retrieval error"}`)
		return
	}

	// create map so we return a JSON object
	reposMap := map[string][]*datastore.Repo{}
	reposMap["repos"] = repos
	js, err := json.Marshal(reposMap)
	if err != nil {
		fmt.Fprintf(w, `{"error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) reposPostHelper(w http.ResponseWriter, r *http.Request) {
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
	subprojectIDf, ok := js["subproject_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'subproject_id'"}`)
		return
	}
	name, ok := js["name"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'name'"}`)
		return
	}
	address, ok := js["address"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'address'"}`)
		return
	}

	// convert subprojectID float64 to uint32
	subprojectID := uint32(subprojectIDf.(float64))

	// add the new repo
	newID, err := env.db.AddRepo(subprojectID, name.(string), address.(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Unable to create repo"}`)
		return
	}

	// success!
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id": %d}`, newID)
}

// ========== HANDLER for /subprojects/{id}/repos

func (env *Env) reposSubHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// we only take GET or POST requests
	switch r.Method {
	case "GET":
		env.reposSubGetHelper(w, r)
	case "POST":
		env.reposSubPostHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) reposSubGetHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	// must be at least viewer
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access; get subproject id from vars
	subprojectID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid subproject ID"}`)
		return
	}

	// get repos from database
	repos, err := env.db.GetAllReposForSubprojectID(subprojectID)
	if err != nil {
		fmt.Fprintf(w, `{"error": "Database retrieval error"}`)
		return
	}

	// create map so we return a JSON object
	reposMap := map[string][]*datastore.Repo{}
	reposMap["repos"] = repos
	js, err := json.Marshal(reposMap)
	if err != nil {
		fmt.Fprintf(w, `{"error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) reposSubPostHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	// must be at least operator
	user := extractUser(w, r, datastore.AccessOperator)
	if user == nil {
		return
	}

	// sufficient access; get subproject id from vars
	subprojectID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid subproject ID"}`)
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

	// and extract data
	name, ok := js["name"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'name'"}`)
		return
	}
	address, ok := js["address"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'address'"}`)
		return
	}

	// add the new repo
	newID, err := env.db.AddRepo(subprojectID, name.(string), address.(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Unable to create repo"}`)
		return
	}

	// success!
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id": %d}`, newID)
}

// ========== HANDLER for /repos/{id}

func (env *Env) reposOneHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// check valid request types
	switch r.Method {
	case "GET":
		env.reposOneGetHelper(w, r)
	case "PUT":
		env.reposOnePutHelper(w, r)
	case "DELETE":
		env.reposOneDeleteHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) reposOneGetHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	repoID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid ID"}`)
		return
	}

	// get repo from database
	repo, err := env.db.GetRepoByID(repoID)
	if err != nil {
		fmt.Fprintf(w, `{"error": "Database retrieval error"}`)
		return
	}

	// create map so we return a JSON object
	jsData := struct {
		Repo *datastore.Repo `json:"repo"`
	}{Repo: repo}
	js, err := json.Marshal(jsData)
	if err != nil {
		fmt.Fprintf(w, `{"error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) reposOnePutHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessOperator)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	repoID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid ID"}`)
		return
	}

	// get existing repo from database
	repo, err := env.db.GetRepoByID(repoID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "Unknown repo ID"}`)
		return
	}

	// parse JSON request
	js := map[string]string{}
	err = json.NewDecoder(r.Body).Decode(&js)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Invalid JSON request"}`)
		return
	}

	// and extract data; if absent, use existing data
	newName, ok := js["name"]
	if !ok {
		newName = repo.Name
	}
	newAddress, ok := js["address"]
	if !ok {
		newAddress = repo.Address
	}
	// NOTE: currently, cannot update the repo's project ID
	// using this API call.

	// modify the repo data
	err = env.db.UpdateRepo(repoID, newName, newAddress)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Unable to update repo"}`)
		return
	}

	// success!
	w.WriteHeader(http.StatusNoContent)
}

func (env *Env) reposOneDeleteHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessAdmin)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	repoID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid ID"}`)
		return
	}

	// delete the repo
	err = env.db.DeleteRepo(repoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Unable to delete repo"}`)
		return
	}

	// success!
	w.WriteHeader(http.StatusNoContent)
}
