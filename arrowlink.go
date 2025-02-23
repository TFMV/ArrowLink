package main

import (
	"log"
	"net"

	"bytes"

	pb "github.com/TFMV/ArrowLink/proto/dataexchange"
	arrow "github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	ipc "github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"google.golang.org/grpc"
)

type server struct {
	pb.ArrowDataServiceServer
}

func (s *server) GetArrowData(req *pb.Empty, stream pb.ArrowDataService_GetArrowDataServer) error {
	// Define Arrow schema
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "value", Type: arrow.PrimitiveTypes.Float64},
	}, nil)

	// Create record builder (for demo purposes, using in-memory data)
	// You would fill with your actual dataset here.
	mem := memory.NewGoAllocator()
	builder := array.NewRecordBuilder(mem, schema)
	defer builder.Release()

	// Populate data (example with one record)
	intBuilder := builder.Field(0).(*array.Int64Builder)
	floatBuilder := builder.Field(1).(*array.Float64Builder)

	intBuilder.Append(1)
	floatBuilder.Append(3.14)

	record := builder.NewRecord()
	defer record.Release()

	// Serialize the record using Arrow IPC
	var buf bytes.Buffer
	writer := ipc.NewWriter(&buf, ipc.WithSchema(schema))
	if err := writer.Write(record); err != nil {
		return err
	}
	writer.Close()

	// Stream the serialized payload
	payload := buf.Bytes()
	err := stream.Send(&pb.ArrowData{Payload: payload})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterArrowDataServiceServer(s, &server{})
	log.Println("Go gRPC server is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
