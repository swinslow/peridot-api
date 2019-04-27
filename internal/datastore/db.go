// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"

	// postgres driver
	_ "github.com/lib/pq"
)

// DB holds the actual database/sql object as well as its related
// database statements.
type DB struct {
	sqldb *sql.DB
}

// NewDB opens and returns an initialized DB object.
func NewDB(srcName string) (*DB, error) {
	sqldb, err := sql.Open("postgres", srcName)
	if err != nil {
		return nil, err
	}
	if err = sqldb.Ping(); err != nil {
		return nil, err
	}

	db := &DB{sqldb: sqldb}
	return db, nil
}
