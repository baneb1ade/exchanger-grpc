package db

import (
	"context"
	"exchanger-microservice/internal/domain/exchanger"
	"exchanger-microservice/internal/domain/models"
	"log/slog"
)

type Storage struct {
	Client exchanger.SQLClient
	Logger *slog.Logger
}

func NewStorage(client exchanger.SQLClient, logger *slog.Logger) *Storage {
	return &Storage{client, logger}
}

func (s *Storage) GetOne(ctx context.Context, currency string) (models.Currency, error) {
	const op = "db.psql.GetOne"
	log := s.Logger.With(slog.String("op", op))

	q := `SELECT code, rate
			FROM currency
			WHERE code = $1`

	var res models.Currency
	if err := s.Client.QueryRow(ctx, q, currency).Scan(&res.Code, &res.Rate); err != nil {
		log.Error(op, "error", err)
		return res, err
	}
	return res, nil
}

func (s *Storage) GetAll(ctx context.Context) ([]models.Currency, error) {
	const op = "db.psql.GetAll"
	log := s.Logger.With(slog.String("op", op))

	q := `SELECT code, rate FROM currency`

	var res []models.Currency
	rows, err := s.Client.Query(ctx, q)
	if err != nil {
		log.Error(op, "error", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var currency models.Currency
		if err := rows.Scan(&currency.Code, &currency.Rate); err != nil {
			log.Error(op, "error", err)
			return nil, err
		}
		res = append(res, currency)
	}
	return res, nil
}

//func (s *Storage) GetOneByID(ctx context.Context, id string) (models.User, error) {
//	const op = "db.psql.GetOneByID"
//	log := s.Logger.With(slog.String("op", op))
//
//	q := `SELECT id, username, email, password
//          FROM "user"
//          WHERE id = $1`
//
//	var u models.User
//	if err := s.Client.QueryRow(ctx, q, id).Scan(&u.ID, &u.Username, &u.Email, &u.PassHash); err != nil {
//		log.Error(op, "error", err)
//		return u, err
//	}
//	return u, nil
//}
