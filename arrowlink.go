package main

import (
	"github.com/TFMV/ArrowLink/arrow"
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

	// Create arrow service
	arrowService := arrow.NewArrowService()

	// Start the gRPC server with our arrowService injected.
	grpcserver.StartGRPCServer(":50051", logger, arrowService)
}
