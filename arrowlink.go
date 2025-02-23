package main

import (
	"bytes"
	"log"
	"net"

	pb "github.com/TFMV/ArrowLink/proto/dataexchange"
	arrow "github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	ipc "github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedArrowDataServiceServer
	logger *zap.Logger
}

func (s *server) GetArrowData(req *pb.Empty, stream pb.ArrowDataService_GetArrowDataServer) error {
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

func main() {
	// Create zap logger (production configuration)
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Create listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()
	pb.RegisterArrowDataServiceServer(s, &server{logger: logger})

	logger.Info("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
