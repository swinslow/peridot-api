// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"time"

	"github.com/swinslow/obsidian-api/internal/datastore"
)

type mockDB struct {
	mockUsers []*datastore.User
}

// ===== Users =====

// GetAllUsers returns a slice of all users in the database.
func (mdb *mockDB) GetAllUsers() ([]*datastore.User, error) {
	return []*datastore.User{}, nil
}

// GetUserByID returns the User with the given user ID, or nil
// and an error if not found.
func (mdb *mockDB) GetUserByID(id uint32) (*datastore.User, error) {
	return nil, nil
}

// GetUserByEmail returns the User with the given email, or nil
// and an error if not found.
func (mdb *mockDB) GetUserByEmail(email string) (*datastore.User, error) {
	return nil, nil
}

// AddUser adds a new User with the given user ID, name, email,
// and access level. It returns nil on success or an error if
// failing.
func (mdb *mockDB) AddUser(id uint32, name string, email string, accessLevel datastore.UserAccessLevel) error {
	return nil
}

// ===== Projects =====

// GetAllProjects returns a slice of all projects in the database.
func (mdb *mockDB) GetAllProjects() ([]*datastore.Project, error) {
	return []*datastore.Project{}, nil
}

// GetProjectByID returns the Project with the given ID, or nil
// and an error if not found.
func (mdb *mockDB) GetProjectByID(id uint32) (*datastore.Project, error) {
	return nil, nil
}

// AddProject adds a new Project with the given short name and
// full name. It returns the new project's ID on success or an
// error if failing.
func (mdb *mockDB) AddProject(name string, fullname string) (uint32, error) {
	return 0, nil
}

// UpdateProject updates an existing Project with the given ID,
// changing to the specified short name and full name. If an
// empty string is passed, the existing value will remain
// unchanged. It returns nil on success or an error if failing.
func (mdb *mockDB) UpdateProject(id uint32, newName string, newFullname string) error {
	return nil
}

// DeleteProject deletes an existing Project with the given ID.
// It returns nil on success or an error if failing.
func (mdb *mockDB) DeleteProject(id uint32) error {
	return nil
}

// ===== Subprojects =====

// GetAllSubprojects returns a slice of all subprojects in the
// database.
func (mdb *mockDB) GetAllSubprojects() ([]*datastore.Subproject, error) {
	return []*datastore.Subproject{}, nil
}

// GetAllSubprojectsForProjectID returns a slice of all
// subprojects in the database for the given project ID.
func (mdb *mockDB) GetAllSubprojectsForProjectID(projectID uint32) ([]*datastore.Subproject, error) {
	return []*datastore.Subproject{}, nil
}

// GetSubprojectByID returns the Subproject with the given ID, or nil
// and an error if not found.
func (mdb *mockDB) GetSubprojectByID(id uint32) (*datastore.Subproject, error) {
	return nil, nil
}

// AddSubproject adds a new subproject with the given short
// name and full name, referencing the designated Project. It
// returns the new subproject's ID on success or an error if
// failing.
func (mdb *mockDB) AddSubproject(projectID uint32, name string, fullname string) (uint32, error) {
	return 0, nil
}

// UpdateSubproject updates an existing Subproject with the
// given ID, changing to the specified short name and full
// name. If an empty string is passed, the existing value will
// remain unchanged. It returns nil on success or an error if
// failing.
func (mdb *mockDB) UpdateSubproject(id uint32, newName string, newFullname string) error {
	return nil
}

// UpdateSubprojectProjectID updates an existing Subproject
// with the given ID, changing its corresponding Project ID.
// It returns nil on success or an error if failing.
func (mdb *mockDB) UpdateSubprojectProjectID(id uint32, newProjectID uint32) error {
	return nil
}

// DeleteSubproject deletes an existing Subproject with the
// given ID. It returns nil on success or an error if failing.
func (mdb *mockDB) DeleteSubproject(id uint32) error {
	return nil
}

// ===== Repos =====

// GetAllRepos returns a slice of all repos in the database.
func (mdb *mockDB) GetAllRepos() ([]*datastore.Repo, error) {
	return []*datastore.Repo{}, nil
}

// GetAllReposForSubprojectID returns a slice of all repos in
// the database for the given subproject ID.
func (mdb *mockDB) GetAllReposForSubprojectID(subprojectID uint32) ([]*datastore.Repo, error) {
	return []*datastore.Repo{}, nil
}

// GetRepoByID returns the Repo with the given ID, or nil
// and an error if not found.
func (mdb *mockDB) GetRepoByID(id uint32) (*datastore.Repo, error) {
	return nil, nil
}

