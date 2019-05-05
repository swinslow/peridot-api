// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

// CreateTableUsers creates the users table if it does not
// already exist.
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

// CreateTableProjects creates the projects table if it
// does not already exist.
func (db *DB) CreateTableProjects() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			fullname TEXT NOT NULL
		)
	`)
	return err
}

// CreateTableSubprojects creates the subprojects table
// if it does not already exist.
func (db *DB) CreateTableSubprojects() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS subprojects (
			id SERIAL PRIMARY KEY,
			project_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			fullname TEXT NOT NULL,
			FOREIGN KEY (project_id) REFERENCES projects (id)
		)
	`)
	return err
}

// CreateTableRepos creates the repos table if it does
// not already exist.
func (db *DB) CreateTableRepos() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS repos (
			id SERIAL PRIMARY KEY,
			subproject_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			FOREIGN KEY (subproject_id) REFERENCES subprojects (id)
		)
	`)
	return err
}

// CreateTableRepoBranches creates the repo_branches table
// if it does not already exist.
func (db *DB) CreateTableRepoBranches() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS repo_branches (
			repo_id INTEGER,
			branch TEXT,
			PRIMARY KEY (repo_id, branch),
			FOREIGN KEY (repo_id) REFERENCES repos (id)
		)
	`)
	return err
}
