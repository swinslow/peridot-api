// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/swinslow/obsidian-api/internal/datastore"
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

	// we only take GET requests
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

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
