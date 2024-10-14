package main

import (
	"context"
	"log"
	"strings"
	"time"

	pb "queueapi/pkg/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Using grpc.DialContext instead of grpc.NewClient
	conn, err := grpc.NewClient("server:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewQueueClient(conn)

	var delay time.Duration
	nilCounter := 0
	maxDelay := 30 * time.Second // Maximum delay of 30 seconds

	for {
		// Set a timeout for receiving messages
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		receiveResp, err := c.ReceiveMessage(ctx, &pb.ReceiveMessageRequest{})
		if err != nil {
			if err == context.DeadlineExceeded {
				log.Println("Could not receive message: request timed out.")
			} else if strings.Contains(err.Error(), "no messages available") {
				nilCounter++
				log.Println("No message received.")

				if nilCounter >= 3 {
					delay += 5 * time.Second
					if delay > maxDelay {
						delay = maxDelay
					}
					log.Printf("Incrementing delay to %v due to multiple nil responses.", delay)
					time.Sleep(delay)
				}
				continue
			} else {
				log.Fatalf("Could not receive message: %v", err)
				break
			}
		}

		nilCounter = 0
		delay = 0 // Reset the delay when a message is received
		log.Printf("Message received: %s", receiveResp.MessageBody)

		// Delete the message after processing
		deleteResp, err := c.DeleteMessage(ctx, &pb.DeleteMessageRequest{ReceiptHandle: receiveResp.ReceiptHandle})
		if err != nil {
			log.Fatalf("Could not delete message: %v", err)
			break
		}
		log.Printf("Message deleted: %v", deleteResp.Success)

		// If there's a delay, sleep before the next attempt
		if delay > 0 {
			log.Printf("Sleeping for %v before the next attempt.", delay)
			time.Sleep(delay)
		}
	}
}
