// Package response provides
package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	errEncode = errors.New("failed to encode JSON")
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("%v: %w", errEncode, err)
	}
	return nil
}
