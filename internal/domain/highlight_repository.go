package domain

import "context"

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/highlight_repository.go -package mock_domain . HighlightRepository
type HighlightRepository interface {
	Create(ctx context.Context, highlight *Highlight) error
	Delete(ctx context.Context, establishmentID, reviewID string) error
	Exists(ctx context.Context, establishmentID, reviewID string) (bool, error)
	CountByEstablishment(ctx context.Context, establishmentID string) (int, error)
	ListByEstablishment(ctx context.Context, establishmentID string) ([]Highlight, error)
}
