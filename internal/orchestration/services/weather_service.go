package services

import (
	"brunobevilaquaa/otel-distributed-tracing/internal/orchestration/adapters/api"
	"brunobevilaquaa/otel-distributed-tracing/internal/orchestration/domain"
	"context"
	"errors"
	"fmt"
)

var (
	ERROR_CANNOT_GET_LOCALE   = errors.New("error on get zipcode")
	ERROR_CANNOT_FIND_ZIPCODE = errors.New("can not find zipcode")
	ERROR_CANNOT_GET_WHEATER  = errors.New("error on get weather")
)

type WeatherServiceInterface interface {
	CheckWeather(ctx context.Context, zipcode string) (*domain.Result, error)
}

type WeatherService struct {
	Client api.ClientInterface
}

func NewWeatherService(client api.ClientInterface) *WeatherService {
	return &WeatherService{
		Client: client,
	}
}

func (w *WeatherService) CheckWeather(ctx context.Context, zipcode string) (*domain.Result, error) {
	locale, err := w.Client.GetLocaleByZipcode(ctx, zipcode)
	if err != nil {
		fmt.Println("err", err)
		return nil, ERROR_CANNOT_GET_LOCALE
	}

	if locale.Locale == "" {
		return nil, ERROR_CANNOT_FIND_ZIPCODE
	}

	weather, err := w.Client.GetWeatherByLocale(ctx, locale.Locale)
	if err != nil {
		return nil, ERROR_CANNOT_GET_WHEATER
	}

	return &domain.Result{
		TempC: weather.TempC,
		TempF: weather.TempC*1.8 + 32,
		TempK: weather.TempC + 273.15,
	}, nil
}
