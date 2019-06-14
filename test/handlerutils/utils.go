// Package handlerutils contains utility functions for testing the
// peridot API handlers.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package handlerutils

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/yudai/gojsondiff"
)

// CheckMatch compares a wanted string and a got byte slice containing
// JSON data, and fails the test if they did not contain equivalent
// content. If the fatal parameter is true, it will fail with t.Fatalf(),
// otherwise with t.Errorf().
func CheckMatch(t *testing.T, wanted string, got []byte, fatal bool) {
	differ := gojsondiff.New()
	d, err := differ.Compare([]byte(wanted), got)
	if err != nil {
		if fatal {
			t.Fatalf("JSON differ.Compare() returned error: %v", err)
			return
		}
		t.Errorf("JSON differ.Compare() returned error: %v", err)
		return
	}

	if d.Modified() {
		t.Logf("WANTED:      %#v\n", wanted)
		t.Logf("GOT (str):   %#v\n", string(got))
		t.Logf("GOT (bytes): %#v\n", got)
		if fatal {
			t.Fatalf("JSON not equivalent")
			return
		}
		t.Errorf("JSON not equivalent")
	}
}

// GetBody reads (and fails the test on error) the body from a
// recorded httptest call, and returns it as a byte string.
func GetBody(t *testing.T, rec *httptest.ResponseRecorder) []byte {
	got, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
		return []byte{}
	}
	return got
}

// CheckResponse is a simple wrapper around the GetBody and
// CheckMatch functions. It assumes that a mismatch should
// return a fatal error.
func CheckResponse(t *testing.T, rec *httptest.ResponseRecorder, wanted string) {
	got := GetBody(t, rec)
	CheckMatch(t, wanted, got, true)
}

// ConfirmOKResponse confirms that the handler returned an
// OK (200) response and that the header is set for JSON content.
func ConfirmOKResponse(t *testing.T, rec *httptest.ResponseRecorder) {
	// check that we got a 200 (OK)
	if 200 != rec.Code {
		t.Errorf("Expected %d, got %d", 200, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}
}

// ConfirmInvalidAuth confirms that the handler returned an
// Unauthorized (401) response and that the correct error
// message appeared in the JSON content.
func ConfirmInvalidAuth(t *testing.T, rec *httptest.ResponseRecorder, errMsg string) {
	// check that we got a 401 (Unauthorized)
	if 401 != rec.Code {
		t.Errorf("Expected %d, got %d", 401, rec.Code)
	}

	// check that we got a WWW-Authenticate header
	// (see https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401)
	header := rec.Result().Header
	wantHeader := "Bearer"
	gotHeader := header.Get("WWW-Authenticate")
	if gotHeader != wantHeader {
		t.Errorf("expected %v, got %v", wantHeader, gotHeader)
	}

	// check that content type was application/json
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that the right "error" JSON string was returned
	wantString := `{"error": "` + errMsg + `"}`
	if rec.Body.String() != wantString {
		t.Fatalf("expected %s, got %s", wantString, rec.Body.String())
	}
}

// ConfirmDisabledAuth confirms that the handler returned a
// Forbidden (403) response and that the correct error
// message appeared in the JSON content.
func ConfirmDisabledAuth(t *testing.T, rec *httptest.ResponseRecorder) {
	// check that we got a 403 (Forbidden)
	if 403 != rec.Code {
		t.Errorf("Expected %d, got %d", 403, rec.Code)
	}

	// check that content type was application/json
	header := rec.Result().Header
	if header.Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", header.Get("Content-Type"))
	}

	// check that the right "error" JSON string was returned
	wantString := `{"error": "Access denied"}`
	if rec.Body.String() != wantString {
		t.Fatalf("expected %s, got %s", wantString, rec.Body.String())
	}
}
