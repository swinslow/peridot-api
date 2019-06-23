// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/swinslow/peridot-api/internal/datastore"
)

func (env *Env) adminDBHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// we only take POST requests
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// get user and check access level
	user := extractUser(w, r, datastore.AccessAdmin)
	if user == nil {
		return
	}

	// sufficient access; check command
	js := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&js)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Invalid JSON request"}`)
		return
	}
	command, ok := js["command"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "No command specified"}`)
		return
	}
	switch command {
	case "resetDB":
		err = env.db.ResetDB()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "Unable to reset database"}`)
			return
		}
		fmt.Fprintf(w, `{"success": true}`)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Unknown command '%s'"}`, command)
		return
	}

}
