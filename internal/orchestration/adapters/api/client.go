package api

import (
	"brunobevilaquaa/otel-distributed-tracing/internal/orchestration/domain"
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	netUrl "net/url"
	"os"
	"time"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ClientInterface interface {
	GetWeatherByLocale(ctx context.Context, locale string) (*domain.WeatherResult, error)
	GetLocaleByZipcode(ctx context.Context, zipcode string) (*domain.LocaleResult, error)
}

type Client struct {
	tracer trace.Tracer
}

func NewClient(tracer trace.Tracer) *Client {
	return &Client{
		tracer: tracer,
	}
}

func (c *Client) GetWeatherByLocale(ctx context.Context, locale string) (*domain.WeatherResult, error) {
	ctx, span := c.tracer.Start(ctx, "get weather by locale")
	time.Sleep(time.Second)
	defer span.End()

	token := os.Getenv("WEATHER_API_KEY")

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", token, netUrl.QueryEscape(locale))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var weatherAPIResponse WeatherAPIResponse
	err = json.NewDecoder(res.Body).Decode(&weatherAPIResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &domain.WeatherResult{
		TempC: weatherAPIResponse.Current.TempC,
	}, nil
}

func (c *Client) GetLocaleByZipcode(ctx context.Context, zipcode string) (*domain.LocaleResult, error) {
	ctx, span := c.tracer.Start(ctx, "get locale by zipcode")
	time.Sleep(time.Second)
	defer span.End()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipcode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var viaCepResponse ViaCepResponse
	err = json.NewDecoder(res.Body).Decode(&viaCepResponse)
	if err != nil {
		return nil, err
	}

	return &domain.LocaleResult{
		Locale: viaCepResponse.Localidade,
	}, nil
}
