// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	hu "github.com/swinslow/obsidian-api/test/handlerutils"
)

func TestCanGetHelloHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	//env := Env{db: db, jwtSecretKey: "keyForTesting"}
	env := Env{db: db}

	http.HandlerFunc(env.helloHandler).ServeHTTP(rec, req)

	// check that we got a 200 (OK)
	if 200 != rec.Code {
		t.Errorf("Expected %d, got %d", 200, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that the correct JSON strings were returned
	wanted := `{"message": "hello", "success": true}`
	got := hu.GetBody(t, rec)
	hu.CheckMatch(t, wanted, got, true)
}

func TestCannotPostHelloHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/hello", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	//env := Env{db: db, jwtSecretKey: "keyForTesting"}
	env := Env{db: db}

	http.HandlerFunc(env.helloHandler).ServeHTTP(rec, req)

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}

	// and check that the hello message was not returned
	got := hu.GetBody(t, rec)
	if strings.Contains(string(got), "hello") {
		t.Errorf("Expected hello to be absent, got %s", string(got))
	}
}

func TestCannotPutHelloHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/hello", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	//env := Env{db: db, jwtSecretKey: "keyForTesting"}
	env := Env{db: db}

	http.HandlerFunc(env.helloHandler).ServeHTTP(rec, req)

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}

	// and check that the hello message was not returned
	got := hu.GetBody(t, rec)
	if strings.Contains(string(got), "hello") {
		t.Errorf("Expected hello to be absent, got %s", string(got))
	}
}

func TestCannotDeleteHelloHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/hello", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	db := &mockDB{}
	//env := Env{db: db, jwtSecretKey: "keyForTesting"}
	env := Env{db: db}

	http.HandlerFunc(env.helloHandler).ServeHTTP(rec, req)

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}

	// and check that the hello message was not returned
	got := hu.GetBody(t, rec)
	if strings.Contains(string(got), "hello") {
		t.Errorf("Expected hello to be absent, got %s", string(got))
	}
}
