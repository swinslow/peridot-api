// Package handlerutils contains utility functions for testing the
// obsidian API handlers.
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
