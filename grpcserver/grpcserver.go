package grpcserver

import (
	"bytes"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	pb "github.com/TFMV/ArrowLink/proto/dataexchange"
	arrow "github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	ipc "github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedArrowDataServiceServer
	logger *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{logger: logger}
}

func (s *Server) GetArrowData(req *pb.Empty, stream pb.ArrowDataService_GetArrowDataServer) error {
	// Define Arrow schema
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "value", Type: arrow.PrimitiveTypes.Float64},
	}, nil)

	// Create record builder
	mem := memory.NewGoAllocator()
	builder := array.NewRecordBuilder(mem, schema)
	defer builder.Release()

	// Populate data
	intBuilder := builder.Field(0).(*array.Int64Builder)
	floatBuilder := builder.Field(1).(*array.Float64Builder)

	intBuilder.Append(1)
	floatBuilder.Append(3.14)

	record := builder.NewRecord()
	defer record.Release()

	// Serialize record
	var buf bytes.Buffer
	writer := ipc.NewWriter(&buf, ipc.WithSchema(schema))
	if err := writer.Write(record); err != nil {
		return err
	}
	writer.Close()

	return stream.Send(&pb.ArrowData{Payload: buf.Bytes()})
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
