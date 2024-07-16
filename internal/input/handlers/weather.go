package handlers

import (
	"brunobevilaquaa/otel-distributed-tracing/internal/input/services"
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"net/http"
)

type CheckWeatherRequest struct {
	Cep string `json:"cep"`
}

type WeatherHandler struct {
	service                 services.ZipcodeServiceInterface
	orchestrationServiceUrl string
	tracer                  trace.Tracer
}

func NewWeatherHandler(service services.ZipcodeServiceInterface, orchestrationServiceUrl string, tracer trace.Tracer) *WeatherHandler {
	return &WeatherHandler{
		service:                 service,
		orchestrationServiceUrl: orchestrationServiceUrl,
		tracer:                  tracer,
	}
}

func (wh *WeatherHandler) CheckWeather(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := wh.tracer.Start(ctx, "call to orchestration service")
	defer span.End()

	var data CheckWeatherRequest

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = wh.service.CheckZipcode(data.Cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	url := fmt.Sprintf("%s/api/v1/weather-check/%s", wh.orchestrationServiceUrl, data.Cep)

	log.Println("Requesting orchestration service at", url)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bodyBytes)
}
