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
	GetPassword(user, website string) ([]*pb.Password, []uint32, error)
	UpdatePassword(user string, id int32, password *pb.Password) error
	DeletePassword(user string, id int32) error

	AddText(user string, text *pb.Text) error
	GetText(user, title string) ([]*pb.Text, []uint32, error)
	UpdateText(user string, id uint32, text *pb.Text) error
	DeleteText(user string, id uint32) error
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

	passwordTable = `CREATE TABLE IF NOT EXISTS passwords (
    id SERIAL PRIMARY KEY,
    website VARCHAR(100),
    login VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    owner VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (owner) REFERENCES users (login)
);`

	textTable = `CREATE TABLE IF NOT EXISTS texts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    text VARCHAR(1000) NOT NULL,
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

	_, err = s.conn.Exec(context.Background(), passwordTable)
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(context.Background(), textTable)
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

func (s *PostgreStorage) GetPassword(user, website string) ([]*pb.Password, []uint32, error) {
	query := `SELECT website, login, password, id FROM passwords WHERE owner = $1 AND website = $2`

	rows, err := s.conn.Query(context.Background(), query, user, website)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "Something went wrong while making request to database")
	}

	var passwords []*pb.Password
	var ids []uint32
	for rows.Next() {
		pass := &pb.Password{}
		var id uint32
		err := rows.Scan(&pass.Website, &pass.Login, &pass.Password, &id)
		if err != nil {
			return nil, nil, status.Error(codes.Internal, "Something went wrong while scanning values from database")
		}

		passwords = append(passwords, pass)
		ids = append(ids, id)
	}

	return passwords, ids, nil
}

func (s *PostgreStorage) UpdatePassword(user string, id int32, password *pb.Password) error {
	query := `UPDATE passwords SET website = $1, login = $2, password = $3 WHERE id = $4 AND owner = $5`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		password.Website,
		password.Login,
		password.Password,
		id,
		user,
	)

	return err
}

func (s *PostgreStorage) DeletePassword(user string, id int32) error {
	query := `DELETE FROM passwords WHERE owner = $1 AND id = $2`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		user,
		id,
	)

	return err
}

func (s *PostgreStorage) AddText(user string, text *pb.Text) error {
	query := `INSERT INTO texts(title, text, owner, created_at) VALUES($1, $2, $3, $4)`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		text.Title,
		text.Text,
		user,
		time.Now(),
	)
	if err != nil {
		return status.Error(codes.Internal, "Couldn't store text to database")
	}

	return status.Error(codes.OK, "")
}

func (s *PostgreStorage) GetText(user, title string) ([]*pb.Text, []uint32, error) {
	query := `SELECT (title, text, id) FROM texts WHERE owner = $1 AND title = $2`

	rows, err := s.conn.Query(context.Background(), query, user, title)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "Something went wrong while making request to database")
	}

	var texts []*pb.Text
	var ids []uint32

	for rows.Next() {
		text := &pb.Text{}
		var id uint32
		err := rows.Scan(&text.Title, &text.Text, &id)
		if err != nil {
			return nil, nil, status.Error(codes.Internal, "Something went wrong while scanning values from database")
		}

		texts = append(texts, text)
		ids = append(ids, id)
	}

	return texts, ids, nil
}

func (s *PostgreStorage) UpdateText(user string, id uint32, text *pb.Text) error {
	query := `UPDATE texts SET title = $1, text = $2 WHERE owner = $3 AND id = $4`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		text.Title,
		text.Text,
		user,
		id,
	)

	return err
}

func (s *PostgreStorage) DeleteText(user string, id uint32) error {
	query := `DELETE FROM texts WHERE owner = $1 AND id = $2`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		user,
		id,
	)

	return err
}
