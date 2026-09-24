package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hello-service/server"
)

func main() {
	// CLI logic placeholder — can be expanded for configuration parsing, etc.
	log.Println("Starting hello-service...")

	if err := setupServer(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func setupServer() error {
	srv := server.NewServer(":8080")
	if err := srv.RegisterDefaultRoutes(); err != nil {
		return err
	}

	// Ожидание сигнала для graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Println("Server listening on :8080")
	return srv.Start()
}
