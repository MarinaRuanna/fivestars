package controller

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_decodeStrictJSONBody_RejectsUnknownField(t *testing.T) {
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{"email":"a@a.com","password":"123456","x":1}`))
	rec := httptest.NewRecorder()

	var dst LoginRequest
	err := decodeStrictJSONBody(rec, req, &dst, 1024)
	require.Error(t, err)
}

func Test_decodeStrictJSONBody_RejectsOversizedPayload(t *testing.T) {
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{"email":"a@a.com","password":"123456"}`))
	rec := httptest.NewRecorder()

	var dst LoginRequest
	err := decodeStrictJSONBody(rec, req, &dst, 8)
	require.Error(t, err)
}

func Test_decodeStrictJSONBody_RejectsMultipleJSONObjects(t *testing.T) {
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{"email":"a@a.com","password":"123456"}{"email":"b@b.com","password":"123456"}`))
	rec := httptest.NewRecorder()

	var dst LoginRequest
	err := decodeStrictJSONBody(rec, req, &dst, 4096)
	require.Error(t, err)
}
