package main

import (
	"bytes"
	"fmt"
	"grpc/models"
	"os"
	"strings"

	"github.com/golang/protobuf/jsonpb"
)

func main() {

	var user1 = &models.User{
		Id:       "001",
		Name:     "user1",
		Password: "123456",
		Gender:   models.UserGender_MALE,
	}

	var garage1 = &models.Garage{
		Id:   "001",
		Name: "garage1",
		Coordinate: &models.GarageCoordinate{
			Latitude:  1.0,
			Longitude: 1.0,
		},
	}

	fmt.Println(user1.Name)
	fmt.Println(garage1.Name)

	fmt.Printf("%v\n", user1)

	fmt.Printf("%v\n", user1.String())

	var buf bytes.Buffer

	err1 := (&jsonpb.Marshaler{}).Marshal(&buf, user1)

	if err1 != nil {
		fmt.Println(err1.Error())
		os.Exit(0)
	}

	jsonString := buf.String()
	fmt.Printf("# as json\n %v\n", jsonString)

	buf2 := strings.NewReader(jsonString)

	protoObject := new(models.User)

	err2 := jsonpb.Unmarshal(buf2, protoObject)

	if err2 != nil {
		fmt.Println(err2.Error())
		os.Exit(0)
	}

	fmt.Printf("# as string\n %s\n", protoObject.String())

	// var userList = &models.UserList{
	// 	List: ,
	// }
}
