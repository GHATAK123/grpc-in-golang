package main

import (
	"context"
	pb "grpc-in-golang/coffeeshop_proto"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to gRPC server")
	}
	defer conn.Close()

	c := pb.NewCoeffeShopClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	menuStream, err := c.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatal("error calling GetMenu")
	}

	done := make(chan bool)

	var items []*pb.Item

	go func() {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			items = resp.Items
			log.Printf("Resp recived %v", resp.Items)
		}
	}()

	<-done

	receipt, _ := c.PlaceOrder(ctx, &pb.Order{Items: items})
	log.Printf("receipt : %v", receipt)

	status, _ := c.GetOrderStatus(ctx, receipt)
	log.Printf("status : %v", status)

}
