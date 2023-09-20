package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func parseID(r *http.Request) (uuid.UUID, error) {
	id := r.URL.Query().Get("id")
	if id == "" {
		return uuid.UUID{}, fmt.Errorf("required param 'id' not found")
	}
	return uuid.Parse(id)
}

func parseBody(r *http.Request, body any) error {
	err := json.NewDecoder(r.Body).Decode(body)
	if errors.Is(err, io.EOF) {
		return fmt.Errorf("request body is empty")
	}
	if err != nil {
		return fmt.Errorf("failed to decode request: %w", err)
	}
	return nil
}
