package controller

import (
	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
	"time"
)

type EstablishmentResponse struct {
	ID        string   `json:"establishment_id"`
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	Category  string   `json:"category"`
	Address   string   `json:"address,omitempty"`
	Lat       *float64 `json:"lat,omitempty"`
	Lng       *float64 `json:"lng,omitempty"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type EstablishmentDetailResponse struct {
	ID         string           `json:"establishment_id"`
	Name       string           `json:"name"`
	Slug       string           `json:"slug"`
	Category   string           `json:"category"`
	Address    string           `json:"address,omitempty"`
	Lat        *float64         `json:"lat,omitempty"`
	Lng        *float64         `json:"lng,omitempty"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
	Highlights []ReviewResponse `json:"highlights"`
}

type EstablishmentListResponse struct {
	Items []EstablishmentResponse `json:"items"`
}

type CreateHighlightRequest struct {
	ReviewID string `json:"review_id" validate:"required,uuid4"`
}

type CreateEstablishmentRequest struct {
	Name     string   `json:"name" validate:"required"`
	Slug     string   `json:"slug" validate:"required"`
	Category string   `json:"category" validate:"required"`
	Address  string   `json:"address,omitempty"`
	Lat      *float64 `json:"lat,omitempty"`
	Lng      *float64 `json:"lng,omitempty"`
	QRCode   string   `json:"qr_code,omitempty"`
}

type ClaimEstablishmentOwnershipRequest struct {
	QRCode string `json:"qr_code" validate:"required"`
}

type ClaimEstablishmentOwnershipResponse struct {
	EstablishmentID string `json:"establishment_id"`
	OwnerID         string `json:"owner_id"`
}

type HighlightResponse struct {
	EstablishmentID string `json:"establishment_id"`
	ReviewID        string `json:"review_id"`
	CreatedByUserID string `json:"created_by_user_id"`
	CreatedAt       string `json:"created_at"`
}

type EstablishmentStatsResponse struct {
	AverageRating          float64 `json:"average_rating"`
	TotalReviews           int     `json:"total_reviews"`
	TotalLikes             int     `json:"total_likes"`
	HighlightedReviewCount int     `json:"highlighted_review_count"`
}

func (r *CreateHighlightRequest) Validate() error {
	if err := validator.Validate(r); err != nil {
		return customerror.NewValidationError(err.Error())
	}

	return nil
}

func (r *CreateEstablishmentRequest) Validate() error {
	if err := validator.Validate(r); err != nil {
		return customerror.NewValidationError(err.Error())
	}

	return nil
}

func (r *ClaimEstablishmentOwnershipRequest) Validate() error {
	if err := validator.Validate(r); err != nil {
		return customerror.NewValidationError(err.Error())
	}

	return nil
}

func FromDomain(e *domain.Establishment) EstablishmentResponse {
	return EstablishmentResponse{
		ID:        e.ID,
		Name:      e.Name,
		Slug:      e.Slug,
		Category:  e.Category,
		Address:   e.Address,
		Lat:       e.Lat,
		Lng:       e.Lng,
		CreatedAt: e.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: e.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func FromDomainList(estabs []domain.Establishment) []EstablishmentResponse {
	responses := make([]EstablishmentResponse, len(estabs))
	for i, estab := range estabs {
		responses[i] = FromDomain(&estab)
	}
	return responses
}

func HighlightFromDomain(h *domain.Highlight) HighlightResponse {
	return HighlightResponse{
		EstablishmentID: h.EstablishmentID,
		ReviewID:        h.ReviewID,
		CreatedByUserID: h.CreatedByUserID,
		CreatedAt:       h.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func EstablishmentStatsFromDomain(stats *domain.EstablishmentStats) EstablishmentStatsResponse {
	return EstablishmentStatsResponse{
		AverageRating:          stats.AverageRating,
		TotalReviews:           stats.TotalReviews,
		TotalLikes:             stats.TotalLikes,
		HighlightedReviewCount: stats.HighlightedReviewCount,
	}
}

func EstablishmentDetailFromDomain(detail *domain.EstablishmentDetail) EstablishmentDetailResponse {
	return EstablishmentDetailResponse{
		ID:         detail.Establishment.ID,
		Name:       detail.Establishment.Name,
		Slug:       detail.Establishment.Slug,
		Category:   detail.Establishment.Category,
		Address:    detail.Establishment.Address,
		Lat:        detail.Establishment.Lat,
		Lng:        detail.Establishment.Lng,
		CreatedAt:  detail.Establishment.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  detail.Establishment.UpdatedAt.UTC().Format(time.RFC3339),
		Highlights: ReviewsFromDomain(detail.Highlights),
	}
}

func ToDomainEstablishment(req CreateEstablishmentRequest, userID string) domain.Establishment {
	return domain.Establishment{
		OwnerID:  userID,
		Name:     req.Name,
		Slug:     req.Slug,
		Category: req.Category,
		Address:  req.Address,
		Lat:      req.Lat,
		Lng:      req.Lng,
		QRCode:   req.QRCode,
	}
}

func ClaimOwnershipFromDomain(e *domain.Establishment) ClaimEstablishmentOwnershipResponse {
	return ClaimEstablishmentOwnershipResponse{
		EstablishmentID: e.ID,
		OwnerID:         e.OwnerID,
	}
}
