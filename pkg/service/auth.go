package service

import (
	"context"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"praktikum-gophkeeper/pkg/auth"
	"praktikum-gophkeeper/pkg/storage"
	pb "praktikum-gophkeeper/proto"
)

type AuthServer struct {
	pb.UnimplementedAuthorizationServer
	storage repository
}

func NewAuthServer(conn *pgx.Conn) (*AuthServer, error) {
	s, err := storage.NewStorage(conn)
	if err != nil {
		return nil, err
	}

	return &AuthServer{storage: s}, nil
}

func (s *AuthServer) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	resp := &pb.RegisterUserResponse{}

	if _, err := s.storage.GetUser(in.User.Login); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, `User with login "%s" already exist`, in.User.Login)
	}

	err := s.storage.AddUser(in.User)
	if err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(in.User.Login)
	if err != nil {
		return nil, err
	}

	resp.Token = token

	return resp, nil
}

func (s *AuthServer) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	resp := &pb.LoginUserResponse{}

	if user, err := s.storage.GetUser(in.User.Login); err != nil {
		return nil, err
	} else if in.User.Password != user.Password {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid login or password.")
	}

	token, err := auth.GenerateToken(in.User.Login)
	if err != nil {
		return nil, err
	}

	resp.Token = token

	return resp, nil
}
