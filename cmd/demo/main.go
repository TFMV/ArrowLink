package main

import (
	"flag"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/TFMV/ArrowLink/arrow"
	"github.com/TFMV/ArrowLink/grpcserver"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Parse command line flags
	dataSize := flag.Int("size", 1000, "Number of rows to generate")
	benchmark := flag.Bool("benchmark", false, "Run performance benchmark")
	visualize := flag.Bool("visualize", false, "Generate visualization of data")
	insecure := flag.Bool("insecure", true, "Use insecure connection (no TLS)")
	flag.Parse()

	// Create the arrow service with configurable data size
	arrowService := arrow.NewDemoArrowService(*dataSize)

	// Start the gRPC server in a goroutine
	go grpcserver.StartGRPCServer(":50051", logger, arrowService)
	logger.Info("Server started on :50051")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Run the Python client
	pythonPath := "python"
	if runtime.GOOS == "windows" {
		pythonPath = "python.exe"
	}

	clientArgs := []string{"/Users/thomasmcgeehan/ArrowLink/ArrowLink/python/main.py"}
	if *benchmark {
		clientArgs = append(clientArgs, "--benchmark")
	}
	if *visualize {
		clientArgs = append(clientArgs, "--visualize")
	}
	if *insecure {
		// Use a non-existent cert path to force insecure mode
		clientArgs = append(clientArgs, "--cert", "non_existent_cert.crt")
	}

	cmd := exec.Command(pythonPath, clientArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Info("Running Python client...")
	if err := cmd.Run(); err != nil {
		logger.Error("Failed to run Python client", zap.Error(err))
		os.Exit(1)
	}
}
