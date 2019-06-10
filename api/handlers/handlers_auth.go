// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"golang.org/x/oauth2"
)

func (env *Env) authLoginHandler(w http.ResponseWriter, r *http.Request) {
	// we only take GET requests
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	url := env.oauthConf.AuthCodeURL(env.oauthState, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
