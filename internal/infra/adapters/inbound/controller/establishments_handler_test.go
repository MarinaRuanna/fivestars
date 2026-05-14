package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/auth"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubGetEstablishmentDetailUseCase struct {
	result *domain.EstablishmentDetail
	err    error
}

func (s stubGetEstablishmentDetailUseCase) Execute(ctx context.Context, establishmentID string) (*domain.EstablishmentDetail, error) {
	return s.result, s.err
}

type stubCreateEstablishmentUseCase struct {
	result *domain.Establishment
	err    error
}

func (s stubCreateEstablishmentUseCase) Execute(ctx context.Context, userID string, input domain.Establishment) (*domain.Establishment, error) {
	return s.result, s.err
}

type stubClaimEstablishmentOwnershipUseCase struct {
	result *domain.Establishment
	err    error
}

func (s stubClaimEstablishmentOwnershipUseCase) Execute(ctx context.Context, userID, establishmentID, claimQRCode string) (*domain.Establishment, error) {
	return s.result, s.err
}

type stubListEstablishmentsUseCase struct {
	result []domain.Establishment
	err    error
}

func (s stubListEstablishmentsUseCase) Execute(ctx context.Context) ([]domain.Establishment, error) {
	return s.result, s.err
}

type stubGetEstablishmentStatsUseCase struct {
	result *domain.EstablishmentStats
	err    error
}

func (s stubGetEstablishmentStatsUseCase) Execute(ctx context.Context, establishmentID string) (*domain.EstablishmentStats, error) {
	return s.result, s.err
}

type stubCreateHighlightUseCase struct {
	result *domain.Highlight
	err    error
}

func (s stubCreateHighlightUseCase) Execute(ctx context.Context, userID, establishmentID, reviewID string) (*domain.Highlight, error) {
	return s.result, s.err
}

type stubDeleteHighlightUseCase struct {
	err error
}

func (s stubDeleteHighlightUseCase) Execute(ctx context.Context, userID, establishmentID, reviewID string) error {
	return s.err
}

func TestEstablishmentsHandler_GetStats(t *testing.T) {
	handler := NewEstablishmentsHandler(
		stubCreateEstablishmentUseCase{},
		stubClaimEstablishmentOwnershipUseCase{},
		stubGetEstablishmentDetailUseCase{},
		stubListEstablishmentsUseCase{},
		stubGetEstablishmentStatsUseCase{
			result: &domain.EstablishmentStats{
				AverageRating:          4.5,
				TotalReviews:           10,
				TotalLikes:             15,
				HighlightedReviewCount: 2,
			},
		},
		stubCreateHighlightUseCase{},
		stubDeleteHighlightUseCase{},
	)

	req := httptest.NewRequest(http.MethodGet, "/establishments/222/stats", nil)
	req = withRouteParam(req, "id", "22222222-2222-4222-8222-222222222222")
	rec := httptest.NewRecorder()

	err := handler.GetStats(rec, req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp EstablishmentStatsResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, 4.5, resp.AverageRating)
	assert.Equal(t, 10, resp.TotalReviews)
}

func TestEstablishmentsHandler_CreateHighlight(t *testing.T) {
	now := time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC)
	handler := NewEstablishmentsHandler(
		stubCreateEstablishmentUseCase{},
		stubClaimEstablishmentOwnershipUseCase{},
		stubGetEstablishmentDetailUseCase{},
		stubListEstablishmentsUseCase{},
		stubGetEstablishmentStatsUseCase{},
		stubCreateHighlightUseCase{
			result: &domain.Highlight{
				EstablishmentID: "22222222-2222-4222-8222-222222222222",
				ReviewID:        "44444444-4444-4444-8444-444444444444",
				CreatedByUserID: "11111111-1111-4111-8111-111111111111",
				CreatedAt:       now,
			},
		},
		stubDeleteHighlightUseCase{},
	)

	body := bytes.NewBufferString(`{"review_id":"44444444-4444-4444-8444-444444444444"}`)
	req := httptest.NewRequest(http.MethodPost, "/establishments/222/highlights", body)
	req = req.WithContext(auth.WithUserID(req.Context(), "11111111-1111-4111-8111-111111111111"))
	req = withRouteParam(req, "id", "22222222-2222-4222-8222-222222222222")
	rec := httptest.NewRecorder()

	err := handler.CreateHighlight(rec, req)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)

	var resp HighlightResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "44444444-4444-4444-8444-444444444444", resp.ReviewID)
}

