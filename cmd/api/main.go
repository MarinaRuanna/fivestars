package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fivestars/internal/infra"
)

func main() {
	// Setup context principal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ====== BUILD APP (Composition Root) ======
	log.Println("Building application...")
	app, err := infra.BuildApp(ctx)
	if err != nil {
		log.Fatalf("Failed to build app: %v", err)
	}

	// ====== SETUP GRACEFUL SHUTDOWN ======
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// ====== START SERVER IN GOROUTINE ======
	log.Printf("Starting server on port %d\n", app.Config.Port)
	errChan := make(chan error, 1)
	go func() {
		errChan <- app.Start(ctx)
	}()

	// ====== WAIT FOR SIGNAL OR ERROR ======
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v\n", sig)
	case err := <-errChan:
		if err != nil {
			log.Printf("Server error: %v\n", err)
		}
	}

	// ====== GRACEFUL SHUTDOWN ======
	log.Println("Shutting down gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.Stop(shutdownCtx); err != nil {
		log.Fatalf("Failed to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
