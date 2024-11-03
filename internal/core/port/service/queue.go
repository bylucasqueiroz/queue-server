package service

import (
	"context"
	"time"

	"queueserver/internal/core/domain"
)

type QueueService interface {
	SendMessage(ctx context.Context, queueName string, body string) string
	ReceiveMessage(ctx context.Context, queueName string, timeout time.Duration) *domain.Message
	DeleteMessage(ctx context.Context, queueName string, receiptHandle string) bool
}
