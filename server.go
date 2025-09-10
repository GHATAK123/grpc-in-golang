package main

import (
	"context"
	pb "grpc-in-golang/coffeeshop_proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoeffeShopServer
}

func (s *server) GetMenu(MenuRequest *pb.MenuRequest, srv pb.CoeffeShop_GetMenuServer) error {
	items := []*pb.Item{
		&pb.Item{
			Id:   "1",
			Name: "Latte",
		},
		&pb.Item{
			Id:   "2",
			Name: "Indian Coffee",
		},
		&pb.Item{
			Id:   "3",
			Name: "Desi Chai",
		},
	}

	for i, _ := range items {
		srv.Send(&pb.Menu{
			Items: items[0 : i+1],
		})
	}

	return nil

}

func (s *server) PlaceOrder(context context.Context, order *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{
		Id: "Rcp123",
	}, nil

}
func (s *server) GetOrderStatus(context context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{
		OrderId: receipt.Id,
		Status:  "IN PROGRESS",
	}, nil

}

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterCoeffeShopServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : % s", err)
	}

	log.Println("Server is running at port 9001")

}
