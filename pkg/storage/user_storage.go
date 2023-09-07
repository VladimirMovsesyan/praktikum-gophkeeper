package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
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
	return err
}

func (s *userStorage) Get(login string) (*pb.User, error) {
	query := `SELECT login, password FROM users WHERE login = $1`

	row := s.conn.QueryRow(context.Background(), query, login)

	user := &pb.User{}
	err := row.Scan(&user.Login, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
