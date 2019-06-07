// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import "fmt"

// User describes a registered user of the platform.
type User struct {
	// ID is the unique ID for this user.
	ID uint32 `json:"id"`
	// Name is this user's name.
	Name string `json:"name"`
	// Email is this user's registered email address.
	Email string `json:"email"`
	// AccessLevel is this user's access level.
	AccessLevel UserAccessLevel `json:"access"`
}

// GetAllUsers returns a slice of all users in the database.
func (db *DB) GetAllUsers() ([]*User, error) {
	rows, err := db.sqldb.Query("SELECT id, email, name, access_level FROM obsidian.users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.AccessLevel)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID returns the User with the given user ID, or nil
// and an error if not found.
func (db *DB) GetUserByID(id uint32) (*User, error) {
	var user User
	var ualInt int
	err := db.sqldb.QueryRow("SELECT id, email, name, access_level FROM obsidian.users WHERE id = $1", id).
		Scan(&user.ID, &user.Email, &user.Name, &ualInt)
	if err != nil {
		return nil, err
	}

	// convert integer to UserAccessLevel
	ual, err := UserAccessLevelFromInt(ualInt)
	if err != nil {
		return nil, err
	}

	user.AccessLevel = ual
	return &user, nil
}

// GetUserByEmail returns the User with the given email, or nil
// and an error if not found.
func (db *DB) GetUserByEmail(email string) (*User, error) {
	var user User
	var ualInt int
	err := db.sqldb.QueryRow("SELECT id, email, name, access_level FROM obsidian.users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Name, &ualInt)
	if err != nil {
		return nil, err
	}

	// convert integer to UserAccessLevel
	ual, err := UserAccessLevelFromInt(ualInt)
	if err != nil {
		return nil, err
	}

	user.AccessLevel = ual
	return &user, nil
}

// AddUser adds a new User with the given user ID, name, email, and
// access level. It returns nil on success or an error if failing.
// Due to PostgreSQL limits on integer size, id must be less than 2147483647.
// It should typically be created via math/rand's Int31() function and then
// cast to uint32.
func (db *DB) AddUser(id uint32, name string, email string, accessLevel UserAccessLevel) error {
	var maxUserID uint32
	maxUserID = 2147483647

	if id > maxUserID {
		return fmt.Errorf("User id cannot be greater than %d; received %d", maxUserID, id)
	}

	ualInt := IntFromUserAccessLevel(accessLevel)

	// move out into one-time-prepared statement?
	stmt, err := db.sqldb.Prepare("INSERT INTO obsidian.users(id, email, name, access_level) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id, email, name, ualInt)
	if err != nil {
		return err
	}
	return nil
}
