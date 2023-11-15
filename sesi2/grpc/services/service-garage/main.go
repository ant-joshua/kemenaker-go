package main

import (
	"context"
	"grpc/models"
	"log"

	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

var localStorage *models.GarageListByUser

func init() {
	localStorage = new(models.GarageListByUser)
	localStorage.List = make(map[string]*models.GarageList)
}

type GarageServer struct {
	models.UnimplementedGaragesServer
}

func (GarageServer) List(ctx context.Context, params *models.GarageID) (*models.GarageList, error) {
	userId := params.Id

	return localStorage.List[userId], nil
}
func (GarageServer) Add(ctx context.Context, params *models.AddGarage) (*emptypb.Empty, error) {
	userId := params.Id
	garage := params.Garage

	if _, ok := localStorage.List[userId]; !ok {
		localStorage.List[userId] = new(models.GarageList)
		localStorage.List[userId].List = make([]*models.Garage, 0)
	}

	localStorage.List[userId].List = append(localStorage.List[userId].List, garage)

	log.Println("GarageServer.Add", garage.String(), "to", userId)

	return new(emptypb.Empty), nil
}
