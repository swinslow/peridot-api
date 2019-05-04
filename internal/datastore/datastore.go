// Package datastore defines the database and in-memory models for all
// data in obsidian.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package datastore

// Datastore defines the interface to be implemented by models
// for database tables, using either a backing database (production)
// or mocks (test).
type Datastore interface {
	// ===== Users =====
	// GetAllUsers returns a slice of all users in the database.
	GetAllUsers() ([]*User, error)
	// GetUserByID returns the User with the given user ID, or nil
	// and an error if not found.
	GetUserByID(id uint32) (*User, error)
	// GetUserByEmail returns the User with the given email, or nil
	// and an error if not found.
	GetUserByEmail(email string) (*User, error)
	// AddUser adds a new User with the given user ID, name, email,
	// and access level. It returns nil on success or an error if
	// failing.
	AddUser(id uint32, name string, email string, accessLevel UserAccessLevel) error

	// ===== Projects =====
	// GetAllProjects returns a slice of all projects in the database.
	GetAllProjects() ([]*Project, error)
	// AddProject adds a new Project with the given short name and
	// full name. It returns the new project's ID on success or an
	// error if failing.
	AddProject(name string, fullname string) (uint32, error)
	// UpdateProject updates an existing Project with the given ID,
	// changing to the specified short name and full name. If an
	// empty string is passed, the existing value will remain
	// unchanged. It returns nil on success or an error if failing.
	UpdateProject(id uint32, newName string, newFullname string) error
	// DeleteProject deletes an existing Project with the given ID.
	// It returns nil on success or an error if failing.
	DeleteProject(id uint32) error
}
