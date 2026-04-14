package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	marketdataapp "github.com/adishjain1107/tradex/pkg/market-data/app"
	"github.com/adishjain1107/tradex/pkg/market-data/binance"
	"github.com/adishjain1107/tradex/pkg/market-data/config"
	api "github.com/adishjain1107/tradex/pkg/market-data/routes"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading env variables:", err)
	} else {
		log.Printf("Env variables loaded Successfully.")
	}

	application := marketdataapp.New(cfg)

	router := api.Routes(application)
	streamCtx, cancelStream := context.WithCancel(context.Background())
	defer cancelStream()
	streamStopped := make(chan struct{})

	go func() {
		defer close(streamStopped)
		binance.StartMultiplexStream(streamCtx, cfg.KafkaBroker, cfg.Symbols)
	}()

	srv := &http.Server{
		Addr:         ":" + cfg.MarketDataPort,
		Handler:      router,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Minute,
		IdleTimeout:  60 * time.Minute,
	}

	ShutdownErr := make(chan error)

	go func() { //this go-routine listens for the shutdown signal and gracefully shuts down the server when it receives one.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		log.Printf("Received shutdown signal: %v", s)
		cancelStream()

		select {
		case <-streamStopped:
			log.Println("Binance stream stopped.")
		case <-time.After(2 * time.Second):
			log.Println("Binance stream stop not confirmed within 2s.")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ShutdownErr <- srv.Shutdown(ctx)
	}()

	log.Printf("Starting server on port %s", cfg.MarketDataPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}

	shutdownerr := <-ShutdownErr
	if shutdownerr != nil {
		log.Fatalf("Error shutting down server: %v", shutdownerr)
	}
	log.Println("Server gracefully stopped.")
}
