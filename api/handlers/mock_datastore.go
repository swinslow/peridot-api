// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"fmt"
	"time"

	"github.com/swinslow/peridot-api/internal/datastore"
)

type mockDB struct {
	mockUsers        []*datastore.User
	mockProjects     []*datastore.Project
	mockSubprojects  []*datastore.Subproject
	mockRepos        []*datastore.Repo
	mockRepoBranches []*datastore.RepoBranch
}

// createMockDB creates mock values for the handler tests to use.
func createMockDB() *mockDB {
	mdb := &mockDB{}

	mdb.mockUsers = []*datastore.User{
		{ID: 1, Name: "Admin", Github: "admin", AccessLevel: datastore.AccessAdmin},
		{ID: 2, Name: "Operator", Github: "operator", AccessLevel: datastore.AccessOperator},
		{ID: 3, Name: "Commenter", Github: "commenter", AccessLevel: datastore.AccessCommenter},
		{ID: 4, Name: "Viewer", Github: "viewer", AccessLevel: datastore.AccessViewer},
		{ID: 10, Name: "Disabled", Github: "disabled", AccessLevel: datastore.AccessDisabled},
	}

	mdb.mockProjects = []*datastore.Project{
		{ID: 1, Name: "prj1", Fullname: "project 1"},
		{ID: 2, Name: "prj2", Fullname: "project 2"},
		{ID: 3, Name: "prj3", Fullname: "project 3"},
	}

	mdb.mockSubprojects = []*datastore.Subproject{
		{ID: 1, ProjectID: 3, Name: "subprj1", Fullname: "subproject 1"},
		{ID: 2, ProjectID: 1, Name: "subprj2", Fullname: "subproject 2"},
		{ID: 3, ProjectID: 1, Name: "subprj3", Fullname: "subproject 3"},
		{ID: 4, ProjectID: 1, Name: "subprj4", Fullname: "subproject 4"},
	}

	mdb.mockRepos = []*datastore.Repo{
		{ID: 1, SubprojectID: 2, Name: "repo1", Address: "https://example.com/repo1.git"},
		{ID: 2, SubprojectID: 4, Name: "repo2", Address: "https://example.com/repo2.git"},
		{ID: 3, SubprojectID: 4, Name: "repo3", Address: "https://example.com/repo3.git"},
		{ID: 4, SubprojectID: 4, Name: "repo4", Address: "https://example.com/repo4.git"},
	}

	mdb.mockRepoBranches = []*datastore.RepoBranch{
		{RepoID: 2, Branch: "master"},
		{RepoID: 2, Branch: "alpha"},
		{RepoID: 4, Branch: "master"},
		{RepoID: 4, Branch: "dev"},
		{RepoID: 2, Branch: "beta"},
		{RepoID: 1, Branch: "master"},
	}

	return mdb
}

// ===== Administrative actions =====
// ResetDB drops the current schema and initializes a new one.
// NOTE that if the initial Github user is not defined in an
// environment variable, the new DB will not have an admin user!
func (mdb *mockDB) ResetDB() error {
	// reset to just admin user
	mdb.mockUsers = []*datastore.User{
		&datastore.User{ID: 1, Name: "Admin", Github: "admin", AccessLevel: datastore.AccessAdmin},
	}
	return nil
}

// ===== Users =====

// GetAllUsers returns a slice of all users in the database.
func (mdb *mockDB) GetAllUsers() ([]*datastore.User, error) {
	return mdb.mockUsers, nil
}

