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

type userRepository interface {
	Add(user *pb.User) error
	Get(login string) (*pb.User, error)
}

type AuthServer struct {
	pb.UnimplementedAuthorizationServer
	user userRepository
}

func NewAuthServer(conn *pgx.Conn) (*AuthServer, error) {
	s, err := storage.NewUserStorage(conn)
	if err != nil {
		return nil, err
	}

	return &AuthServer{user: s}, nil
}

func (s *AuthServer) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	resp := &pb.RegisterUserResponse{}

	if _, err := s.user.Get(in.User.Login); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, `User with login "%s" already exist`, in.User.Login)
	}

	err := s.user.Add(in.User)
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

	if user, err := s.user.Get(in.User.Login); err != nil {
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
