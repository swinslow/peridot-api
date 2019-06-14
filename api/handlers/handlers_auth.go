// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/swinslow/peridot-api/internal/auth"
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

// NOTE that authGithubCallbackHandler DOES NOT return JSON.
// Instead, it returns HTML with JavaScript that is intended
// to save the JWT in the browser's local storage, and then
// redirect back to the webapp root location. An API user
// would need to obtain the JWT through the webapp and then
// use it in other peridot API calls as needed.
func (env *Env) authGithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ghUser, err := auth.ValidateGithub(r, env.oauthConf, env.oauthState)
	if err != nil {
		// FIXME return HTML with error message + redirect
		fmt.Fprintf(w, "<html>\n<body>\n<p>Error: Couldn't validate GitHub credentials</p>\n</body>\n</html>\n")
	}

	// code was valid and we have the user name
	// encode it into a JWT
	tkn, err := auth.EncodeToken(env.jwtSecretKey, ghUser)
	if err != nil {
		// FIXME return HTML with error message + redirect
		fmt.Fprintf(w, "<html>\n<body>\n<p>Error: Couldn't create token</p>\n</body>\n</html>\n")
	}

	// return HTML with JWT + JS for localstorage + redirect
	fmt.Fprintf(w, "<html>\n<script>\nwindow.localStorage.setItem('apitoken', '%s');\nwindow.location.href = '/';\n</script>\n</html>\n", tkn)
}
