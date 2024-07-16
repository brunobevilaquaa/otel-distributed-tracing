package main

import (
	"brunobevilaquaa/otel-distributed-tracing/internal/input/handlers"
	"brunobevilaquaa/otel-distributed-tracing/internal/input/services"
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
	orchestrationServiceUrl := os.Getenv("ORCHESTRATION_SERVICE_URL")
	if orchestrationServiceUrl == "" {
		log.Fatal("ORCHESTRATION_SERVICE_URL is required")
	}

	collectorUrl := os.Getenv("COLLECTOR_URL")
	if collectorUrl == "" {
		log.Fatal("COLLECTOR_URL is required")
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := pkg.InitProvider("input", collectorUrl)
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

	zipCodeService := services.NewZipcodeService()

	weatherHandler := handlers.NewWeatherHandler(zipCodeService, orchestrationServiceUrl, tracer)

	router.HandleFunc("/api/v1/weather-check", weatherHandler.CheckWeather).Methods(http.MethodPost)

	go func() {
		log.Println("Server starting on port 8080...")

		if err := http.ListenAndServe(":8080", router); err != nil {
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