func TestEstablishmentsHandler_DeleteHighlight(t *testing.T) {
	handler := NewEstablishmentsHandler(
		stubCreateEstablishmentUseCase{},
		stubClaimEstablishmentOwnershipUseCase{},
		stubGetEstablishmentDetailUseCase{},
		stubListEstablishmentsUseCase{},
		stubGetEstablishmentStatsUseCase{},
		stubCreateHighlightUseCase{},
		stubDeleteHighlightUseCase{},
	)

	req := httptest.NewRequest(http.MethodDelete, "/establishments/222/highlights/444", nil)
	req = req.WithContext(auth.WithUserID(req.Context(), "11111111-1111-4111-8111-111111111111"))
	req = withRouteParams(req, map[string]string{
		"id":       "22222222-2222-4222-8222-222222222222",
		"reviewId": "44444444-4444-4444-8444-444444444444",
	})
	rec := httptest.NewRecorder()

	err := handler.DeleteHighlight(rec, req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestEstablishmentsHandler_GetEstablishment(t *testing.T) {
	establishment := domain.Establishment{
		ID:        "22222222-2222-4222-8222-222222222222",
		Name:      "Cafe Central",
		Slug:      "cafe-central",
		Category:  "cafe",
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	review := domain.Review{
		ID:              "44444444-4444-4444-8444-444444444444",
		UserID:          "11111111-1111-4111-8111-111111111111",
		EstablishmentID: establishment.ID,
		CheckinID:       "33333333-3333-4333-8333-333333333333",
		Rating:          5,
		Text:            "Atendimento excelente e ambiente ótimo.",
		CreatedAt:       time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	handler := NewEstablishmentsHandler(
		stubCreateEstablishmentUseCase{},
		stubClaimEstablishmentOwnershipUseCase{},
		stubGetEstablishmentDetailUseCase{
			result: &domain.EstablishmentDetail{
				Establishment: establishment,
				Highlights:    []domain.Review{review},
			},
		},
		stubListEstablishmentsUseCase{},
		stubGetEstablishmentStatsUseCase{},
		stubCreateHighlightUseCase{},
		stubDeleteHighlightUseCase{},
	)

	req := httptest.NewRequest(http.MethodGet, "/establishments/222", nil)
	req = withRouteParam(req, "id", establishment.ID)
	rec := httptest.NewRecorder()

	err := handler.GetEstablishment(rec, req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp EstablishmentDetailResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, establishment.ID, resp.ID)
	require.Len(t, resp.Highlights, 1)
	assert.Equal(t, review.ID, resp.Highlights[0].ID)
}

func TestEstablishmentsHandler_CreateHighlight_Unauthorized(t *testing.T) {
	handler := NewEstablishmentsHandler(
		stubCreateEstablishmentUseCase{},
		stubClaimEstablishmentOwnershipUseCase{},
		stubGetEstablishmentDetailUseCase{},
		stubListEstablishmentsUseCase{},
		stubGetEstablishmentStatsUseCase{},
		stubCreateHighlightUseCase{},
		stubDeleteHighlightUseCase{},
	)

	req := httptest.NewRequest(http.MethodPost, "/establishments/222/highlights", bytes.NewBufferString(`{"review_id":"44444444-4444-4444-8444-444444444444"}`))
	req = withRouteParam(req, "id", "22222222-2222-4222-8222-222222222222")
	rec := httptest.NewRecorder()

	err := handler.CreateHighlight(rec, req)

	require.Error(t, err)
	errorType, ok := customerror.TypeOf(err)
	require.True(t, ok)
	assert.Equal(t, customerror.UnauthorizedErrorType, errorType)
}

func TestEstablishmentsHandler_CreateEstablishment(t *testing.T) {
	establishment := domain.Establishment{
		ID:        "22222222-2222-4222-8222-222222222222",
		OwnerID:   "11111111-1111-4111-8111-111111111111",
		Name:      "Cafe Central",
		Slug:      "cafe-central",
		Category:  "cafe",
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	handler := NewEstablishmentsHandler(
		stubCreateEstablishmentUseCase{result: &establishment},
		stubClaimEstablishmentOwnershipUseCase{},
		stubGetEstablishmentDetailUseCase{},
		stubListEstablishmentsUseCase{},
		stubGetEstablishmentStatsUseCase{},
		stubCreateHighlightUseCase{},
		stubDeleteHighlightUseCase{},
	)

	body := bytes.NewBufferString(`{"name":"Cafe Central","slug":"cafe-central","category":"cafe"}`)
	req := httptest.NewRequest(http.MethodPost, "/establishments", body)
	req = req.WithContext(auth.WithUserID(req.Context(), "11111111-1111-4111-8111-111111111111"))
	rec := httptest.NewRecorder()

	err := handler.CreateEstablishment(rec, req)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)

	var resp EstablishmentResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, establishment.ID, resp.ID)
	assert.Equal(t, establishment.Name, resp.Name)
}

func TestEstablishmentsHandler_ClaimOwnership(t *testing.T) {
	establishment := domain.Establishment{
		ID:      "22222222-2222-4222-8222-222222222222",
		OwnerID: "11111111-1111-4111-8111-111111111111",
	}
	handler := NewEstablishmentsHandler(
		stubCreateEstablishmentUseCase{},
		stubClaimEstablishmentOwnershipUseCase{result: &establishment},
		stubGetEstablishmentDetailUseCase{},
		stubListEstablishmentsUseCase{},
		stubGetEstablishmentStatsUseCase{},
		stubCreateHighlightUseCase{},
		stubDeleteHighlightUseCase{},
	)

	req := httptest.NewRequest(http.MethodPost, "/establishments/222/claim", bytes.NewBufferString(`{"qr_code":"qr-cafe-central"}`))
	req = req.WithContext(auth.WithUserID(req.Context(), "11111111-1111-4111-8111-111111111111"))
	req = withRouteParam(req, "id", establishment.ID)
	rec := httptest.NewRecorder()

	err := handler.ClaimOwnership(rec, req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp ClaimEstablishmentOwnershipResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, establishment.ID, resp.EstablishmentID)
	assert.Equal(t, establishment.OwnerID, resp.OwnerID)
}

func withRouteParam(req *http.Request, key, value string) *http.Request {
	return withRouteParams(req, map[string]string{key: value})
}

func withRouteParams(req *http.Request, params map[string]string) *http.Request {
	routeCtx := chi.NewRouteContext()
	for k, v := range params {
		routeCtx.URLParams.Add(k, v)
	}

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}
