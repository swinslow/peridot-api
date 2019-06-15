// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"github.com/gorilla/mux"
)

// RegisterHandlers registers the api handler endpoints with the
// specified router, for the given environment.
func (env *Env) RegisterHandlers(router *mux.Router) {
	// /hello -- ping and hello
	router.HandleFunc("/hello", env.helloHandler).Methods("GET")

	// /admin -- administrative actions
	router.HandleFunc("/admin/db", env.adminDBHandler).Methods("POST")

	// /auth -- authentication / OAuth flow
	router.HandleFunc("/auth/login", env.authLoginHandler).Methods("GET")
	router.HandleFunc("/auth/redirect", env.authGithubCallbackHandler).Methods("GET")

	// /users -- user data
	router.HandleFunc("/users", env.validateTokenMiddleware(env.usersHandler)).Methods("GET")
}
