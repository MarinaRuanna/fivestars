package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"

	"fivestars/internal/infra/auth"
	"fivestars/internal/infra/config"
	"fivestars/internal/infra/controller"
	"fivestars/internal/infra/repository"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required for auth")
	}

	ctx := context.Background()
	pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	// Run migrations from embedded FS
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Printf("migrate (non-fatal): %v", err)
		// Continue anyway so app can start if DB is already migrated
	}

	establishmentRepo := repository.NewEstablishmentRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	healthHandler := controller.NewHealthHandler(pool)
	establishmentsHandler := controller.NewEstablishmentsHandler(establishmentRepo)
	authHandler := controller.NewAuthHandler(userRepo, cfg.JWTSecret)
	userHandler := controller.NewUserHandler(userRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		status, body := healthHandler.Ping(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write(body)
	})
	mux.HandleFunc("/establishments", func(w http.ResponseWriter, r *http.Request) {
		establishmentsHandler.List(w, r)
	})
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.Handle("GET /users/me", auth.RequireAuth(cfg.JWTSecret)(http.HandlerFunc(userHandler.Me)))

	withCORS := controller.CORS(mux)
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, withCORS); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func runMigrations(databaseURL string) error {
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("iofs: %w", err)
	}
	// driver postgres expects "postgres://..." or "postgresql://..."
	m, err := migrate.NewWithSourceInstance("iofs", source, databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
