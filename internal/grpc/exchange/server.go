package exchange

import (
	"context"
	"errors"
	"exchanger-microservice/internal/domain/exchanger"
	"exchanger-microservice/internal/domain/models"
	proto "github.com/baneb1ade/exchanger-protos/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	GetExchangeRates(ctx context.Context) ([]models.Currency, error)
	GetExchangeRateForCurrency(ctx context.Context, fromCurrency, toCurrency string) (float32, error)
}

type serverAPI struct {
	proto.UnimplementedExchangeServiceServer
	service Service
}

func Register(gRPC *grpc.Server, service Service) {
	proto.RegisterExchangeServiceServer(gRPC, &serverAPI{service: service})
}

func (s *serverAPI) GetExchangeRates(ctx context.Context, req *proto.Empty) (*proto.ExchangeRatesResponse, error) {
	res := make(map[string]float32)
	rates, err := s.service.GetExchangeRates(ctx)
	if err != nil {
		return nil, validateError(err)
	}
	for _, rate := range rates {
		res[rate.Code] = float32(rate.Rate)
	}
	return &proto.ExchangeRatesResponse{Rates: res}, nil
}

func (s *serverAPI) GetExchangeRateForCurrency(ctx context.Context, req *proto.CurrencyRequest) (*proto.ExchangeRateResponse, error) {
	rate, err := s.service.GetExchangeRateForCurrency(ctx, req.GetFromCurrency(), req.GetToCurrency())
	if err != nil {
		return nil, validateError(err)
	}
	return &proto.ExchangeRateResponse{
		FromCurrency: req.GetFromCurrency(),
		ToCurrency:   req.GetToCurrency(),
		Rate:         rate,
	}, nil
}

func validateError(err error) error {
	switch {
	case errors.Is(err, exchanger.ErrCurrencyNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, exchanger.ErrCurrencyRateIsZero):
		return status.Error(codes.Aborted, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
