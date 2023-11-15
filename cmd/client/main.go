package main

import (
	"context"
	"encoding/json"
	"log"
	"testGRPC/pkg/api"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := api.NewDirectoryInformerClient(conn)
	res, err := c.Dir(context.Background(), &api.DirectoryRequest{
		Path: "./cmd",
	})
	if err != nil {
		log.Fatal(err)
	}

	bytes := res.GetDirectoryInfo()

	desc := api.Description{}
	if err = json.Unmarshal(bytes, &desc); err != nil {
		log.Fatal(err)
	}

	log.Println(desc)
}
