package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"fivestars/internal/domain/customerror"
)

func decodeStrictJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}, maxBodyBytes int64) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return customerror.NewValidationError("invalid body")
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return customerror.NewValidationError("body must contain a single JSON object")
	}

	return nil
}
