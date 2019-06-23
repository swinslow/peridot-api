// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func extractIDasU32(r *http.Request) (uint32, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return 0, fmt.Errorf("Missing ID in endpoint")
	}
	p, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Invalid ID in endpoint")
	}

	return uint32(p), nil
}
