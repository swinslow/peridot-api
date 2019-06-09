// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import "os"

// CreateTableUsersAndAddInitialAdminUser creates the users table
// if it does not already exist. Also, if there are not yet any
// users, AND the environment variable INITIALADMINGITHUB is set,
// then it creates an initial admin user with ID 1 and the Github
// user name specified in that variable.
func (db *DB) CreateTableUsersAndAddInitialAdminUser() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS obsidian.users (
			id INTEGER NOT NULL PRIMARY KEY,
			github TEXT NOT NULL,
			name TEXT NOT NULL,
			access_level INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// if there are no users yet, and if INITIALADMINGITHUB env var
	// is also set, we'll create an initial administrative user
	// with ID 1
	users, err := db.GetAllUsers()
	if err == nil && len(users) == 0 {
		INITIALADMINGITHUB := os.Getenv("INITIALADMINGITHUB")
		if INITIALADMINGITHUB != "" {
			err = db.AddUser(1, "Admin", INITIALADMINGITHUB, AccessAdmin)
		}
	}
	return err
}

// CreateTableProjects creates the projects table if it
// does not already exist.
func (db *DB) CreateTableProjects() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS obsidian.projects (
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
		CREATE TABLE IF NOT EXISTS obsidian.subprojects (
			id SERIAL PRIMARY KEY,
			project_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			fullname TEXT NOT NULL,
			FOREIGN KEY (project_id) REFERENCES obsidian.projects (id)
		)
	`)
	return err
}

// CreateTableRepos creates the repos table if it does
// not already exist.
func (db *DB) CreateTableRepos() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS obsidian.repos (
			id SERIAL PRIMARY KEY,
			subproject_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			FOREIGN KEY (subproject_id) REFERENCES obsidian.subprojects (id)
		)
	`)
	return err
}

// CreateTableRepoBranches creates the repo_branches table
// if it does not already exist.
func (db *DB) CreateTableRepoBranches() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS obsidian.repo_branches (
			repo_id INTEGER,
			branch TEXT,
			PRIMARY KEY (repo_id, branch),
			FOREIGN KEY (repo_id) REFERENCES obsidian.repos (id)
		)
	`)
	return err
}

// CreateTableRepoPulls creates the repo_pulls table if it
// does not already exist.
func (db *DB) CreateTableRepoPulls() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS obsidian.repo_pulls (
			id SERIAL PRIMARY KEY,
			repo_id INTEGER NOT NULL,
			branch TEXT NOT NULL,
			pulled_at TIMESTAMP WITH TIME ZONE,
			commit TEXT,
			tag TEXT,
			spdx_id TEXT,
			FOREIGN KEY (repo_id, branch) REFERENCES obsidian.repo_branches (repo_id, branch)
		)
	`)
	return err
}

// CreateTableFileHashes creates the file_hashes table if it
// does not already exist.
func (db *DB) CreateTableFileHashes() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS obsidian.file_hashes (
			id SERIAL PRIMARY KEY,
			hash_s256 TEXT,
			hash_s1 TEXT
		)
	`)
	return err
}

// CreateTableFileInstances creates the file_instances table if it
// does not already exist.
func (db *DB) CreateTableFileInstances() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS obsidian.file_instances (
			id SERIAL PRIMARY KEY,
			repopull_id INTEGER NOT NULL,
			filehash_id INTEGER NOT NULL,
			path TEXT NOT NULL,
			FOREIGN KEY (repopull_id) REFERENCES obsidian.repo_pulls (id),
			FOREIGN KEY (filehash_id) REFERENCES obsidian.file_hashes (id)
		)
	`)
	return err
}
