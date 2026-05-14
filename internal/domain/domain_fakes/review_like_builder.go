package domain_fakes

import (
	"fivestars/internal/domain"
	"time"
)

type ReviewLikeBuilder struct {
	*domain.Builder[domain.ReviewLike]
}

func NewReviewLikeBuilder() *ReviewLikeBuilder {
	now := time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC)
	like := &domain.ReviewLike{
		UserID:    "11111111-1111-4111-8111-111111111111",
		ReviewID:  "44444444-4444-4444-8444-444444444444",
		CreatedAt: now,
	}

	return &ReviewLikeBuilder{Builder: domain.NewBuilder[domain.ReviewLike](*like)}
}

func (b *ReviewLikeBuilder) WithUserID(userID string) *ReviewLikeBuilder {
	b.Builder.Value.UserID = userID
	return b
}

func (b *ReviewLikeBuilder) WithReviewID(reviewID string) *ReviewLikeBuilder {
	b.Builder.Value.ReviewID = reviewID
	return b
}

func (b *ReviewLikeBuilder) WithoutTimestamps() *ReviewLikeBuilder {
	b.Builder.Value.CreatedAt = time.Time{}
	return b
}
