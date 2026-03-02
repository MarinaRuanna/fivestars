package controller

import (
	"net/http"

	"fivestars/internal/infra/adapters/inbound/httperror"
)

func EncodeError(w http.ResponseWriter, err error) {
	httperror.Encode(w, err)
}
