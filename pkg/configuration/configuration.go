package configuration

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
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
	DB      *pgx.Conn
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

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return Server{}, err
	}

	srv := grpc.NewServer()

	authServer, err := service.NewAuthServer(conn)
	if err != nil {
		return Server{}, err
	}
	pb.RegisterAuthorizationServer(srv, authServer)

	gophkeeperServer, err := service.NewGophKeeperServer(conn)
	if err != nil {
		return Server{}, err
	}
	pb.RegisterGophKeeperServer(srv, gophkeeperServer)

	return Server{
		Address: address,
		DSN:     dsn,
		DB:      conn,
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
