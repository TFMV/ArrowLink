package main

import (
	"github.com/TFMV/ArrowLink/grpcserver"
	"go.uber.org/zap"
)

func main() {
	// Create zap logger (production configuration)
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Start gRPC server using the new module
	grpcserver.StartGRPCServer(":50051", logger)
}
