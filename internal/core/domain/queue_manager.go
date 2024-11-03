package domain

type QueueManager struct {
	QueueName string
	Messages  []*Message
}
