package main

import (
	"brunobevilaquaa/otel-distributed-tracing/internal/orchestration/adapters/api"
	"brunobevilaquaa/otel-distributed-tracing/internal/orchestration/adapters/handlers"
	"brunobevilaquaa/otel-distributed-tracing/internal/orchestration/services"
	"brunobevilaquaa/otel-distributed-tracing/pkg"
	"context"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	collectorUrl := os.Getenv("COLLECTOR_URL")
	if collectorUrl == "" {
		log.Fatal("COLLECTOR_URL is required")
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := pkg.InitProvider("orchestration", collectorUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("microservice-tracer")

	router := mux.NewRouter()

	apiClient := api.NewClient(tracer)

	weatherService := services.NewWeatherService(apiClient)

	weatherHandler := handlers.NewWeatherHandler(weatherService, tracer)

	router.HandleFunc("/api/v1/weather-check/{zipcode}", weatherHandler.CheckWeather).Methods(http.MethodGet)

	go func() {
		log.Println("Server starting on port 8081...")

		if err := http.ListenAndServe(":8081", router); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}
