package repository

import (
	"context"

	"queueserver/internal/core/domain"
)

type MessageRepository interface {
	Save(ctx context.Context, message *domain.Message) error
	GetByMessageID(ctx context.Context, messageId string) (*domain.Message, error)
	Delete(ctx context.Context, messageId string) error
}
