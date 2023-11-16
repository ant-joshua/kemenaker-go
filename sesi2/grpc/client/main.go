package main

import (
	"context"
	"encoding/json"
	"grpc/models"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func serviceGarage() models.GaragesClient {
	conn, err := grpc.Dial("localhost:9002", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("could not connect to localhost:9002", err.Error())
	}

	return models.NewGaragesClient(conn)
}

func serviceUser() models.UsersClient {
	conn, err := grpc.Dial("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("could not connect to localhost:9001", err.Error())
	}

	return models.NewUsersClient(conn)
}

func main() {
	user1 := models.User{
		Id:       "001",
		Name:     "user1",
		Password: "123456",
		Gender:   models.UserGender_FEMALE,
	}

	user2 := models.User{
		Id:       "002",
		Name:     "user2",
		Password: "123456",
		Gender:   models.UserGender_MALE,
	}

	garage1 := models.Garage{
		Id:   "G001",
		Name: "garage1",
		Coordinate: &models.GarageCoordinate{
			Latitude:  1.0,
			Longitude: 1.0,
		},
	}

	garage2 := models.Garage{
		Id:   "G002",
		Name: "garage2",
		Coordinate: &models.GarageCoordinate{
			Latitude:  1.0,
			Longitude: 1.0,
		},
	}

	user := serviceUser()

	user.Register(context.Background(), &user1)
	user.Register(context.Background(), &user2)

	res1, err := user.List(context.Background(), new(emptypb.Empty))

	if err != nil {
		log.Fatal("could not list users", err.Error())
	}

	res1String, _ := json.Marshal(res1.List)
	log.Println("users", string(res1String))

	garage := serviceGarage()

	garage.Add(context.Background(), &models.AddGarage{
		Id:     "001",
		Garage: &garage1,
	})

	garage.Add(context.Background(), &models.AddGarage{
		Id:     "002",
		Garage: &garage2,
	})

	res2, err := garage.List(context.Background(), &models.GarageID{
		Id: "001",
	})

	if err != nil {
		log.Fatal("could not list users", err.Error())
	}

	res2String, _ := json.Marshal(res2.List)
	log.Println("")
	log.Println("garage", string(res2String))
}
