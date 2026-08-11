package reviewlikes

import (
	"time"

	"fivestars/internal/domain"
)

type ReviewLikeDTO struct {
	UserID    string    `json:"user_id"`
	ReviewID  string    `json:"review_id"`
	CreatedAt time.Time `json:"created_at"`
}

func FromDomain(like *domain.ReviewLike) (*ReviewLikeDTO, error) {
	if err := like.Validate(); err != nil {
		return nil, err
	}

	dto := &ReviewLikeDTO{
		UserID:    like.UserID,
		ReviewID:  like.ReviewID,
		CreatedAt: like.CreatedAt,
	}

	if dto.CreatedAt.IsZero() {
		dto.CreatedAt = time.Now().UTC()
	}

	return dto, nil
}
