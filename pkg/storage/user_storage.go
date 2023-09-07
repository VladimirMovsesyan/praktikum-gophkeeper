package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "praktikum-gophkeeper/proto"
	"time"
)

type userStorage struct {
	conn *pgx.Conn
}

func NewUserStorage(conn *pgx.Conn) (*userStorage, error) {
	s := &userStorage{
		conn: conn,
	}

	err := s.ensureTableExist()
	if err != nil {
		return nil, err
	}

	return s, nil
}

const (
	userTable = `CREATE TABLE IF NOT EXISTS users (
    login VARCHAR(100) PRIMARY KEY,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL
);`
)

func (s *userStorage) ensureTableExist() error {
	_, err := s.conn.Exec(context.Background(), userTable)
	return err
}

func (s *userStorage) Add(user *pb.User) error {
	query := `INSERT INTO users VALUES($1, $2, $3)`

	_, err := s.conn.Exec(context.Background(), query, user.Login, user.Password, time.Now())
	if err != nil {
		return status.Errorf(codes.AlreadyExists, `User with login "%s" already exist`, user.Login)
	}

	return status.Error(codes.OK, "")
}

func (s *userStorage) Get(login string) (*pb.User, error) {
	query := `SELECT login, password FROM users WHERE login = $1`

	row := s.conn.QueryRow(context.Background(), query, login)

	user := &pb.User{}
	err := row.Scan(&user.Login, &user.Password)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, `User with login %s doesn't exist`, login)
	}

	return user, status.Error(codes.OK, "")
}
