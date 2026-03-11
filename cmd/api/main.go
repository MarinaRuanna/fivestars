package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"fivestars/internal/infra"
)

func main() {
	if err := loadDotEnvLocal(".env.local"); err != nil {
		log.Fatalf("Failed to load .env.local: %v", err)
	}

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
	log.Printf("Starting server on port %d\n", app.Config.AppPort)
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

// loadDotEnvLocal loads variables from a local .env file if it exists.
// It does not override variables already set in the environment.
func loadDotEnvLocal(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}
		if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) ||
			(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
			val = strings.Trim(val, "\"'")
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}

	return nil
}
