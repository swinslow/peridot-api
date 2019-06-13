// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package auth

import (
	"context"
	"fmt"
	"net/http"
	
	"github.com/google/go-github/v25/github"
	"golang.org/x/oauth2"
)

// ValidateGithub parses a (presumed) Github OAuth redirect
// request, and tries to use it to obtain the user's Github
// user name. It returns the user name if successful or an
// error if not.
func ValidateGithub(r *http.Request, oauthConf *oauth2.Config, oauthState string) (string, error) {
	ctx := context.Background()

	// first, check and confirm the state matches
	state := r.FormValue("state")
	if state != oauthState {
		return "", fmt.Errorf("invalid oauth state, expected %s, got %s", oauthState, state)
	}

	// now, check the code
	code := r.FormValue("code")
	token, err := oauthConf.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("could not access Github API: %v", err)
	}

	// finally, use the code to get user data
	oauthClient := oauthConf.Client(ctx, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return "", fmt.Errorf("could not get user's Github data: %v", err)
	}

	return *user.Login, nil
}