package grpc

import (
	"context"
	"fmt"
	"time"

	"queueapi/internal/core/port/service"

	proto "queueapi/api"
)

// https://github.com/phamtai97/go-experienced-series/blob/main/hexagonal/internal/core/service/user.go

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
	messageID := s.queueService.SendMessage(req.GetMessageBody())

	// if messageID != "" {
	// 	messagesProduced.Inc()
	// }

	return &proto.SendMessageResponse{MessageId: messageID}, nil
}

// Implement the ReceiveMessage method with visibility timeout using the Queue
func (s *queueController) ReceiveMessage(ctx context.Context, req *proto.ReceiveMessageRequest) (*proto.ReceiveMessageResponse, error) {
	message := s.queueService.ReceiveMessage(time.Second * 30) // 30-second visibility timeout
	if message == nil {
		return nil, fmt.Errorf("no messages available")
	}

	// messagesConsumed.Inc()

	return &proto.ReceiveMessageResponse{
		MessageId:     message.ID,
		MessageBody:   message.Body,
		ReceiptHandle: message.ReceiptHandle,
	}, nil
}

// DeleteMessage gRPC method
func (s *queueController) DeleteMessage(ctx context.Context, req *proto.DeleteMessageRequest) (*proto.DeleteMessageResponse, error) {
	success := s.queueService.DeleteMessage(req.GetReceiptHandle())

	if success {
		// messagesDeleted.Inc()
	}

	return &proto.DeleteMessageResponse{Success: success}, nil
}
