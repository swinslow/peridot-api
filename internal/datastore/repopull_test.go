// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetAllRepoPullsForOneRepoBranch(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	t11 := time.Date(2019, 5, 2, 13, 53, 41, 671764, time.UTC)
	t15 := time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC)
	t16 := time.Date(2019, 5, 5, 12, 0, 0, 0, time.UTC)

	c11 := "0123456789012345678901234567890123456789"
	c15 := "4567890123456789012345678901234567890123"
	c16 := "8901234567890123456789012345678901234567"

	spdxID11 := "SPDXRef-xyzzy-11"
	spdxID15 := "SPDXRef-xyzzy-15"
	spdxID16 := "SPDXRef-xyzzy-16"

	sentRows := sqlmock.NewRows([]string{"id", "repo_id", "branch", "pulled_at", "commit", "tag", "spdx_id"}).
		AddRow(11, 3, "dev-1.1", t11, c11, "", spdxID11).
		AddRow(15, 3, "dev-1.1", t15, c15, "v1.1-rc0", spdxID15).
		AddRow(16, 3, "dev-1.1", t16, c16, "v1.1-rc1", spdxID16)
	mock.ExpectQuery(`SELECT id, repo_id, branch, pulled_at, commit, tag, spdx_id FROM repo_pulls WHERE repo_id = \$1 AND branch = \$2 ORDER BY id`).
		WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllRepoPullsForRepoBranch(3, "dev-1.1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if len(gotRows) != 3 {
		t.Fatalf("expected len %d, got %d", 3, len(gotRows))
	}
	rp0 := gotRows[0]
	if rp0.ID != 11 {
		t.Errorf("expected %v, got %v", 11, rp0.ID)
	}
	if rp0.RepoID != 3 {
		t.Errorf("expected %v, got %v", 3, rp0.RepoID)
	}
	if rp0.Branch != "dev-1.1" {
		t.Errorf("expected %v, got %v", "dev-1.1", rp0.Branch)
	}
	if rp0.PulledAt != t11 {
		t.Errorf("expected %v, got %v", t11, rp0.PulledAt)
	}
	if rp0.Commit != c11 {
		t.Errorf("expected %v, got %v", c11, rp0.Commit)
	}
	if rp0.Tag != "" {
		t.Errorf("expected %v, got %v", "", rp0.Tag)
	}
	if rp0.SPDXID != spdxID11 {
		t.Errorf("expected %v, got %v", spdxID11, rp0.SPDXID)
	}
	rp2 := gotRows[2]
	if rp2.ID != 16 {
		t.Errorf("expected %v, got %v", 16, rp2.ID)
	}
	if rp2.RepoID != 3 {
		t.Errorf("expected %v, got %v", 3, rp2.RepoID)
	}
	if rp2.Branch != "dev-1.1" {
		t.Errorf("expected %v, got %v", "dev-1.1", rp2.Branch)
	}
	if rp2.PulledAt != t16 {
		t.Errorf("expected %v, got %v", t16, rp2.PulledAt)
	}
	if rp2.Commit != c16 {
		t.Errorf("expected %v, got %v", c16, rp2.Commit)
	}
	if rp2.Tag != "v1.1-rc1" {
		t.Errorf("expected %v, got %v", "v1.1-rc1", rp2.Tag)
	}
	if rp2.SPDXID != spdxID16 {
		t.Errorf("expected %v, got %v", spdxID16, rp2.SPDXID)
	}
}

func TestShouldGetRepoPullByID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	t15 := time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC)
	c15 := "4567890123456789012345678901234567890123"
	spdxID15 := "SPDXRef-xyzzy-15"

	sentRows := sqlmock.NewRows([]string{"id", "repo_id", "branch", "pulled_at", "commit", "tag", "spdx_id"}).
		AddRow(15, 3, "dev-1.1", t15, c15, "v1.1-rc0", spdxID15)
	mock.ExpectQuery(`[SELECT id, repo_id, branch, pulled_at, commit, tag, spdx_id FROM repo_pulls WHERE id = \$1]`).
		WithArgs(15).
		WillReturnRows(sentRows)

	// run the tested function
	rp, err := db.GetRepoPullByID(15)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if rp.ID != 15 {
		t.Errorf("expected %v, got %v", 15, rp.ID)
	}
	if rp.RepoID != 3 {
		t.Errorf("expected %v, got %v", 3, rp.RepoID)
	}
	if rp.Branch != "dev-1.1" {
		t.Errorf("expected %v, got %v", "dev-1.1", rp.Branch)
	}
	if rp.PulledAt != t15 {
		t.Errorf("expected %v, got %v", t15, rp.PulledAt)
	}
	if rp.Commit != c15 {
		t.Errorf("expected %v, got %v", c15, rp.Commit)
	}
	if rp.Tag != "v1.1-rc0" {
		t.Errorf("expected %v, got %v", "v1.1-rc0", rp.Tag)
	}
	if rp.SPDXID != spdxID15 {
		t.Errorf("expected %v, got %v", spdxID15, rp.SPDXID)
	}
}

func TestShouldFailGetRepoPullByIDForUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	mock.ExpectQuery(`[SELECT id, repo_id, branch, pulled_at, commit, tag, spdx_id FROM repo_pulls WHERE id = \$1]`).
		WithArgs(413).
		WillReturnRows(sqlmock.NewRows([]string{}))

	// run the tested function
	rp, err := db.GetRepoPullByID(413)
	if rp != nil {
		t.Fatalf("expected nil repo pull, got %v", rp)
	}
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldAddRepoPull(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	t15 := time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC)
	c15 := "4567890123456789012345678901234567890123"
	spdxID15 := "SPDXRef-xyzzy-15"

	regexStmt := `[INSERT INTO repo_pulls(repo_id, branch, pulled_at, commit, tag, spdx_id) VALUES (\$1, \$2, \$3, \$4, \$5, \$6) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO repo_pulls"
	mock.ExpectQuery(stmt).
		WithArgs(15, "master", t15, c15, "v1.15-rc0", spdxID15).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(36))

	// run the tested function
	rpID, err := db.AddRepoPull(15, "master", t15, c15, "v1.15-rc0", spdxID15)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if rpID != 36 {
		t.Errorf("expected %v, got %v", 36, rpID)
	}
}

func TestShouldFailAddRepoPullWithUnknownRepoBranch(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	t0 := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	c0 := "4567890123456789012345678901234567890123"
	spdxID0 := "SPDXRef-oops"

	regexStmt := `[INSERT INTO repo_pulls(repo_id, branch, pulled_at, commit, tag, spdx_id) VALUES (\$1, \$2, \$3, \$4, \$5, \$6) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO repo_pulls"
	mock.ExpectQuery(stmt).
		WithArgs(413, "unknown-branch", t0, c0, "", spdxID0).
		WillReturnError(fmt.Errorf("pq: insert or update on table \"repo_pulls\" violates foreign key constraint \"repo_pulls_repo_id_fkey\""))

	// run the tested function
	_, err = db.AddRepoPull(413, "unknown-branch", t0, c0, "", spdxID0)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldDeleteRepoPull(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM repo_pulls WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM repo_pulls"
	mock.ExpectExec(stmt).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.DeleteRepoPull(1)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailDeleteRepoPullWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM repo_pulls WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM repo_pulls"
	mock.ExpectExec(stmt).
		WithArgs(413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function
	err = db.DeleteRepoPull(413)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// ===== JSON marshalling and unmarshalling =====
func TestCanMarshalRepoPullToJSON(t *testing.T) {
	rp := &RepoPull{
		ID:       17,
		RepoID:   5,
		Branch:   "master",
		PulledAt: time.Date(2019, 5, 2, 13, 53, 41, 0, time.UTC),
		Commit:   "0123456789012345678901234567890123456789",
		Tag:      "v1.12-rc3",
		SPDXID:   "SPDXRef-xyzzy-5",
	}

	js, err := json.Marshal(rp)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// read back in as empty interface to check values
	// should be a map whose keys are strings, values are empty interface values
	// per https://blog.golang.org/json-and-go
	var mapGot interface{}
	err = json.Unmarshal(js, &mapGot)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	mGot := mapGot.(map[string]interface{})

	// check for expected values
	if float64(rp.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(rp.ID), mGot["id"].(float64))
	}
	if float64(rp.RepoID) != mGot["repo_id"].(float64) {
		t.Errorf("expected %v, got %v", float64(rp.RepoID), mGot["repo_id"].(float64))
	}
	if rp.Branch != mGot["branch"].(string) {
		t.Errorf("expected %v, got %v", rp.Branch, mGot["branch"].(string))
	}
	if rp.PulledAt.Format(time.RFC3339) != mGot["pulled_at"].(string) {
		t.Errorf("expected %v, got %v", rp.PulledAt.Format(time.RFC3339), mGot["pulled_at"].(string))
	}
	if rp.Commit != mGot["commit"].(string) {
		t.Errorf("expected %v, got %v", rp.Commit, mGot["commit"].(string))
	}
	if rp.Tag != mGot["tag"].(string) {
		t.Errorf("expected %v, got %v", rp.Tag, mGot["tag"].(string))
	}
	if rp.SPDXID != mGot["spdx_id"].(string) {
		t.Errorf("expected %v, got %v", rp.SPDXID, mGot["spdx_id"].(string))
	}
}

func TestCanUnmarshalRepoPullFromJSON(t *testing.T) {
	rp := &RepoPull{}
	js := []byte(`{"id":17, "repo_id":1, "branch":"dev", "pulled_at":"2019-01-02T15:04:05Z", "commit":"4567890123456789012345678901234567890123", "tag":"t7", "spdx_id":"SPDXRef-xyzzy-17"}`)

	err := json.Unmarshal(js, rp)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if rp.ID != 17 {
		t.Errorf("expected %v, got %v", 17, rp.ID)
	}
	if rp.RepoID != 1 {
		t.Errorf("expected %v, got %v", 1, rp.RepoID)
	}
	if rp.Branch != "dev" {
		t.Errorf("expected %v, got %v", "dev", rp.Branch)
	}
	if rp.PulledAt.Format(time.RFC3339) != "2019-01-02T15:04:05Z" {
		t.Errorf("expected %v, got %v", "2019-01-02T15:04:05Z", rp.PulledAt.Format(time.RFC3339))
	}
	if rp.Commit != "4567890123456789012345678901234567890123" {
		t.Errorf("expected %v, got %v", "4567890123456789012345678901234567890123", rp.Commit)
	}
	if rp.Tag != "t7" {
		t.Errorf("expected %v, got %v", "t7", rp.Tag)
	}
	if rp.SPDXID != "SPDXRef-xyzzy-17" {
		t.Errorf("expected %v, got %v", "SPDXRef-xyzzy-17", rp.SPDXID)
	}

}

func TestCannotUnmarshalRepoPullWithNegativeIDFromJSON(t *testing.T) {
	rp := &RepoPull{}
	js := []byte(`{"id":-9283, "repo_id":1, "branch":"dev", "pulled_at":"2019-01-02T15:04:05Z", "commit":"4567890123456789012345678901234567890123", "tag":"t7", "spdx_id":"SPDXRef-xyzzy-17"}`)

	err := json.Unmarshal(js, rp)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
