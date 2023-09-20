package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
)

func parseParam(r *http.Request, name string) (string, error) {
	param := r.URL.Query().Get(name)
	if param == "" {
		return "", fmt.Errorf("required param '%v' not found", name)
	}
	return param, nil
}

func parseParamUint64(r *http.Request, name string) (uint64, error) {
	paramS, err := parseParam(r, name)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(paramS, 0, 64)
}

func parseParamUuid(r *http.Request, name string) (uuid.UUID, error) {
	paramS, err := parseParam(r, name)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.Parse(paramS)
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

func checkError(w http.ResponseWriter, err error) bool {
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}

	return false
}
