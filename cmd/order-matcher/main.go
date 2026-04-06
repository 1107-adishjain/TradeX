package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ordermatcherapp "github.com/adishjain1107/tradex/pkg/order-matcher/app"
	"github.com/adishjain1107/tradex/pkg/order-matcher/config"
	api "github.com/adishjain1107/tradex/pkg/order-matcher/routes"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading env variables:", err)
	} else {
		log.Printf("Env variables loaded Successfully.")
	}

	application := ordermatcherapp.New(cfg)

	router := api.Routes(application)

	srv := &http.Server{
		Addr:         ":" + cfg.OrderMatcherPort,
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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ShutdownErr <- srv.Shutdown(ctx)
	}()

	log.Printf("Starting server on port %s", cfg.OrderMatcherPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}

	shutdownerr := <-ShutdownErr
	if shutdownerr != nil {
		log.Fatalf("Error shutting down server: %v", shutdownerr)
	}
	log.Println("Server gracefully stopped.")
}
