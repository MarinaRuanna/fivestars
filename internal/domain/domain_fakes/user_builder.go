package domain_fakes

import (
	"fivestars/internal/domain"
	"time"
)

type UserBuilder struct {
	*domain.Builder[domain.User]
}

func NewUserBuilder() *UserBuilder {
	dateTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	user := &domain.User{
		ID:           "test-user-id",
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "test-password",
		AvatarURL:    "https://example.com/avatar.png",
		Level:        1,
		CreatedAt:    dateTime,
		UpdatedAt:    dateTime,
	}
	return &UserBuilder{
		Builder: domain.NewBuilder[domain.User](*user),
	}
}

func (b *UserBuilder) WithID(id string) *UserBuilder {
	b.Builder.Value.ID = id
	return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.Builder.Value.Email = email
	return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.Builder.Value.Name = name
	return b
}

func (b *UserBuilder) WithPasswordHash(passwordHash string) *UserBuilder {
	b.Builder.Value.PasswordHash = passwordHash
	return b
}

func (b *UserBuilder) WithAvatarURL(avatarURL string) *UserBuilder {
	b.Builder.Value.AvatarURL = avatarURL
	return b
}

func (b *UserBuilder) WithLevel(level int) *UserBuilder {
	b.Builder.Value.Level = level
	return b
}

func (b *UserBuilder) WithCreatedAt(createdAt time.Time) *UserBuilder {
	b.Builder.Value.CreatedAt = createdAt
	return b
}

func (b *UserBuilder) WithUpdatedAt(updatedAt time.Time) *UserBuilder {
	b.Builder.Value.UpdatedAt = updatedAt
	return b
}
