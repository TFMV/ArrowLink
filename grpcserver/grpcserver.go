package grpcserver

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/TFMV/ArrowLink/arrow"
	pb "github.com/TFMV/ArrowLink/proto/dataexchange"
	"go.uber.org/zap"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

// Server implements the ArrowDataServiceServer interface and holds an arrow service instance.
type Server struct {
	pb.ArrowDataServiceServer
	logger       *zap.Logger
	arrowService arrow.ArrowService
}

// NewServer creates a new Server instance.
func NewServer(logger *zap.Logger, arrowService arrow.ArrowService) *Server {
	return &Server{
		logger:       logger,
		arrowService: arrowService,
	}
}

// GetArrowData retrieves the Arrow data and streams it to the client.
func (s *Server) GetArrowData(req *pb.Empty, stream pb.ArrowDataService_GetArrowDataServer) error {
	data, err := s.arrowService.GetData()
	if err != nil {
		s.logger.Error("failed to get arrow data", zap.Error(err))
		return err
	}
	return stream.Send(&pb.ArrowData{Payload: data})
}

// StartGRPCServer sets up and runs the gRPC server with middleware and graceful shutdown.
func StartGRPCServer(address string, logger *zap.Logger, arrowService arrow.ArrowService) {
	grpc_zap.ReplaceGrpcLogger(logger)

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterArrowDataServiceServer(grpcServer, NewServer(logger, arrowService))

	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	// Run server in a separate goroutine.
	go func() {
		logger.Info("gRPC server is running", zap.String("address", address))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	// Listen for termination signals to gracefully shut down.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down gRPC server...")
	grpcServer.GracefulStop()
	logger.Info("gRPC server shutdown complete")
}
