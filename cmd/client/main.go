package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "praktikum-gophkeeper/proto"
)

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := pb.NewAuthorizationClient(conn)
	Test(client)
}

func Test(client pb.AuthorizationClient) {
	cc := []*pb.User{
		{Login: "Login1", Password: "Password1"},
		{Login: "Login2", Password: "Password2"},
		{Login: "Login3", Password: "Password3"},
		{Login: "Login1", Password: "Password1"},
	}

	for _, c := range cc {
		resp, err := client.RegisterUser(context.Background(), &pb.RegisterUserRequest{User: c})
		if err != nil {
			log.Println(err)
		}
		log.Println(resp.GetToken())
	}

	log.Println("Login")

	for _, c := range cc {
		resp, err := client.LoginUser(context.Background(), &pb.LoginUserRequest{User: c})
		if err != nil {
			log.Println(err)
		}
		log.Println(resp.GetToken())
	}
}
