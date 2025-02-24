package main

import (
	"fmt"
	"os"

	"github.com/TFMV/ArrowLink/arrow"
	"github.com/TFMV/ArrowLink/grpcserver"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		// Here we avoid panic to make debugging easier
		fmt.Fprintln(os.Stderr, "Failed to initialize logger:", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Create the arrow service
	arrowService := arrow.NewArrowService()

	// Start the gRPC server with the arrow service injected
	grpcserver.StartGRPCServer(":50051", logger, arrowService)
}
