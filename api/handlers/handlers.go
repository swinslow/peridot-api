// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"github.com/gorilla/mux"
)

// RegisterHandlers registers the api handler endpoints with the
// specified router, for the given environment.
func (env *Env) RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/hello", env.helloHandler).Methods("GET")
	router.HandleFunc("/auth/login", env.authLoginHandler).Methods("GET")
}
