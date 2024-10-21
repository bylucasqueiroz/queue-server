package main

import (
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pb "queueapi/api" // Import the generated package

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
