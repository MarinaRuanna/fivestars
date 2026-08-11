package controller

import (
	"time"

	"fivestars/internal/domain"
)

type CreateReviewRequest struct {
	CheckinID string `json:"checkin_id"`
	Rating    int    `json:"rating"`
	Text      string `json:"text"`
}

type ReviewResponse struct {
	ID              string `json:"review_id"`
	UserID          string `json:"user_id"`
	EstablishmentID string `json:"establishment_id"`
	CheckinID       string `json:"checkin_id"`
	Rating          int    `json:"rating"`
	Text            string `json:"text"`
	LikeCount       int    `json:"like_count"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func ToDomainReview(req CreateReviewRequest, userID string) *domain.Review {
	return &domain.Review{
		UserID:    userID,
		CheckinID: req.CheckinID,
		Rating:    req.Rating,
		Text:      req.Text,
	}
}

func ReviewFromDomain(r *domain.Review) ReviewResponse {
	return ReviewResponse{
		ID:              r.ID,
		UserID:          r.UserID,
		EstablishmentID: r.EstablishmentID,
		CheckinID:       r.CheckinID,
		Rating:          r.Rating,
		Text:            r.Text,
		LikeCount:       r.LikeCount,
		CreatedAt:       r.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:       r.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func ReviewsFromDomain(list []domain.Review) []ReviewResponse {
	res := make([]ReviewResponse, len(list))
	for i := range list {
		res[i] = ReviewFromDomain(&list[i])
	}
	return res
}
