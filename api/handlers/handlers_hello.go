// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"fmt"
	"net/http"
)

func (env *Env) helloHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// we only take GET requests
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	fmt.Fprintf(w, `{"success": true, "message": "hello"}`)
}
