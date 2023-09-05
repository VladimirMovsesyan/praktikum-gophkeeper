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
	GetPassword(user, website string) ([]*pb.Password, []uint32, error)
	UpdatePassword(user string, id int32, password *pb.Password) error
	DeletePassword(user string, id int32) error

	AddText(user string, text *pb.Text) error
	GetText(user, title string) ([]*pb.Text, []uint32, error)
	UpdateText(user string, id uint32, text *pb.Text) error
	DeleteText(user string, id uint32) error
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

func (s *GophKeeperServer) GetPassword(ctx context.Context, in *pb.GetPasswordRequest) (*pb.GetPasswordResponse, error) {
	resp := &pb.GetPasswordResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	passwords, ids, err := s.storage.GetPassword(login, in.Website)
	if err != nil {
		return nil, err
	}

	resp.Passwords = passwords
	resp.Ids = ids

	return resp, nil
}

func (s *GophKeeperServer) UpdatePassword(ctx context.Context, in *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	resp := &pb.UpdatePasswordResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.storage.UpdatePassword(login, in.Id, in.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "Couldn't update password in database")
	}

	return resp, nil
}

func (s *GophKeeperServer) DeletePassword(ctx context.Context, in *pb.DeletePasswordRequest) (*pb.DeletePasswordResponse, error) {
	resp := &pb.DeletePasswordResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.storage.DeletePassword(login, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Couldn't delete password from database")
	}

	return resp, nil
}

func (s *GophKeeperServer) AddText(ctx context.Context, in *pb.AddTextRequest) (*pb.AddTextResponse, error) {
	resp := &pb.AddTextResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.storage.AddText(login, in.Text)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GophKeeperServer) GetText(ctx context.Context, in *pb.GetTextRequest) (*pb.GetTextResponse, error) {
	resp := &pb.GetTextResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	texts, ids, err := s.storage.GetText(login, in.Title)
	if err != nil {
		return nil, err
	}

	resp.Texts = texts
	resp.Ids = ids

	return resp, nil
}

func (s *GophKeeperServer) UpdateText(ctx context.Context, in *pb.UpdateTextRequest) (*pb.UpdateTextResponse, error) {
	resp := &pb.UpdateTextResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.storage.UpdateText(login, in.Id, in.Text)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GophKeeperServer) DeleteText(ctx context.Context, in *pb.DeleteTextRequest) (*pb.DeleteTextResponse, error) {
	resp := &pb.DeleteTextResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.storage.DeleteText(login, in.Id)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
