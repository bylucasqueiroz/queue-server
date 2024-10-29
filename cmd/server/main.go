package main

import (
	"log"
	"net"
	"net/http"

	"queueserver/internal/adapter/config"
	"queueserver/internal/adapter/repository"
	"queueserver/internal/core/service"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pb "queueserver/api" // Import the generated package
	grpcCtrl "queueserver/internal/controller/grpc"

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

	// Create a new Config
	config := config.NewConfig()

	// Create a Message Repository
	messageRepo, err := repository.NewPostgresMessageRepository(config)
	if err != nil {
		panic("error to create a Message Repository")
	}

	// Create a Queue Repository
	queueRepo, err := repository.NewPostgresQueueRepository(config)
	if err != nil {
		panic("error to create a Queue Repository")
	}

	// Create a new Service
	queueService := service.NewQueueService(queueRepo, messageRepo)

	// Create a new Controller
	userController := grpcCtrl.NewQueueController(queueService)

	// Create a new gRPC server
	s := grpc.NewServer()
	pb.RegisterQueueServer(s, userController)

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
