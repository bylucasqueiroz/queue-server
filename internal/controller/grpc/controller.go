package grpc

import (
	"context"
	"fmt"
	"time"

	"queueserver/internal/core/port/service"

	proto "queueserver/api"
)

type queueController struct {
	proto.UnimplementedQueueServer
	queueService service.QueueService
}

func NewQueueController(queueService service.QueueService) proto.QueueServer {
	return &queueController{
		queueService: queueService,
	}
}

// SendMessage gRPC method
func (s *queueController) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	messageID := s.queueService.SendMessage(ctx, req.QueueName, req.GetMessageBody())

	return &proto.SendMessageResponse{MessageId: messageID}, nil
}

// Implement the ReceiveMessage method with visibility timeout using the Queue
func (s *queueController) ReceiveMessage(ctx context.Context, req *proto.ReceiveMessageRequest) (*proto.ReceiveMessageResponse, error) {
	message := s.queueService.ReceiveMessage(ctx, req.QueueName, time.Second*30) // 30-second visibility timeout
	if message == nil {
		return nil, fmt.Errorf("no messages available")
	}

	return &proto.ReceiveMessageResponse{
		MessageId:     message.ID,
		MessageBody:   message.Body,
		ReceiptHandle: message.ReceiptHandle,
	}, nil
}

// DeleteMessage gRPC method
func (s *queueController) DeleteMessage(ctx context.Context, req *proto.DeleteMessageRequest) (*proto.DeleteMessageResponse, error) {
	success := s.queueService.DeleteMessage(ctx, req.QueueName, req.GetReceiptHandle())

	return &proto.DeleteMessageResponse{Success: success}, nil
}
