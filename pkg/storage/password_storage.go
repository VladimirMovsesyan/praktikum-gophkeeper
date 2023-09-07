package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "praktikum-gophkeeper/proto"
	"time"
)

type passwordStorage struct {
	conn *pgx.Conn
}

func NewPasswordStorage(conn *pgx.Conn) (*passwordStorage, error) {
	s := &passwordStorage{
		conn: conn,
	}

	err := s.ensureTableExist()
	if err != nil {
		return nil, err
	}

	return s, nil
}

const (
	passwordTable = `CREATE TABLE IF NOT EXISTS passwords (
    id SERIAL PRIMARY KEY,
    website VARCHAR(100),
    login VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    owner VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (owner) REFERENCES users (login)
);`
)

func (s *passwordStorage) ensureTableExist() error {
	_, err := s.conn.Exec(context.Background(), passwordTable)

	return err
}

func (s *passwordStorage) Add(user string, password *pb.Password) error {
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

func (s *passwordStorage) Get(user, website string) (passwords []*pb.Password, ids []uint32, err error) {
	query := `SELECT website, login, password, id FROM passwords WHERE owner = $1 AND website = $2`

	rows, err := s.conn.Query(context.Background(), query, user, website)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "Something went wrong while making request to database")
	}

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

func (s *passwordStorage) Update(user string, id int32, password *pb.Password) error {
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

func (s *passwordStorage) Delete(user string, id int32) error {
	query := `DELETE FROM passwords WHERE owner = $1 AND id = $2`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		user,
		id,
	)

	return err
}
