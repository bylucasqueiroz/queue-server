package main

import (
	"fmt"
	"log"

	"queueserver/internal/adapter/config"
	"queueserver/internal/adapter/repository"
	grpcCtrl "queueserver/internal/controller/grpc"
	grpcConfig "queueserver/internal/core/config"
	"queueserver/internal/core/server"
	"queueserver/internal/core/server/grpc"
	"queueserver/internal/core/service"

	pb "queueserver/api"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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
	// Carrega o arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	// Create a new Config
	config := config.NewConfig()

	// Create a Message Repository
	messageRepo, err := repository.NewPostgresMessageRepository(config)
	if err != nil {
		panic(fmt.Sprintf("error to create a Message Repository: %v", err))
	}

	// Create a Queue Repository
	queueRepo, err := repository.NewPostgresQueueRepository(config)
	if err != nil {
		panic(fmt.Sprintf("error to create a Queue Repository: %v", err))
	}

	// Create a new Service
	queueService := service.NewQueueService(queueRepo, messageRepo)

	// Create a new Controller
	userController := grpcCtrl.NewQueueController(queueService)

	// Create the gRPC server
	grpcServer, err := grpc.NewGrpcServer(
		grpcConfig.GrpcServerConfig{
			Port: 50051,
			KeepaliveParams: keepalive.ServerParameters{
				MaxConnectionIdle:     100,
				MaxConnectionAge:      7200,
				MaxConnectionAgeGrace: 60,
				Time:                  10,
				Timeout:               3,
			},
			KeepalivePolicy: keepalive.EnforcementPolicy{
				MinTime:             10,
				PermitWithoutStream: true,
			},
		},
	)
	if err != nil {
		log.Fatalf("failed to new grpc server err=%s\n", err.Error())
	}

	// Start the gRPC server
	go grpcServer.Start(
		func(server *googleGrpc.Server) {
			pb.RegisterQueueServer(server, userController)
		},
	)

	// Add shutdown hook to trigger closer resources of service
	server.AddShutdownHook(grpcServer)
}
