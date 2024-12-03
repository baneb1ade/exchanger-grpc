package exchanger

import (
	"context"
	"exchanger-microservice/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SQLClient interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Storage interface {
	GetOne(ctx context.Context, currency string) (models.Currency, error)
	GetAll(ctx context.Context) ([]models.Currency, error)
}
