// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/swinslow/peridot-api/internal/auth"
	"github.com/swinslow/peridot-db/pkg/datastore"
)

const (
	// ErrAuthBearer signifies that a valid Bearer token
	// was not included with the request headers.
	ErrAuthBearer = "Authorization header with valid Bearer token required"

	// ErrAuthGithub signifies that a Bearer token was
	// included with the request headers, but that the
	// Github user named in the token is not registered
	// with the peridot database.
	ErrAuthGithub = "Github user is not registered"

	// ErrAuthAccess signifies that a Bearer token was
	// included with the request headers, and that it
	// identifies a Github user registered with the
	// peridot database, but that user does not have
	// sufficient access rights for the requested action.
	ErrAuthAccess = "Access denied"
)

func sendAuthFail(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, `{"error": %s}`, ErrAuthBearer)
}

type userContextKey int

func (env *Env) validateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// look for and extract the token from header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendAuthFail(w)
			return
		}

		// check that the auth header has the expected format
		// e.g. Authorization: Bearer ....
		if !strings.HasPrefix(authHeader, "Bearer ") {
			sendAuthFail(w)
			return
		}
		remainder := strings.TrimPrefix(authHeader, "Bearer ")

		// decrypt and validate the token
		ghUsername, err := auth.DecodeToken(env.jwtSecretKey, remainder)
		if err != nil {
			sendAuthFail(w)
			return
		}

		// make sure this email also exists in the User database
		user, err := env.db.GetUserByGithub(ghUsername)
		if err != nil {
			user = &datastore.User{
				ID:          0,
				Github:      ghUsername,
				Name:        "",
				AccessLevel: datastore.AccessDisabled,
			}
		}

		// good to go! set context and move on
		ctx := r.Context()
		ctx = context.WithValue(ctx, userContextKey(0), user)
		next(w, r.WithContext(ctx))
	})
}

// pull out the user from context (after auth), and confirm
// they have AT LEAST the requested access level. If they
// don't, return a JSON "access denied" error. If they do
// have sufficient access, the calling handler can still
// take different actions based on their actual access
// level (e.g., send a different response to admins vs.
// normal users).
func extractUser(w http.ResponseWriter, r *http.Request, minLevel datastore.UserAccessLevel) *datastore.User {
	// pull User from context
	ctxCheck := r.Context().Value(userContextKey(0))
	if ctxCheck == nil {
		w.Header().Set("WWW-Authenticate", "Bearer")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%s"}`, ErrAuthBearer)
		return nil
	}
	user := ctxCheck.(*datastore.User)
	if user.ID == 0 {
		w.Header().Set("WWW-Authenticate", "Bearer")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%s"}`, ErrAuthGithub)
		return nil
	}

	// check minimum access required for this resource
	if user.AccessLevel < minLevel {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `{"error": "%s"}`, ErrAuthAccess)
		return nil
	}

	return user
}
