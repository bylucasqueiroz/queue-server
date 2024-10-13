package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "queueapi/api" // Import the generated package

	"google.golang.org/grpc"
)

// QueueServer is the gRPC server that implements the Queue service
type QueueServer struct {
	pb.UnimplementedQueueServer
	queue Queue
}

// SendMessage gRPC method
func (s *QueueServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	messageID := s.queue.SendMessage(req.GetMessageBody())
	return &pb.SendMessageResponse{MessageId: messageID}, nil
}

// Implement the ReceiveMessage method with visibility timeout using the Queue
func (s *QueueServer) ReceiveMessage(ctx context.Context, req *pb.ReceiveMessageRequest) (*pb.ReceiveMessageResponse, error) {
	message := s.queue.ReceiveMessage(time.Second * 30) // 30-second visibility timeout
	if message == nil {
		return nil, fmt.Errorf("no messages available")
	}
	return &pb.ReceiveMessageResponse{
		MessageId:     message.ID,
		MessageBody:   message.Body,
		ReceiptHandle: message.ReceiptHandle,
	}, nil
}

// DeleteMessage gRPC method
func (s *QueueServer) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	success := s.queue.DeleteMessage(req.GetReceiptHandle())
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

	log.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
