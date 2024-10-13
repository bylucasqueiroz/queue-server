package main

import (
	"context"
	"log"
	"time"

	pb "queueapi/pkg/api" // Import the generated package

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("server:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewQueueClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	receiveResp, err := c.ReceiveMessage(ctx, &pb.ReceiveMessageRequest{})
	if err != nil {
		log.Fatalf("Could not receive message: %v", err)
	}
	log.Printf("Message received: %s", receiveResp.MessageBody)

	deleteResp, err := c.DeleteMessage(ctx, &pb.DeleteMessageRequest{ReceiptHandle: receiveResp.ReceiptHandle})
	if err != nil {
		log.Fatalf("Could not delete message: %v", err)
	}
	log.Printf("Message deleted: %v", deleteResp.Success)
}
