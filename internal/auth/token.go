// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package auth

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
)

// EncodeToken takes the secret key for this server and
// Github user name, and creates and returns a JWT signed
// token for that user, or an error message on error.
func EncodeToken(jwtSecretKey string, github string) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"github": github,
	})
	tknString, err := tkn.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("error encoding login token")
	}

	return tknString, nil
}

// DecodeToken takes the secret key for this server and
// a JWT token, and tries to decode the Github user name
// from it. It returns the github user name if it can be
// parsed. It does NOT check the database to confirm
// whether this user exists or what their ID or access is.
func DecodeToken(jwtSecretKey string, tknRecv string) (string, error) {
	// decrypt and validate the token
	token, err := jwt.Parse(tknRecv, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Couldn't parse token method")
		}
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("error decoding login token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("error decoding login token")
	}
	github, ok := claims["github"].(string)
	if !ok {
		return "", fmt.Errorf("error decoding login token")
	}

	return github, nil
}
