package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "praktikum-gophkeeper/proto"
	"time"
)

type repository interface {
	AddUser(user *pb.User) error
	GetUser(login string) (*pb.User, error)
	AddPassword(user string, password *pb.Password) error
}

var _ repository = &PostgreStorage{}

type PostgreStorage struct {
	conn *pgx.Conn
}

func NewStorage(conn *pgx.Conn) (*PostgreStorage, error) {
	s := &PostgreStorage{
		conn: conn,
	}

	err := s.ensureTablesExist()
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

	passwordsTable = `CREATE TABLE IF NOT EXISTS passwords (
    id SERIAL PRIMARY KEY,
    website VARCHAR(100),
    login VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    owner VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (owner) REFERENCES users (login)
);`
)

func (s *PostgreStorage) ensureTablesExist() error {
	_, err := s.conn.Exec(context.Background(), userTable)
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(context.Background(), passwordsTable)
	return err
}

func (s *PostgreStorage) AddUser(user *pb.User) error {
	query := `INSERT INTO users VALUES($1, $2, $3)`

	_, err := s.conn.Exec(context.Background(), query, user.Login, user.Password, time.Now())
	if err != nil {
		return status.Errorf(codes.AlreadyExists, `User with login "%s" already exist`, user.Login)
	}

	return status.Error(codes.OK, "")
}

func (s *PostgreStorage) GetUser(login string) (*pb.User, error) {
	query := `SELECT login, password FROM users WHERE login = $1`

	row := s.conn.QueryRow(context.Background(), query, login)

	user := &pb.User{}
	err := row.Scan(&user.Login, &user.Password)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, `User with login %s doesn't exist`, login)
	}

	return user, status.Error(codes.OK, "")
}

func (s *PostgreStorage) AddPassword(user string, password *pb.Password) error {
	query := `INSERT INTO passwords(website, login, password, owner, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		password.Website,
		password.Login,
		password.Password,
		user,
		time.Now(),
	)
	if err != nil {
		return status.Error(codes.Internal, "Couldn't store password to database")
	}

	return status.Error(codes.OK, "")
}
