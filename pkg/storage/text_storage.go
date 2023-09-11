package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	pb "praktikum-gophkeeper/proto"
	"time"
)

type textStorage struct {
	conn *pgx.Conn
}

func NewTextStorage(conn *pgx.Conn) (*textStorage, error) {
	s := &textStorage{
		conn: conn,
	}

	err := s.ensureTableExist()
	if err != nil {
		return nil, err
	}

	return s, nil
}

const (
	textTable = `CREATE TABLE IF NOT EXISTS texts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    text VARCHAR(1000) NOT NULL,
    owner VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (owner) REFERENCES users (login)
);`
)

func (s *textStorage) ensureTableExist() error {
	_, err := s.conn.Exec(context.Background(), textTable)

	return err
}

func (s *textStorage) Add(user string, text *pb.Text) error {
	query := `INSERT INTO texts(title, text, owner, created_at) VALUES($1, $2, $3, $4)`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		text.Title,
		text.Text,
		user,
		time.Now(),
	)

	return err
}

func (s *textStorage) Get(user, title string) (texts []*pb.Text, ids []uint32, err error) {
	query := `SELECT (title, text, id) FROM texts WHERE owner = $1 AND title = $2`

	rows, err := s.conn.Query(context.Background(), query, user, title)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		text := &pb.Text{}
		var id uint32
		err := rows.Scan(&text.Title, &text.Text, &id)
		if err != nil {
			return nil, nil, err
		}

		texts = append(texts, text)
		ids = append(ids, id)
	}

	return texts, ids, nil
}

func (s *textStorage) Update(user string, id uint32, text *pb.Text) error {
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

func (s *textStorage) Delete(user string, id uint32) error {
	query := `DELETE FROM texts WHERE owner = $1 AND id = $2`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		user,
		id,
	)

	return err
}
