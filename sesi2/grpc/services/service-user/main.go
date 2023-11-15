package main

import (
	"context"
	"grpc/models"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var localStorage *models.UserList

func init() {
	localStorage = new(models.UserList)
	localStorage.List = make([]*models.User, 0)
}

type UserServer struct {
	models.UnimplementedUsersServer
}

func (UserServer) Register(ctx context.Context, params *models.User) (*emptypb.Empty, error) {
	localStorage.List = append(localStorage.List, params)

	log.Println("UserServer.Register", params.String())

	return new(emptypb.Empty), nil
}

func (UserServer) List(context.Context, *emptypb.Empty) (*models.UserList, error) {
	return localStorage, nil
}

func main() {
	srv := grpc.NewServer()

	var userServ UserServer

	models.RegisterUsersServer(srv, userServ)

	log.Println("Starting RPC Server at", "localhost:9001")

	listen, err := net.Listen("tcp", "localhost:9001")

	if err != nil {
		log.Fatalf("could not listen to localhost:9001: %v", err)
	}

	log.Fatal(srv.Serve(listen))
}
