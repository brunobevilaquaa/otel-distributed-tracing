package handlers

import (
	"brunobevilaquaa/otel-distributed-tracing/internal/orchestration/services"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type CheckWeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type WeatherHandler struct {
	service services.WeatherServiceInterface
	tracer  trace.Tracer
}

func NewWeatherHandler(service services.WeatherServiceInterface, tracer trace.Tracer) *WeatherHandler {
	return &WeatherHandler{
		service: service,
		tracer:  tracer,
	}
}

func (wh *WeatherHandler) CheckWeather(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := wh.tracer.Start(ctx, "orchestration process")
	defer span.End()

	zipcode := mux.Vars(r)["zipcode"]

	weather, err := wh.service.CheckWeather(ctx, zipcode)

	if err != nil {
		if errors.Is(err, services.ERROR_CANNOT_FIND_ZIPCODE) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res := CheckWeatherResponse{
		City:  weather.City,
		TempC: weather.TempC,
		TempF: weather.TempF,
		TempK: weather.TempK,
	}

	json.NewEncoder(w).Encode(res)
}
