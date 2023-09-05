package configuration

import (
	"errors"
	"google.golang.org/grpc"
	"os"
	"praktikum-gophkeeper/pkg/service"
	pb "praktikum-gophkeeper/proto"
)

const (
	envAddress = "RUN_ADDRESS"
	envDSN     = "DSN"
)

type Server struct {
	Address string
	DSN     string
	Server  *grpc.Server
}

func NewServer(flAddress, flDSN *string) (Server, error) {
	address, err := parseStringVar(flAddress, envAddress)
	if err != nil {
		return Server{}, err
	}

	dsn, err := parseStringVar(flDSN, envDSN)
	if err != nil {
		return Server{}, err
	}

	srv := grpc.NewServer()
	pb.RegisterAuthorizationServer(srv, &service.AuthServer{})

	return Server{
		Address: address,
		DSN:     dsn,
		Server:  srv,
	}, nil
}

func parseStringVar(flag *string, envName string) (string, error) {
	if *flag != "" {
		return *flag, nil
	}

	value := os.Getenv(envName)
	if value == "" {
		return "", errors.New("not enough parameters to run service")
	}
	return value, nil
}
