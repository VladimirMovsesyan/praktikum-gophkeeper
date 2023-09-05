package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"praktikum-gophkeeper/pkg/auth"
	pb "praktikum-gophkeeper/proto"
	"sync"
)

type AuthServer struct {
	pb.UnimplementedAuthorizationServer
	users sync.Map
}

func (s *AuthServer) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	resp := &pb.RegisterUserResponse{}

	if _, ok := s.users.Load(in.User.Login); ok {
		return nil, status.Errorf(codes.AlreadyExists, "User with such login already exist.")
	}

	s.users.Store(in.User.Login, in.User.Password)

	token, err := auth.GenerateToken(in.User.Login)
	if err != nil {
		return nil, err
	}

	resp.Token = token

	return resp, nil
}

func (s *AuthServer) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	resp := &pb.LoginUserResponse{}

	if pass, ok := s.users.Load(in.User.Login); !ok {
		return nil, status.Errorf(codes.NotFound, "User with such login doesn't exist.")
	} else if in.User.Password != pass {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid login or password.")
	}

	token, err := auth.GenerateToken(in.User.Login)
	if err != nil {
		return nil, err
	}

	resp.Token = token

	return resp, nil
}
