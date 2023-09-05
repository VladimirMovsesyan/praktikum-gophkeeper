package service

import (
	"context"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"praktikum-gophkeeper/pkg/storage"
	pb "praktikum-gophkeeper/proto"
)

type repository interface {
	AddUser(user *pb.User) error
	GetUser(login string) (*pb.User, error)
	AddPassword(user string, password *pb.Password) error
}

type GophKeeperServer struct {
	pb.UnimplementedGophKeeperServer
	storage repository
}

func NewGophKeeperServer(conn *pgx.Conn) (*GophKeeperServer, error) {
	s, err := storage.NewStorage(conn)
	if err != nil {
		return nil, err
	}

	return &GophKeeperServer{storage: s}, nil
}

func (s *GophKeeperServer) AddPassword(ctx context.Context, in *pb.AddPasswordRequest) (*pb.AddPasswordResponse, error) {
	resp := &pb.AddPasswordResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.storage.AddPassword(login, in.Password)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
