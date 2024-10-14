package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pb "queueapi/pkg/api" // Import the generated package

	"google.golang.org/grpc"
)

var (
	messagesProduced = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "messages_produced_total",
		Help: "Total number of messages produced to the queue.",
	})

	messagesConsumed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "messages_consumed_total",
		Help: "Total number of messages consumed from the queue.",
	})

	messagesDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "messages_deleted_total",
		Help: "Total number of messages consumed from the queue.",
	})
)

func init() {
	prometheus.MustRegister(messagesProduced)
	prometheus.MustRegister(messagesConsumed)
	prometheus.MustRegister(messagesDeleted)
}

// QueueServer is the gRPC server that implements the Queue service
type QueueServer struct {
	pb.UnimplementedQueueServer
	queue Queue
}

// SendMessage gRPC method
func (s *QueueServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	messageID := s.queue.SendMessage(req.GetMessageBody())

	if messageID != "" {
		messagesProduced.Inc()
	}

	return &pb.SendMessageResponse{MessageId: messageID}, nil
}

// Implement the ReceiveMessage method with visibility timeout using the Queue
func (s *QueueServer) ReceiveMessage(ctx context.Context, req *pb.ReceiveMessageRequest) (*pb.ReceiveMessageResponse, error) {
	message := s.queue.ReceiveMessage(time.Second * 30) // 30-second visibility timeout
	if message == nil {
		return nil, fmt.Errorf("no messages available")
	}

	messagesConsumed.Inc()

	return &pb.ReceiveMessageResponse{
		MessageId:     message.ID,
		MessageBody:   message.Body,
		ReceiptHandle: message.ReceiptHandle,
	}, nil
}

// DeleteMessage gRPC method
func (s *QueueServer) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	success := s.queue.DeleteMessage(req.GetReceiptHandle())

	if success {
		messagesDeleted.Inc()
	}

	return &pb.DeleteMessageResponse{Success: success}, nil
}

func main() {
	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Initialize the queue and register the gRPC server
	queueServer := &QueueServer{queue: Queue{messages: []*Message{}}}

	pb.RegisterQueueServer(s, queueServer)

	// Start Prometheus metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	log.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Message represents a message in the queue
type Message struct {
	ID                string
	Body              string
	ReceiptHandle     string
	VisibilityTimeout time.Time
}

// Queue represents an in-memory queue
type Queue struct {
	messages []*Message
	mu       sync.Mutex // for thread-safe access
}

// SendMessage pushes a message onto the queue
func (q *Queue) SendMessage(body string) string {
	q.mu.Lock()
	defer q.mu.Unlock()

	message := &Message{
		ID:                generateID(),
		Body:              body,
		ReceiptHandle:     generateReceiptHandle(),
		VisibilityTimeout: time.Now(), // Initial visibility timeout set to now
	}
	q.messages = append(q.messages, message)
	return message.ID
}

// ReceiveMessage retrieves a message from the queue with a visibility timeout
func (q *Queue) ReceiveMessage(timeout time.Duration) *Message {
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
func (q *Queue) DeleteMessage(receiptHandle string) bool {
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
