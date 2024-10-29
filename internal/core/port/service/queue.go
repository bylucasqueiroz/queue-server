package service

import (
	"context"
	"time"

	"queueserver/internal/core/domain"
)

type QueueService interface {
	SendMessage(ctx context.Context, body string) string
	ReceiveMessage(ctx context.Context, timeout time.Duration) *domain.Message
	DeleteMessage(ctx context.Context, receiptHandle string) bool
}
