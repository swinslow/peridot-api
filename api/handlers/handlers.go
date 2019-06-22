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

	// /auth -- authentication / OAuth flow
	router.HandleFunc("/auth/login", env.authLoginHandler).Methods("GET")
	router.HandleFunc("/auth/redirect", env.authGithubCallbackHandler).Methods("GET")

	// /admin -- administrative actions
	router.HandleFunc("/admin/db", env.validateTokenMiddleware(env.adminDBHandler)).Methods("POST")

	// /users -- user data
	router.HandleFunc("/users", env.validateTokenMiddleware(env.usersHandler)).Methods("GET", "POST")
	router.HandleFunc("/users/{id:[0-9]+}", env.validateTokenMiddleware(env.usersOneHandler)).Methods("GET", "PUT")

	// /projects -- project data
	router.HandleFunc("/projects", env.validateTokenMiddleware(env.projectsHandler)).Methods("GET", "POST")
	router.HandleFunc("/projects/{id:[0-9]+}", env.validateTokenMiddleware(env.projectsOneHandler)).Methods("GET", "PUT", "DELETE")
}
