package repository

import (
	"context"

	"queueserver/internal/core/domain"
)

type QueueRepository interface {
	Save(ctx context.Context, message *domain.Queue) error
	GetByName(ctx context.Context, name string) (*domain.Queue, error)
	Delete(ctx context.Context, name string) error
}
