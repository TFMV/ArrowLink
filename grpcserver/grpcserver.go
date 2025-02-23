package grpcserver

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	pb "github.com/TFMV/ArrowLink/proto/dataexchange"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedArrowDataServiceServer
	logger *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{logger: logger}
}

func StartGRPCServer(address string, logger *zap.Logger) {
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
	pb.RegisterArrowDataServiceServer(grpcServer, NewServer(logger))

	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	go func() {
		logger.Info("gRPC server is running", zap.String("address", address))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down gRPC server...")
	grpcServer.GracefulStop()
	logger.Info("gRPC server shutdown complete")
}
