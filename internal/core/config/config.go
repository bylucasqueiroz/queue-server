package config

import "google.golang.org/grpc/keepalive"

type GrpcServerConfig struct {
	Port            uint32
	KeepaliveParams keepalive.ServerParameters
	KeepalivePolicy keepalive.EnforcementPolicy
}