// GetUserByID returns the User with the given user ID, or nil
// and an error if not found.
func (mdb *mockDB) GetUserByID(id uint32) (*datastore.User, error) {
	for _, user := range mdb.mockUsers {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("User not found with ID %d", id)
}

// GetUserByGithub returns the User with the given Github user
// name, or nil and an error if not found.
func (mdb *mockDB) GetUserByGithub(github string) (*datastore.User, error) {
	for _, user := range mdb.mockUsers {
		if user.Github == github {
			return user, nil
		}
	}
	return nil, fmt.Errorf("User not found with Github username %s", github)
}

// AddUser adds a new User with the given user ID, name, github
// user name, and access level. It returns nil on success or an
// error if failing.
func (mdb *mockDB) AddUser(id uint32, name string, github string, accessLevel datastore.UserAccessLevel) error {
	for _, u := range mdb.mockUsers {
		if u.ID == id {
			return fmt.Errorf("User with ID %d already exists in database", id)
		}
	}
	user := &datastore.User{
		ID:          id,
		Name:        name,
		Github:      github,
		AccessLevel: accessLevel,
	}

	mdb.mockUsers = append(mdb.mockUsers, user)
	return nil
}

// UpdateUser updates an existing User with the given ID,
// changing to the specified username, Github ID and and access
// level. It returns nil on success or an error if failing.
func (mdb *mockDB) UpdateUser(id uint32, newName string, newGithub string, newAccessLevel datastore.UserAccessLevel) error {
	for _, user := range mdb.mockUsers {
		if user.ID == id {
			user.Name = newName
			user.Github = newGithub
			user.AccessLevel = newAccessLevel
			return nil
		}
	}
	return fmt.Errorf("User not found with ID %d", id)
}

// UpdateUserNameOnly updates an existing User with the given ID,
// changing to the specified username. It returns nil on success
// or an error if failing.
func (mdb *mockDB) UpdateUserNameOnly(id uint32, newName string) error {
	for _, user := range mdb.mockUsers {
		if user.ID == id {
			user.Name = newName
			return nil
		}
	}
	return fmt.Errorf("User not found with ID %d", id)
}

// ===== Projects =====

// GetAllProjects returns a slice of all projects in the database.
func (mdb *mockDB) GetAllProjects() ([]*datastore.Project, error) {
	return mdb.mockProjects, nil
}

// GetProjectByID returns the Project with the given ID, or nil
// and an error if not found.
func (mdb *mockDB) GetProjectByID(id uint32) (*datastore.Project, error) {
	for _, prj := range mdb.mockProjects {
		if prj.ID == id {
			return prj, nil
		}
	}
	return nil, fmt.Errorf("Project not found with ID %d", id)
}

// AddProject adds a new Project with the given short name and
// full name. It returns the new project's ID on success or an
// error if failing.
func (mdb *mockDB) AddProject(name string, fullname string) (uint32, error) {
	// get max mock project ID
	var maxID uint32
	for _, p := range mdb.mockProjects {
		if p.Name == name {
			return 0, fmt.Errorf("Project with name %s already exists in database", name)
		}
		if p.ID > maxID {
			maxID = p.ID
		}
	}

	newID := maxID + 1
	prj := &datastore.Project{
		ID:       newID,
		Name:     name,
		Fullname: fullname,
	}

	mdb.mockProjects = append(mdb.mockProjects, prj)
	return newID, nil
}

// UpdateProject updates an existing Project with the given ID,
// changing to the specified short name and full name. If an
// empty string is passed, the existing value will remain
// unchanged. It returns nil on success or an error if failing.
func (mdb *mockDB) UpdateProject(id uint32, newName string, newFullname string) error {
	for _, p := range mdb.mockProjects {
		if p.ID == id {
			if newName != "" {
				p.Name = newName
			}
			if newFullname != "" {
				p.Fullname = newFullname
			}

			return nil
		}
	}
	return fmt.Errorf("Project not found with ID %d", id)
}

// DeleteProject deletes an existing Project with the given ID.
// It returns nil on success or an error if failing.
func (mdb *mockDB) DeleteProject(id uint32) error {
	found := false
	newMockProjects := []*datastore.Project{}
	for _, p := range mdb.mockProjects {
		if p.ID == id {
			found = true
		} else {
			newMockProjects = append(newMockProjects, p)
		}
	}
	if found {
		mdb.mockProjects = newMockProjects
		// and cascade delete any subprojects under this project
		for _, sp := range mdb.mockSubprojects {
			if sp.ProjectID == id {
				err := mdb.DeleteSubproject(sp.ID)
				if err != nil {
					return fmt.Errorf("Error with cascade delete of subprojects: %v", err)
				}
			}
		}
		return nil
	}
	return fmt.Errorf("Project not found with ID %d", id)
}

// ===== Subprojects =====

// GetAllSubprojects returns a slice of all subprojects in the
// database.
func (mdb *mockDB) GetAllSubprojects() ([]*datastore.Subproject, error) {
	return mdb.mockSubprojects, nil
}

// GetAllSubprojectsForProjectID returns a slice of all
// subprojects in the database for the given project ID.
func (mdb *mockDB) GetAllSubprojectsForProjectID(projectID uint32) ([]*datastore.Subproject, error) {
	subps := []*datastore.Subproject{}
	for _, subp := range mdb.mockSubprojects {
		if subp.ProjectID == projectID {
			subps = append(subps, subp)
		}
	}
	return subps, nil
}

// GetSubprojectByID returns the Subproject with the given ID, or nil
// and an error if not found.
func (mdb *mockDB) GetSubprojectByID(id uint32) (*datastore.Subproject, error) {
	for _, subp := range mdb.mockSubprojects {
		if subp.ID == id {
			return subp, nil
		}
	}
	return nil, fmt.Errorf("Subproject not found with ID %d", id)
}

// AddSubproject adds a new subproject with the given short
// name and full name, referencing the designated Project. It
// returns the new subproject's ID on success or an error if
// failing.
func (mdb *mockDB) AddSubproject(projectID uint32, name string, fullname string) (uint32, error) {
	// make sure project ID is valid
	_, err := mdb.GetProjectByID(projectID)
	if err != nil {
		return 0, fmt.Errorf("Project not found with ID %d", projectID)
	}

	// get max mock subproject ID
	var maxID uint32
	for _, sp := range mdb.mockSubprojects {
		if sp.Name == name && sp.ProjectID == projectID {
			return 0, fmt.Errorf("Subproject with name %s for project %d already exists in database", name, projectID)
		}
		if sp.ID > maxID {
			maxID = sp.ID
		}
	}

	newID := maxID + 1
	subp := &datastore.Subproject{
		ID:        newID,
		ProjectID: projectID,
		Name:      name,
		Fullname:  fullname,
	}

	mdb.mockSubprojects = append(mdb.mockSubprojects, subp)
	return newID, nil
}

// UpdateSubproject updates an existing Subproject with the
// given ID, changing to the specified short name and full
// name. If an empty string is passed, the existing value will
// remain unchanged. It returns nil on success or an error if
// failing.
func (mdb *mockDB) UpdateSubproject(id uint32, newName string, newFullname string) error {
	for _, sp := range mdb.mockSubprojects {
		if sp.ID == id {
			if newName != "" {
				sp.Name = newName
			}
			if newFullname != "" {
				sp.Fullname = newFullname
			}

			return nil
		}
	}
	return fmt.Errorf("Subproject not found with ID %d", id)
}

// UpdateSubprojectProjectID updates an existing Subproject
// with the given ID, changing its corresponding Project ID.
// It returns nil on success or an error if failing.
func (mdb *mockDB) UpdateSubprojectProjectID(id uint32, newProjectID uint32) error {
	// make sure project ID is valid
	_, err := mdb.GetProjectByID(newProjectID)
	if err != nil {
		return fmt.Errorf("Project not found with ID %d", newProjectID)
	}

	for _, sp := range mdb.mockSubprojects {
		if sp.ID == id {
			sp.ProjectID = newProjectID
			return nil
		}
	}
	return fmt.Errorf("Subproject not found with ID %d", id)
}

// DeleteSubproject deletes an existing Subproject with the
// given ID. It returns nil on success or an error if failing.
func (mdb *mockDB) DeleteSubproject(id uint32) error {
	found := false
	newMockSubprojects := []*datastore.Subproject{}
	for _, sp := range mdb.mockSubprojects {
		if sp.ID == id {
			found = true
		} else {
			newMockSubprojects = append(newMockSubprojects, sp)
		}
	}
	if found {
		mdb.mockSubprojects = newMockSubprojects
		return nil
	}
	return fmt.Errorf("Subproject not found with ID %d", id)
}

// ===== Repos =====

// GetAllRepos returns a slice of all repos in the database.
func (mdb *mockDB) GetAllRepos() ([]*datastore.Repo, error) {
	return mdb.mockRepos, nil
}

// GetAllReposForSubprojectID returns a slice of all repos in
// the database for the given subproject ID.
func (mdb *mockDB) GetAllReposForSubprojectID(subprojectID uint32) ([]*datastore.Repo, error) {
	repos := []*datastore.Repo{}
	for _, repo := range mdb.mockRepos {
		if repo.SubprojectID == subprojectID {
			repos = append(repos, repo)
		}
	}
	return repos, nil
}

// GetRepoByID returns the Repo with the given ID, or nil
// and an error if not found.
func (mdb *mockDB) GetRepoByID(id uint32) (*datastore.Repo, error) {
	for _, repo := range mdb.mockRepos {
		if repo.ID == id {
			return repo, nil
		}
	}
	return nil, fmt.Errorf("Repo not found with ID %d", id)
}

// AddRepo adds a new repo with the given name and address,
// referencing the designated Subproject. It returns the new
// repo's ID on success or an error if failing.
func (mdb *mockDB) AddRepo(subprojectID uint32, name string, address string) (uint32, error) {
	// make sure subproject ID is valid
	_, err := mdb.GetSubprojectByID(subprojectID)
	if err != nil {
		return 0, fmt.Errorf("Subproject not found with ID %d", subprojectID)
	}

	// get max mock repo ID
	var maxID uint32
	for _, repo := range mdb.mockRepos {
		if repo.Name == name && repo.SubprojectID == subprojectID {
			return 0, fmt.Errorf("Repo with name %s for subproject %d already exists in database", name, subprojectID)
		}
		if repo.ID > maxID {
			maxID = repo.ID
		}
	}

	newID := maxID + 1
	repo := &datastore.Repo{
		ID:           newID,
		SubprojectID: subprojectID,
		Name:         name,
		Address:      address,
	}

	mdb.mockRepos = append(mdb.mockRepos, repo)
	return newID, nil
}

// UpdateRepo updates an existing Repo with the given ID,
// changing to the specified name and address. If an empty
// string is passed, the existing value will remain unchanged.
// It returns nil on success or an error if failing.
func (mdb *mockDB) UpdateRepo(id uint32, newName string, newAddress string) error {
	for _, repo := range mdb.mockRepos {
		if repo.ID == id {
			if newName != "" {
				repo.Name = newName
			}
			if newAddress != "" {
				repo.Address = newAddress
			}

			return nil
		}
	}
	return fmt.Errorf("Repo not found with ID %d", id)
}

// UpdateRepoSubprojectID updates an existing Repo with the
// given ID, changing its corresponding Subproject ID.
// It returns nil on success or an error if failing.
func (mdb *mockDB) UpdateRepoSubprojectID(id uint32, newSubprojectID uint32) error {
	// make sure subproject ID is valid
	_, err := mdb.GetSubprojectByID(newSubprojectID)
	if err != nil {
		return fmt.Errorf("Subproject not found with ID %d", newSubprojectID)
	}

	for _, repo := range mdb.mockRepos {
		if repo.ID == id {
			repo.SubprojectID = newSubprojectID
			return nil
		}
	}
	return fmt.Errorf("Repo not found with ID %d", id)
}

// DeleteRepo deletes an existing Repo with the given ID.
// It returns nil on success or an error if failing.
func (mdb *mockDB) DeleteRepo(id uint32) error {
	found := false
	newMockRepos := []*datastore.Repo{}
	for _, repo := range mdb.mockRepos {
		if repo.ID == id {
			found = true
		} else {
			newMockRepos = append(newMockRepos, repo)
		}
	}
	if found {
		mdb.mockRepos = newMockRepos
		return nil
	}
	return fmt.Errorf("Repo not found with ID %d", id)
}

// ===== RepoBranches =====

// GetAllRepoBranchesForRepoID returns a slice of all repo
// branches in the database for the given Repo ID.
func (mdb *mockDB) GetAllRepoBranchesForRepoID(repoID uint32) ([]*datastore.RepoBranch, error) {
	rbs := []*datastore.RepoBranch{}
	for _, rb := range mdb.mockRepoBranches {
		if rb.RepoID == repoID {
			rbs = append(rbs, rb)
		}
	}
	return rbs, nil
}

// AddRepoBranch adds a new repo branch as specified,
// referencing the designated Repo. It returns nil on
// success or an error if failing.
func (mdb *mockDB) AddRepoBranch(repoID uint32, branch string) error {
	// make sure repo ID is valid
	_, err := mdb.GetRepoByID(repoID)
	if err != nil {
		return fmt.Errorf("Repo not found with ID %d", repoID)
	}

	// see if branch is already present for this repo
	for _, rb := range mdb.mockRepoBranches {
		if rb.RepoID == repoID && rb.Branch == branch {
			return fmt.Errorf("Branch %s for repo ID %d already exists in database", branch, repoID)
		}
	}

	rb := &datastore.RepoBranch{
		RepoID: repoID,
		Branch: branch,
	}

	mdb.mockRepoBranches = append(mdb.mockRepoBranches, rb)
	return nil
}

// DeleteRepoBranch deletes an existing RepoBranch with
// the given branch name for the given repo ID.
// It returns nil on success or an error if failing.
func (mdb *mockDB) DeleteRepoBranch(repoID uint32, branch string) error {
	found := false
	newMockRepoBranches := []*datastore.RepoBranch{}
	for _, rb := range mdb.mockRepoBranches {
		if rb.RepoID == repoID && rb.Branch == branch {
			found = true
		} else {
			newMockRepoBranches = append(newMockRepoBranches, rb)
		}
	}
	if found {
		mdb.mockRepoBranches = newMockRepoBranches
		return nil
	}
	return fmt.Errorf("Branch %s not found for repo ID %d", branch, repoID)
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
