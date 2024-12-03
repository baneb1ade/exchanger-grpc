package exchanger

import (
	"context"
	"database/sql"
	"errors"
	"exchanger-microservice/internal/domain/models"
	"log/slog"
)

type Service struct {
	logger  *slog.Logger
	storage Storage
}

func NewService(logger *slog.Logger, storage Storage) *Service {
	return &Service{logger, storage}
}

func (s *Service) GetExchangeRates(ctx context.Context) ([]models.Currency, error) {
	const op = "exchangerService.GetExchangeRates"
	log := s.logger.With(slog.String("op", op))

	cur, err := s.storage.GetAll(ctx)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return cur, nil
}

func (s *Service) GetExchangeRateForCurrency(ctx context.Context, fromCurrency, toCurrency string) (float32, error) {
	const op = "exchangerService.GetExchangeRates"
	log := s.logger.With(slog.String("op", op))

	fromCur, err := s.storage.GetOne(ctx, fromCurrency)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrCurrencyNotFound
		}
		log.Error(err.Error())
		return 0, err
	}
	toCur, err := s.storage.GetOne(ctx, toCurrency)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrCurrencyNotFound
		}
		log.Error(err.Error())
		return 0, err
	}

	if fromCur.Rate == 0 {
		return 0, ErrCurrencyRateIsZero
	}
	exchangeRate := toCur.Rate / fromCur.Rate

	return float32(exchangeRate), nil
}
