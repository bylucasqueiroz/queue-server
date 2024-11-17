package service

import (
	"context"
	"sync"
	"time"

	"queueserver/internal/adapter/repository"
	"queueserver/internal/core/domain"
	"queueserver/internal/core/port/service"

	"github.com/google/uuid"
)

type queueService struct {
	messages     []*domain.Message
	queueRepo    *repository.PostgresQueueRepository
	messageRepos *repository.PostgresMessageRepository
	mu           sync.Mutex // for thread-safe access
}

func NewQueueService(queueRepo *repository.PostgresQueueRepository, messageRepo *repository.PostgresMessageRepository) service.QueueService {
	return &queueService{
		messages:     make([]*domain.Message, 0),
		queueRepo:    queueRepo,
		messageRepos: messageRepo,
	}
}

// SendMessage pushes a message onto the queue
func (q *queueService) SendMessage(ctx context.Context, queueName string, body string) string {
	q.mu.Lock()
	defer q.mu.Unlock()

	message := &domain.Message{
		ID:                generateID(),
		Body:              body,
		ReceiptHandle:     generateReceiptHandle(),
		QueueName:         queueName,
		VisibilityTimeout: time.Now(), // Initial visibility timeout set to now
	}

	q.messages = append(q.messages, message)

	err := q.messageRepos.Save(ctx, message)
	if err != nil {
		return ""
	}

	return message.ID
}

// ReceiveMessage retrieves a message from the queue with a visibility timeout
func (q *queueService) ReceiveMessage(ctx context.Context, queueName string, timeout time.Duration) *domain.Message {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, msg := range q.messages {
		if time.Now().After(msg.VisibilityTimeout) {
			msg.VisibilityTimeout = time.Now().Add(timeout)
			return msg
		}
	}

	return nil // No available message
}

// DeleteMessage deletes a message using its receipt handle
func (q *queueService) DeleteMessage(ctx context.Context, queueName string, receiptHandle string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	for j, msg := range q.messages {
		if msg.ReceiptHandle == receiptHandle {
			q.messages = append(q.messages[:j], q.messages[j+1:]...)
			return true
		}
	}

	return false
}

// Utility functions to generate IDs and receipt handles
func generateID() string {
	return uuid.New().String()
}

func generateReceiptHandle() string {
	return uuid.New().String()
}
