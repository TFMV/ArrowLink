package main

import (
	"log"
	"net"

	"github.com/TFMV/ArrowLink/arrow"
	pb "github.com/TFMV/ArrowLink/proto/dataexchange"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedArrowDataServiceServer
	logger       *zap.Logger
	arrowService arrow.ArrowService
}

func (s *server) GetArrowData(req *pb.Empty, stream pb.ArrowDataService_GetArrowDataServer) error {
	data, err := s.arrowService.GetData()
	if err != nil {
		s.logger.Error("failed to get arrow data", zap.Error(err))
		return err
	}

	return stream.Send(&pb.ArrowData{Payload: data})
}

func main() {
	// Create zap logger (production configuration)
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Create arrow service
	arrowService := arrow.NewArrowService()

	// Create listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()
	pb.RegisterArrowDataServiceServer(s, &server{
		logger:       logger,
		arrowService: arrowService,
	})

	logger.Info("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
