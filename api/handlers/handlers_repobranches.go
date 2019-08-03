// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/swinslow/peridot-api/internal/datastore"
)

// ========== HANDLER for /repos/{id}/branches

func (env *Env) repoBranchesSubHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// we only take GET or POST requests
	switch r.Method {
	case "GET":
		env.repoBranchesSubGetHelper(w, r)
	case "POST":
		env.repoBranchesSubPostHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) repoBranchesSubGetHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	// must be at least viewer
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access; get repo id from vars
	repoID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid repo ID"}`)
		return
	}

	// get repo branches from database
	branches, err := env.db.GetAllRepoBranchesForRepoID(repoID)
	if err != nil {
		fmt.Fprintf(w, `{"error": "Database retrieval error"}`)
		return
	}

	// want to return as a sorted array of strings with just
	// the branch names
	branchesArr := []string{}
	for _, b := range branches {
		branchesArr = append(branchesArr, b.Branch)
	}
	sort.Strings(branchesArr)

	// create map so we return a JSON object
	branchesMap := map[string][]string{}
	branchesMap["branches"] = branchesArr
	js, err := json.Marshal(branchesMap)
	if err != nil {
		fmt.Fprintf(w, `{"error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) repoBranchesSubPostHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	// must be at least operator
	user := extractUser(w, r, datastore.AccessOperator)
	if user == nil {
		return
	}

	// sufficient access; get repo id from vars
	repoID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid repo ID"}`)
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
	branch, ok := js["branch"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'branch'"}`)
		return
	}

	// add the new repo branch
	err = env.db.AddRepoBranch(repoID, branch.(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Unable to create repo branch"}`)
		return
	}

	// success!
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"branch": "%s"}`, branch.(string))
}
