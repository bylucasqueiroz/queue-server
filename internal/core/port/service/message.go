package service

import "time"

// Message represents a message in the queue
type Message struct {
	ID                string
	Body              string
	ReceiptHandle     string
	VisibilityTimeout time.Time
}
