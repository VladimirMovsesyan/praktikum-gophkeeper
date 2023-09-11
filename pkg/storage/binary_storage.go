package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	pb "praktikum-gophkeeper/proto"
	"time"
)

type binaryStorage struct {
	conn *pgx.Conn
}

func NewBinaryStorage(conn *pgx.Conn) (*binaryStorage, error) {
	s := &binaryStorage{
		conn: conn,
	}

	err := s.ensureTableExist()
	if err != nil {
		return nil, err
	}

	return s, nil
}

const (
	binaryTable = `CREATE TABLE IF NOT EXISTS binaries (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    file bytea NOT NULL,
    owner VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (owner) REFERENCES users (login)
);`
)

func (s *binaryStorage) ensureTableExist() error {
	_, err := s.conn.Exec(context.Background(), binaryTable)
	return err
}

func (s *binaryStorage) Add(user string, binary *pb.Binary) error {
	query := `INSERT INTO binaries(title, file, owner, created_at) VALUES($1, $2, $3, $4)`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		binary.Title,
		binary.File,
		user,
		time.Now(),
	)

	return err
}

func (s *binaryStorage) Get(user, title string) (binaries []*pb.Binary, ids []uint32, err error) {
	query := `SELECT title, file, id FROM binaries WHERE owner = $1 AND title = $2`

	rows, err := s.conn.Query(
		context.Background(),
		query,
		user,
		title,
	)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		binary := &pb.Binary{}
		var id uint32
		err := rows.Scan(&binary.Title, &binary.File, &id)
		if err != nil {
			return nil, nil, err
		}

		ids = append(ids, id)
		binaries = append(binaries, binary)
	}

	return binaries, ids, nil
}

func (s *binaryStorage) Update(user string, id uint32, binary *pb.Binary) error {
	query := `UPDATE binaries SET title = $1, file = $2 WHERE owner = $3 AND id = $4`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		binary.Title,
		binary.File,
		user,
		id,
	)

	return err
}

func (s *binaryStorage) Delete(user string, id uint32) error {
	query := `DELETE FROM binaries WHERE owner = $1 AND id = $2`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		user,
		id,
	)

	return err
}
