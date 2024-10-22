package service

import "time"

type QueueService interface {
	SendMessage(body string) string
	ReceiveMessage(timeout time.Duration) *Message
	DeleteMessage(receiptHandle string) bool
}
