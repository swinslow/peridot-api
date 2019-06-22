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

// used to strip down the info that is returned to
// non-admin users
type limitedUser struct {
	ID     uint32 `json:"id"`
	Github string `json:"github"`
}

func (env *Env) usersHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// we only take GET or POST requests
	switch r.Method {
	case "GET":
		env.usersGetHelper(w, r)
	case "POST":
		env.usersPostHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) usersGetHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access; get users from database
	users, err := env.db.GetAllUsers()
	if err != nil {
		fmt.Fprintf(w, `{"error": "Database retrieval error"}`)
		return
	}

	// return different message depending whether the
	// logged-in user is admin or lesser
	if user.AccessLevel == datastore.AccessAdmin {
		// admin user just does full JSON marshalling
		// create map so we return a JSON object
		usersMap := map[string][]*datastore.User{}
		usersMap["users"] = users
		js, err := json.Marshal(usersMap)
		if err != nil {
			fmt.Fprintf(w, `{"error": "JSON marshalling error"}`)
			return
		}
		w.Write(js)
		return
	}

	// for non-admin recipient, need to return just id
	// and Github username
	ltdUsers := []*limitedUser{}
	for _, u := range users {
		lu := &limitedUser{ID: u.ID, Github: u.Github}
		ltdUsers = append(ltdUsers, lu)
	}

	// create map so we return a JSON object
	ltdUsersMap := map[string][]*limitedUser{}
	ltdUsersMap["users"] = ltdUsers

	// now write JSON for limited data
	js, err := json.Marshal(ltdUsersMap)
	if err != nil {
		fmt.Fprintf(w, `{"error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) usersPostHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessAdmin)
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
	ghUsername, ok := js["github"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Missing required value for 'github'"}`)
		return
	}
	access, ok := js["access"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Missing required value for 'access'"}`)
		return
	}

	// check that access value is valid
	ual, err := datastore.UserAccessLevelFromString(access)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Invalid value for 'access'"}`)
		return
	}

	// FIXME: we should be able to add a user without specifying
	// an ID. Since db.AddUser currently requires an ID, we'll
	// manually check the maximum existing user ID, and choose
	// the next highest.
	users, err := env.db.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"success": false, "error": "Error in user database"}`)
		return
	}
	var maxCurrentUserID uint32
	for _, u := range users {
		if u.ID > maxCurrentUserID {
			maxCurrentUserID = u.ID
		}
	}
	newID := maxCurrentUserID + 1

	// add the new user
	err = env.db.AddUser(newID, name, ghUsername, ual)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"success": false, "error": "Unable to create user"}`)
		return
	}

	// success!
	fmt.Fprintf(w, `{"success": true, "id": %d}`, newID)
}

func (env *Env) usersOneHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// check valid request types
	switch r.Method {
	case "GET":
		env.usersOneGetHelper(w, r)
	case "PUT":
		env.usersOnePutHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, PUT")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) usersOneGetHelper(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(w, `{"success": false, "error": "Missing or invalid user ID"}`)
		return
	}
	u, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Invalid user ID"}`)
		return
	}
	userID := uint32(u)

	// get user from database
	argUser, err := env.db.GetUserByID(userID)
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "error": "Database retrieval error"}`)
		return
	}

	// return different message depending whether the
	// logged-in user is admin / self or other
	if user.AccessLevel == datastore.AccessAdmin || user.ID == argUser.ID {
		// admin user just does full JSON marshalling
		// create map so we return a JSON object
		jsData := struct {
			Success bool            `json:"success"`
			User    *datastore.User `json:"user"`
		}{Success: true, User: argUser}
		js, err := json.Marshal(jsData)
		if err != nil {
			fmt.Fprintf(w, `{"success": false, "error": "JSON marshalling error"}`)
			return
		}
		w.Write(js)
		return
	}

	// for non-admin, non-self recipient, need to return just id
	// and Github username
	jsData := struct {
		Success bool         `json:"success"`
		LtdUser *limitedUser `json:"user"`
	}{Success: true, LtdUser: &limitedUser{ID: argUser.ID, Github: argUser.Github}}
	js, err := json.Marshal(jsData)
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) usersOnePutHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access generally, but maybe not for the requested user?
	// extract ID for request
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Missing or invalid user ID"}`)
		return
	}
	u, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false, "error": "Invalid user ID"}`)
		return
	}
	userID := uint32(u)

	// if not admin and not self, access will be denied
	if user.AccessLevel != datastore.AccessAdmin && user.ID != userID {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `{"success": false, "error": "Access denied"}`)
		return
	}

	// get existing user from database
	existingUser, err := env.db.GetUserByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"success": false, "error": "Unknown user ID"}`)
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
		newName = existingUser.Name
	}
	newGithub, ok := js["github"]
	if ok {
		// unless we're admin, access will be denied
		if user.AccessLevel != datastore.AccessAdmin {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, `{"success": false, "error": "Access denied"}`)
			return
		}
	} else {
		newGithub = existingUser.Github
	}
	var newUal datastore.UserAccessLevel
	newAccess, ok := js["access"]
	if ok {
		// unless we're admin, access will be denied
		if user.AccessLevel != datastore.AccessAdmin {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, `{"success": false, "error": "Access denied"}`)
			return
		}
		newUal, err = datastore.UserAccessLevelFromString(newAccess)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"success": false, "error": "Invalid value for 'access'"}`)
			return
		}
	} else {
		newUal = existingUser.AccessLevel
	}

	// modify the user data - if admin, for all; if not, name only
	if user.AccessLevel == datastore.AccessAdmin {
		err = env.db.UpdateUser(userID, newName, newGithub, newUal)
	} else {
		err = env.db.UpdateUserNameOnly(userID, newName)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"success": false, "error": "Unable to update user"}`)
		return
	}

	// success!
	fmt.Fprintf(w, `{"success": true}`)
}
