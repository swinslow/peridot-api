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
	// and subprojects within a project
	router.HandleFunc("/projects/{id:[0-9]+}/subprojects", env.validateTokenMiddleware(env.subprojectsSubHandler)).Methods("GET", "POST")

	// /subprojects -- subproject data
	router.HandleFunc("/subprojects", env.validateTokenMiddleware(env.subprojectsHandler)).Methods("GET", "POST")
	router.HandleFunc("/subprojects/{id:[0-9]+}", env.validateTokenMiddleware(env.subprojectsOneHandler)).Methods("GET", "PUT", "DELETE")
	// and repos within a subproject
	router.HandleFunc("/subprojects/{id:[0-9]+}/repos", env.validateTokenMiddleware(env.reposSubHandler)).Methods("GET", "POST")

	// /repos -- repo data
	router.HandleFunc("/repos", env.validateTokenMiddleware(env.reposHandler)).Methods("GET", "POST")
	router.HandleFunc("/repos/{id:[0-9]+}", env.validateTokenMiddleware(env.reposOneHandler)).Methods("GET", "PUT", "DELETE")
	// and a repo's branches
	router.HandleFunc("/repos/{id:[0-9]+}/branches", env.validateTokenMiddleware(env.repoBranchesSubHandler)).Methods("GET", "POST")
	// and a specific branch, to POST a new repo pull
	// FIXME the pattern here does not sync with the various rules for branch naming in git
	router.HandleFunc(`/repos/{id:[0-9]+}/branches/{branch:[0-9a-zA-Z_\-\.]+}`, env.validateTokenMiddleware(env.repoPullsSubHandler)).Methods("GET", "POST")

	// /repopulls -- repo pull data
	router.HandleFunc("/repopulls/{id:[0-9]+}", env.validateTokenMiddleware(env.repoPullsOneHandler)).Methods("GET", "DELETE")

	// /agents -- registered peridot agents
	router.HandleFunc("/agents", env.validateTokenMiddleware(env.agentsHandler)).Methods("GET", "POST")
}
