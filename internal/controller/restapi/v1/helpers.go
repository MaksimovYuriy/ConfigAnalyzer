package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/restapi/v1/response"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
)

func writeRequestError(w http.ResponseWriter, err error) {
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		writeJSON(w, http.StatusRequestEntityTooLarge, response.Error{Error: "request body exceeds 1 MiB"})
		return
	}
	writeJSON(w, http.StatusBadRequest, response.Error{Error: "expected multipart/form-data with a text content field"})
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func writeAnalysisError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidConfiguration):
		writeJSON(w, http.StatusBadRequest, response.Error{Error: "invalid configuration"})
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		writeJSON(w, http.StatusRequestTimeout, response.Error{Error: "request canceled or timed out"})
	default:
		writeJSON(w, http.StatusInternalServerError, response.Error{Error: "configuration analysis failed"})
	}
}
