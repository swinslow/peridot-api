// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	hu "github.com/swinslow/peridot-api/test/handlerutils"
)

func TestCanGetAuthLoginHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/auth/login", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	env := getTestEnv()
	hu.ServeHandler(rec, req, http.HandlerFunc(env.authLoginHandler), "/auth/login")

	// check that we got a 307 (redirect)
	if 307 != rec.Code {
		t.Errorf("Expected %d, got %d", 307, rec.Code)
	}
}

func TestCannotPostAuthLoginHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/auth/login", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	env := getTestEnv()
	hu.ServeHandler(rec, req, http.HandlerFunc(env.authLoginHandler), "/auth/login")

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}
}

func TestCannotPutAuthLoginHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/auth/login", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	env := getTestEnv()
	hu.ServeHandler(rec, req, http.HandlerFunc(env.authLoginHandler), "/auth/login")

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}
}

func TestCannotDeleteAuthLoginHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/auth/login", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	env := getTestEnv()
	hu.ServeHandler(rec, req, http.HandlerFunc(env.authLoginHandler), "/auth/login")

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}
}
