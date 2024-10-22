package service

import (
	"sync"
	"time"

	"queueserver/internal/core/port/service"

	"github.com/google/uuid"
)

type queueService struct {
	messages []*service.Message
	mu       sync.Mutex // for thread-safe access
}

func NewQueueService() service.QueueService {
	return &queueService{
		messages: make([]*service.Message, 0),
	}
}

// SendMessage pushes a message onto the queue
func (q *queueService) SendMessage(body string) string {
	q.mu.Lock()
	defer q.mu.Unlock()

	message := &service.Message{
		ID:                generateID(),
		Body:              body,
		ReceiptHandle:     generateReceiptHandle(),
		VisibilityTimeout: time.Now(), // Initial visibility timeout set to now
	}
	q.messages = append(q.messages, message)
	return message.ID
}

// ReceiveMessage retrieves a message from the queue with a visibility timeout
func (q *queueService) ReceiveMessage(timeout time.Duration) *service.Message {
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
func (q *queueService) DeleteMessage(receiptHandle string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i, msg := range q.messages {
		if msg.ReceiptHandle == receiptHandle {
			q.messages = append(q.messages[:i], q.messages[i+1:]...)
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
