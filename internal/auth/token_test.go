// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package auth

import (
	"testing"
)

func TestShouldEncodeToken(t *testing.T) {
	jwtSecretKey := "keyForTesting"
	
	tknWanted := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJnaXRodWIiOiJzd2luc2xvdyJ9.OKZnGOVSvSk3T7sEXL4SCHNQewdlRIHz6NPl_2gTRIE"
	gh := "swinslow"
	tknGot, err := EncodeToken(jwtSecretKey, gh)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if tknWanted != tknGot {
		t.Errorf("expected %s, got %s", tknWanted, tknGot)
	}
	
	tknWanted = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJnaXRodWIiOiJvY3RvY2F0In0.-TVt3uUrT6wZYMabZaH3wPbkUyBzSFdeyiT7NixsYpY"
	gh = "octocat"
	tknGot, err = EncodeToken(jwtSecretKey, gh)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if tknWanted != tknGot {
		t.Errorf("expected %s, got %s", tknWanted, tknGot)
	}
}

func TestShouldDecodeToken(t *testing.T) {
	jwtSecretKey := "keyForTesting"

	tknRecv := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJnaXRodWIiOiJzd2luc2xvdyJ9.OKZnGOVSvSk3T7sEXL4SCHNQewdlRIHz6NPl_2gTRIE"
	ghWanted := "swinslow"
	ghGot, err := DecodeToken(jwtSecretKey, tknRecv)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if ghWanted != ghGot {
		t.Errorf("expected %s, got %s", ghWanted, ghGot)
	}

	tknRecv = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJnaXRodWIiOiJvY3RvY2F0In0.-TVt3uUrT6wZYMabZaH3wPbkUyBzSFdeyiT7NixsYpY"
	ghWanted = "octocat"
	ghGot, err = DecodeToken(jwtSecretKey, tknRecv)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if ghWanted != ghGot {
		t.Errorf("expected %s, got %s", ghWanted, ghGot)
	}
}

func TestShouldFailToDecodeInvalidToken(t *testing.T) {
	jwtSecretKey := "keyForTesting"

	// fail if not a token at all
	tknRecv := "oops"
	_, err := DecodeToken(jwtSecretKey, tknRecv)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}

	// fail if token doesn't have a "github" field
	// this one is just {"test": "swinslow"}
	tknRecv = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXN0Ijoic3dpbnNsb3cifQ.Ri_stIs5zPl_zNySTx2nLHLXa_rGT4s9xQ26zfRU3HU"
	_, err = DecodeToken(jwtSecretKey, tknRecv)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}
