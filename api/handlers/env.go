// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"github.com/swinslow/obsidian-api/internal/datastore"
)

// Env is the environment for the web handlers.
type Env struct {
	db datastore.Datastore
	//	jwtSecretKey string
}

// SetupEnv sets up systems (such as the data store) and variables
// (such as the JWT signing key) that are used across web requests.
func SetupEnv() (*Env, error) {
	// set up datastore
	db, err := datastore.NewDB("host=db sslmode=disable dbname=dev user=postgres-dev")
	if err != nil {
		return nil, err
	}

	err = datastore.InitNewDB(db)
	if err != nil {
		return nil, err
	}

	// set up JWT secret key (from environment)
	// JWTSECRETKEY := os.Getenv("JWTSECRETKEY")
	// if JWTSECRETKEY == "" {
	// 	return nil, fmt.Errorf("No secret key found; set environment variable JWTSECRETKEY before starting")
	// }

	env := &Env{
		db: db,
		//		jwtSecretKey: JWTSECRETKEY,
	}
	return env, nil
}
