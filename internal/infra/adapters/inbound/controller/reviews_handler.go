package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fivestars/internal/application/usecases"
	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/auth"

	"github.com/go-chi/chi/v5"
)

const reviewBodyMaxBytes int64 = 32 << 10 // 32KB

type ReviewsHandler struct {
	createUC usecases.CreateReviewUseCase
	getUC    usecases.GetReviewUseCase
	listUC   usecases.ListReviewsByEstablishmentUseCase
	likeUC   usecases.LikeReviewUseCase
	unlikeUC usecases.UnlikeReviewUseCase
}

func NewReviewsHandler(
	createUC usecases.CreateReviewUseCase,
	getUC usecases.GetReviewUseCase,
	listUC usecases.ListReviewsByEstablishmentUseCase,
	likeUC usecases.LikeReviewUseCase,
	unlikeUC usecases.UnlikeReviewUseCase,
) *ReviewsHandler {
	return &ReviewsHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		likeUC:   likeUC,
		unlikeUC: unlikeUC,
	}
}

func (h *ReviewsHandler) CreateReview(w http.ResponseWriter, r *http.Request) error {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		return customerror.NewUnauthorizedError("user not authenticated")
	}

	var req CreateReviewRequest
	if err := decodeStrictJSONBody(w, r, &req, reviewBodyMaxBytes); err != nil {
		return err
	}

	review := ToDomainReview(req, userID)

	res, err := h.createUC.Execute(r.Context(), *review)
	if err != nil {
		return err
	}

	resp := ReviewFromDomain(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func (h *ReviewsHandler) GetReview(w http.ResponseWriter, r *http.Request) error {
	reviewID := chi.URLParam(r, "id")

	review, err := h.getUC.Execute(r.Context(), reviewID)
	if err != nil {
		return err
	}

	resp := ReviewFromDomain(review)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func (h *ReviewsHandler) ListByEstablishment(w http.ResponseWriter, r *http.Request) error {
	establishmentID := chi.URLParam(r, "id")

	options, err := parseReviewListOptions(r)
	if err != nil {
		return err
	}

	list, err := h.listUC.Execute(r.Context(), options, establishmentID)
	if err != nil {
		return err
	}

	resp := ReviewsFromDomain(list)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"items": resp}); err != nil {
		return err
	}

	return nil
}

func parseReviewListOptions(r *http.Request) (domain.ReviewListOptions, error) {
	var opts domain.ReviewListOptions

	query := r.URL.Query()

	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return opts, customerror.NewValidationError("invalid limit")
		}
		opts.Limit = limit
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return opts, customerror.NewValidationError("invalid offset")
		}
		opts.Offset = offset
	}

	if minStr := query.Get("min_rating"); minStr != "" {
		minRating, err := strconv.Atoi(minStr)
		if err != nil {
			return opts, customerror.NewValidationError("invalid min_rating")
		}
		opts.MinRating = &minRating
	}

	if maxStr := query.Get("max_rating"); maxStr != "" {
		maxRating, err := strconv.Atoi(maxStr)
		if err != nil {
			return opts, customerror.NewValidationError("invalid max_rating")
		}
		opts.MaxRating = &maxRating
	}

	return opts, nil
}

func (h *ReviewsHandler) LikeReview(w http.ResponseWriter, r *http.Request) error {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		return customerror.NewUnauthorizedError("user not authenticated")
	}

	reviewID := chi.URLParam(r, "id")

	if err := h.likeUC.Execute(r.Context(), userID, reviewID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *ReviewsHandler) UnlikeReview(w http.ResponseWriter, r *http.Request) error {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		return customerror.NewUnauthorizedError("user not authenticated")
	}

	reviewID := chi.URLParam(r, "id")

	if err := h.unlikeUC.Execute(r.Context(), userID, reviewID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