// AddRepo adds a new repo with the given name and address,
// referencing the designated Subproject. It returns the new
// repo's ID on success or an error if failing.
func (mdb *mockDB) AddRepo(subprojectID uint32, name string, address string) (uint32, error) {
	return 0, nil
}

// UpdateRepo updates an existing Repo with the given ID,
// changing to the specified name and address. If an empty
// string is passed, the existing value will remain unchanged.
// It returns nil on success or an error if failing.
func (mdb *mockDB) UpdateRepo(id uint32, newName string, newAddress string) error {
	return nil
}

// UpdateRepoSubprojectID updates an existing Repo with the
// given ID, changing its corresponding Subproject ID.
// It returns nil on success or an error if failing.
func (mdb *mockDB) UpdateRepoSubprojectID(id uint32, newSubprojectID uint32) error {
	return nil
}

// DeleteRepo deletes an existing Repo with the given ID.
// It returns nil on success or an error if failing.
func (mdb *mockDB) DeleteRepo(id uint32) error {
	return nil
}

// ===== RepoBranches =====

// GetAllRepoBranchesForRepoID returns a slice of all repo
// branches in the database for the given Repo ID.
func (mdb *mockDB) GetAllRepoBranchesForRepoID(repoID uint32) ([]*datastore.RepoBranch, error) {
	return []*datastore.RepoBranch{}, nil
}

// AddRepoBranch adds a new repo branch as specified,
// referencing the designated Repo. It returns nil on
// success or an error if failing.
func (mdb *mockDB) AddRepoBranch(repoID uint32, branch string) error {
	return nil
}

// DeleteRepoBranch deletes an existing RepoBranch with
// the given branch name for the given repo ID.
// It returns nil on success or an error if failing.
func (mdb *mockDB) DeleteRepoBranch(repoID uint32, branch string) error {
	return nil
}

// ===== RepoPulls =====

// GetAllRepoPullsForRepoBranch returns a slice of all repo
// pulls in the database for the given Repo ID and branch.
func (mdb *mockDB) GetAllRepoPullsForRepoBranch(repoID uint32, branch string) ([]*datastore.RepoPull, error) {
	return []*datastore.RepoPull{}, nil
}

// GetRepoPullByID returns the RepoPull with the given ID,
// or nil and an error if not found.
func (mdb *mockDB) GetRepoPullByID(id uint32) (*datastore.RepoPull, error) {
	return nil, nil
}

// AddRepoPull adds a new repo pull as specified,
// referencing the designated Repo, branch and other data.
// It returns the new repo pull's ID on success or an
// error if failing.
func (mdb *mockDB) AddRepoPull(repoID uint32, branch string, pulledAt time.Time, commit string, tag string, spdxID string) (uint32, error) {
	return 0, nil
}

// DeleteRepoPull deletes an existing RepoPull with the
// given ID. It returns nil on success or an error if
// failing.
func (mdb *mockDB) DeleteRepoPull(id uint32) error {
	return nil
}

// ===== FileHashes =====

// GetFileHashByID returns the FileHash with the given ID,
// or nil and an error if not found.
func (mdb *mockDB) GetFileHashByID(id uint64) (*datastore.FileHash, error) {
	return nil, nil
}

// GetFileHashesByIDs returns a slice of FileHashes with
// the given IDs, or an empty slice if none are found.
// NOT CURRENTLY TESTED; NEED TO MODIFY FOR USING pq.Array
/*GetFileHashesByIDs(ids []uint64) ([]*FileHash, error)*/

// AddFileHash adds a new file hash as specified,
// requiring its SHA256 and SHA1 values. It returns the
// new file hash's ID on success or an error if failing.
func (mdb *mockDB) AddFileHash(sha256 string, sha1 string) (uint64, error) {
	return 0, nil
}

// FIXME will also want one to add a slice of file hashes
// FIXME all at once

// DeleteFileHash deletes an existing file hash with
// the given ID. It returns nil on success or an error if
// failing.
func (mdb *mockDB) DeleteFileHash(id uint64) error {
	return nil
}

// ===== FileInstancees =====

// GetFileInstanceByID returns the FileInstance with the given ID,
// or nil and an error if not found.
func (mdb *mockDB) GetFileInstanceByID(id uint64) (*datastore.FileInstance, error) {
	return nil, nil
}

// AddFileInstance adds a new file instance as specified,
// requiring its parent RepoPull ID and path within it,
// and the corresponding FileHash ID. It returns the new
// file instance's ID on success or an error if failing.
func (mdb *mockDB) AddFileInstance(repoPullID uint32, fileHashID uint64, path string) (uint64, error) {
	return 0, nil
}

// DeleteFileInstance deletes an existing file instance
// with the given ID. It returns nil on success or an
// if failing.
func (mdb *mockDB) DeleteFileInstance(id uint64) error {
	return nil
}
