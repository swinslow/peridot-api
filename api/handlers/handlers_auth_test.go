// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCanGetAuthLoginHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/auth/login", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	env := getTestEnv()
	http.HandlerFunc(env.authLoginHandler).ServeHTTP(rec, req)

	// check that we got a 307 (redirect)
	if 307 != rec.Code {
		t.Errorf("Expected %d, got %d", 307, rec.Code)
	}
}

func TestCannotPostAuthLoginHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/hello", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	env := getTestEnv()
	http.HandlerFunc(env.authLoginHandler).ServeHTTP(rec, req)

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}
}

func TestCannotPutAuthLoginHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/hello", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	env := getTestEnv()
	http.HandlerFunc(env.authLoginHandler).ServeHTTP(rec, req)

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}
}

func TestCannotDeleteAuthLoginHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/hello", nil)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	env := getTestEnv()
	http.HandlerFunc(env.authLoginHandler).ServeHTTP(rec, req)

	// check that we got a 405
	if 405 != rec.Code {
		t.Errorf("Expected %d, got %d", 405, rec.Code)
	}
}
