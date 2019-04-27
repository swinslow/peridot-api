// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

// CreateTableUsers creates the users table if it does not already exist.
func (db *DB) CreateTableUsers() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY,
			email TEXT NOT NULL,
			name TEXT NOT NULL,
			access_level INTEGER NOT NULL
		)
	`)
	return err
}
