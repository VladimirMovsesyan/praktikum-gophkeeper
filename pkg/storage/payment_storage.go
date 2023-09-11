package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	pb "praktikum-gophkeeper/proto"
	"time"
)

type paymentStorage struct {
	conn *pgx.Conn
}

func NewPaymentStorage(conn *pgx.Conn) (*paymentStorage, error) {
	s := &paymentStorage{
		conn: conn,
	}

	err := s.ensureTableExist()
	if err != nil {
		return nil, err
	}

	return s, nil
}

const (
	paymentTable = `CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    cardholder VARCHAR(100) NOT NULL,
    number VARCHAR(100) NOT NULL,
    exp_date VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL,
    owner VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (owner) REFERENCES users (login)
);`
)

func (s *paymentStorage) ensureTableExist() error {
	_, err := s.conn.Exec(context.Background(), paymentTable)
	return err
}

func (s *paymentStorage) Add(user string, payment *pb.Payment) error {
	query := `INSERT INTO payments(name, cardholder, number, exp_date, code, owner, created_at) VALUES($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		payment.Name,
		payment.Cardholder,
		payment.Number,
		payment.ExpDate,
		payment.Code,
		user,
		time.Now(),
	)

	return err
}

func (s *paymentStorage) Get(user, name string) (payments []*pb.Payment, ids []uint32, err error) {
	query := `SELECT (name, cardholder, number, exp_date, code, id) FROM payments WHERE owner = $1 AND name = $2`

	rows, err := s.conn.Query(
		context.Background(),
		query,
		user,
		name,
	)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		payment := &pb.Payment{}
		var id uint32
		err := rows.Scan(&payment.Name, &payment.Cardholder, &payment.Number, &payment.ExpDate, &payment.Code, &id)
		if err != nil {
			return nil, nil, err
		}

		payments = append(payments, payment)
		ids = append(ids, id)
	}

	return payments, ids, nil
}

func (s *paymentStorage) Update(user string, id uint32, payment *pb.Payment) error {
	query := `UPDATE payments SET name = $1, cardholder = $2, number = $3, exp_date = $4, code = $5 WHERE owner = $6 AND id = $7`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		payment.Name,
		payment.Cardholder,
		payment.Number,
		payment.ExpDate,
		payment.Code,
		user,
		id,
	)

	return err
}

func (s *paymentStorage) Delete(user string, id uint32) error {
	query := `DELETE FROM payments WHERE owner = $1 AND id = $2`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		user,
		id,
	)

	return err
}
