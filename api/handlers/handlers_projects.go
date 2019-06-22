// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/swinslow/peridot-api/internal/datastore"
)

func (env *Env) projectsHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// we only take GET or POST requests
	switch r.Method {
	case "GET":
		env.projectsGetHelper(w, r)
	case "POST":
		env.projectsPostHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) projectsGetHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	// must be at least viewer
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access; get projects from database
	projects, err := env.db.GetAllProjects()
	if err != nil {
		fmt.Fprintf(w, `{"error": "Database retrieval error"}`)
		return
	}

	// create map so we return a JSON object
	projectsMap := map[string][]*datastore.Project{}
	projectsMap["projects"] = projects
	js, err := json.Marshal(projectsMap)
	if err != nil {
		fmt.Fprintf(w, `{"error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) projectsPostHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	// must be at least operator
	user := extractUser(w, r, datastore.AccessOperator)
	if user == nil {
		return
	}

	// sufficient access; parse JSON request
	js := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&js)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Invalid JSON request"}`)
		return
	}

	// and extract data
	name, ok := js["name"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Missing required value for 'name'"}`)
		return
	}
	fullname, ok := js["fullname"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Missing required value for 'fullname'"}`)
		return
	}

	// add the new project
	newID, err := env.db.AddProject(name, fullname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"success": false, "error": "Unable to create project"}`)
		return
	}

	// success!
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"success": true, "id": %d}`, newID)
}

func (env *Env) projectsOneHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// check valid request types
	switch r.Method {
	case "GET":
		env.projectsOneGetHelper(w, r)
	case "PUT":
		env.projectsOnePutHelper(w, r)
	case "DELETE":
		env.projectsOneDeleteHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) projectsOneGetHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Missing or invalid project ID"}`)
		return
	}
	p, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Invalid project ID"}`)
		return
	}
	projectID := uint32(p)

	// get project from database
	argProject, err := env.db.GetProjectByID(projectID)
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "error": "Database retrieval error"}`)
		return
	}

	// create map so we return a JSON object
	jsData := struct {
		Success bool               `json:"success"`
		Project *datastore.Project `json:"project"`
	}{Success: true, Project: argProject}
	js, err := json.Marshal(jsData)
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) projectsOnePutHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessOperator)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Missing or invalid project ID"}`)
		return
	}
	p, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Invalid project ID"}`)
		return
	}
	projectID := uint32(p)

	// get existing project from database
	project, err := env.db.GetProjectByID(projectID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"success": false, "error": "Unknown project ID"}`)
		return
	}

	// parse JSON request
	js := map[string]string{}
	err = json.NewDecoder(r.Body).Decode(&js)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Invalid JSON request"}`)
		return
	}

	// and extract data; if absent, use existing data
	newName, ok := js["name"]
	if !ok {
		newName = project.Name
	}
	newFullname, ok := js["fullname"]
	if !ok {
		newFullname = project.Fullname
	}

	// modify the project data
	err = env.db.UpdateProject(projectID, newName, newFullname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"success": false, "error": "Unable to update project"}`)
		return
	}

	// success!
	fmt.Fprintf(w, `{"success": true}`)
}

func (env *Env) projectsOneDeleteHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessAdmin)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Missing or invalid project ID"}`)
		return
	}
	p, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Invalid project ID"}`)
		return
	}
	projectID := uint32(p)

	// delete the project
	err = env.db.DeleteProject(projectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"success": false, "error": "Unable to delete project"}`)
		return
	}

	// success!
	fmt.Fprintf(w, `{"success": true}`)
}
