package domain

import "time"

type Message struct {
	ID                string
	Body              string
	QueueName         string
	ReceiptHandle     string
	VisibilityTimeout time.Time
}
