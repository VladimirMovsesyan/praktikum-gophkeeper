package service

import (
	"context"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"praktikum-gophkeeper/pkg/storage"
	pb "praktikum-gophkeeper/proto"
)

type passwordRepository interface {
	Add(user string, password *pb.Password) error
	Get(user, website string) (passwords []*pb.Password, ids []uint32, err error)
	Update(user string, id int32, password *pb.Password) error
	Delete(user string, id int32) error
}

type textRepository interface {
	Add(user string, text *pb.Text) error
	Get(user, title string) (texts []*pb.Text, ids []uint32, err error)
	Update(user string, id uint32, text *pb.Text) error
	Delete(user string, id uint32) error
}

type binaryRepository interface {
	Add(user string, binary *pb.Binary) error
	Get(user, title string) (binaries []*pb.Binary, ids []uint32, err error)
	Update(user string, id uint32, binary *pb.Binary) error
	Delete(user string, id uint32) error
}

type paymentRepository interface {
	Add(user string, payment *pb.Payment) error
	Get(user, name string) (payments []*pb.Payment, ids []uint32, err error)
	Update(user string, id uint32, payment *pb.Payment) error
	Delete(user string, id uint32) error
}

type GophKeeperServer struct {
	pb.UnimplementedGophKeeperServer
	password passwordRepository
	text     textRepository
	binary   binaryRepository
	payment  paymentRepository
}

func NewGophKeeperServer(conn *pgx.Conn) (*GophKeeperServer, error) {
	pass, err := storage.NewPasswordStorage(conn)
	if err != nil {
		return nil, err
	}

	text, err := storage.NewTextStorage(conn)
	if err != nil {
		return nil, err
	}

	binary, err := storage.NewBinaryStorage(conn)
	if err != nil {
		return nil, err
	}

	payment, err := storage.NewPaymentStorage(conn)
	if err != nil {
		return nil, err
	}

	return &GophKeeperServer{
		password: pass,
		text:     text,
		binary:   binary,
		payment:  payment,
	}, nil
}

func (s *GophKeeperServer) AddPassword(ctx context.Context, in *pb.AddPasswordRequest) (*pb.AddPasswordResponse, error) {
	resp := &pb.AddPasswordResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.password.Add(login, in.Password)
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

	passwords, ids, err := s.password.Get(login, in.Website)
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

	err := s.password.Update(login, in.Id, in.Password)
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

	err := s.password.Delete(login, in.Id)
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

	err := s.text.Add(login, in.Text)
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

	texts, ids, err := s.text.Get(login, in.Title)
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

	err := s.text.Update(login, in.Id, in.Text)
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

	err := s.text.Delete(login, in.Id)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GophKeeperServer) AddBinary(ctx context.Context, in *pb.AddBinaryRequest) (*pb.AddBinaryResponse, error) {
	resp := &pb.AddBinaryResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.binary.Add(login, in.Binary)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GophKeeperServer) GetBinary(ctx context.Context, in *pb.GetBinaryRequest) (*pb.GetBinaryResponse, error) {
	resp := &pb.GetBinaryResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	binaries, ids, err := s.binary.Get(login, in.Title)
	if err != nil {
		return nil, err
	}

	resp.Binaries = binaries
	resp.Ids = ids

	return resp, nil
}

func (s *GophKeeperServer) UpdateBinary(ctx context.Context, in *pb.UpdateBinaryRequest) (*pb.UpdateBinaryResponse, error) {
	resp := &pb.UpdateBinaryResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.binary.Update(login, in.Id, in.Binary)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GophKeeperServer) DeleteBinary(ctx context.Context, in *pb.DeleteBinaryRequest) (*pb.DeleteBinaryResponse, error) {
	resp := &pb.DeleteBinaryResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.binary.Delete(login, in.Id)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GophKeeperServer) AddPayment(ctx context.Context, in *pb.AddPaymentRequest) (*pb.AddPaymentResponse, error) {
	resp := &pb.AddPaymentResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.payment.Add(login, in.Payment)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GophKeeperServer) GetPayment(ctx context.Context, in *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	resp := &pb.GetPaymentResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	payments, ids, err := s.payment.Get(login, in.Name)
	if err != nil {
		return nil, err
	}

	resp.Payments = payments
	resp.Ids = ids

	return resp, nil
}

func (s *GophKeeperServer) UpdatePayment(ctx context.Context, in *pb.UpdatePaymentRequest) (*pb.UpdatePaymentResponse, error) {
	resp := &pb.UpdatePaymentResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.payment.Update(login, in.Id, in.Payment)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GophKeeperServer) DeletePayment(ctx context.Context, in *pb.DeletePaymentRequest) (*pb.DeletePaymentResponse, error) {
	resp := &pb.DeletePaymentResponse{}

	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Error(codes.Internal, "Login value doesn't found in context")
	}

	err := s.payment.Delete(login, in.Id)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
