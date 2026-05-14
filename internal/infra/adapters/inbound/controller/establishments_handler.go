package controller

import (
	"encoding/json"
	"fivestars/internal/application/usecases"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const highlightBodyMaxBytes int64 = 8 << 10 // 8KB

type EstablishmentsHandler struct {
	getDetailUC       usecases.GetEstablishmentDetailUseCase
	listUC            usecases.ListEstablishmentsUseCase
	statsUC           usecases.GetEstablishmentStatsUseCase
	createHighlightUC usecases.CreateHighlightUseCase
	deleteHighlightUC usecases.DeleteHighlightUseCase
}

func NewEstablishmentsHandler(
	getDetailUC usecases.GetEstablishmentDetailUseCase,
	listUC usecases.ListEstablishmentsUseCase,
	statsUC usecases.GetEstablishmentStatsUseCase,
	createHighlightUC usecases.CreateHighlightUseCase,
	deleteHighlightUC usecases.DeleteHighlightUseCase,
) *EstablishmentsHandler {
	return &EstablishmentsHandler{
		getDetailUC:       getDetailUC,
		listUC:            listUC,
		statsUC:           statsUC,
		createHighlightUC: createHighlightUC,
		deleteHighlightUC: deleteHighlightUC,
	}
}

func (c *EstablishmentsHandler) GetEstablishment(w http.ResponseWriter, r *http.Request) error {
	establishmentID := chi.URLParam(r, "id")

	detail, err := c.getDetailUC.Execute(r.Context(), establishmentID)
	if err != nil {
		return err
	}

	resp := EstablishmentDetailFromDomain(detail)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func (c *EstablishmentsHandler) ListEstablishments(w http.ResponseWriter, r *http.Request) error {
	items, err := c.listUC.Execute(r.Context())
	if err != nil {
		return err
	}

	listDTO := FromDomainList(items)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"items": listDTO}); err != nil {
		return err
	}

	return nil
}

func (c *EstablishmentsHandler) GetStats(w http.ResponseWriter, r *http.Request) error {
	establishmentID := chi.URLParam(r, "id")

	stats, err := c.statsUC.Execute(r.Context(), establishmentID)
	if err != nil {
		return err
	}

	resp := EstablishmentStatsFromDomain(stats)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func (c *EstablishmentsHandler) CreateHighlight(w http.ResponseWriter, r *http.Request) error {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		return customerror.NewUnauthorizedError("user not authenticated")
	}

	establishmentID := chi.URLParam(r, "id")

	var req CreateHighlightRequest
	if err := decodeStrictJSONBody(w, r, &req, highlightBodyMaxBytes); err != nil {
		return err
	}
	if err := req.Validate(); err != nil {
		return err
	}

	highlight, err := c.createHighlightUC.Execute(r.Context(), userID, establishmentID, req.ReviewID)
	if err != nil {
		return err
	}

	resp := HighlightFromDomain(highlight)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func (c *EstablishmentsHandler) DeleteHighlight(w http.ResponseWriter, r *http.Request) error {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		return customerror.NewUnauthorizedError("user not authenticated")
	}

	establishmentID := chi.URLParam(r, "id")
	reviewID := chi.URLParam(r, "reviewId")

	if err := c.deleteHighlightUC.Execute(r.Context(), userID, establishmentID, reviewID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
