package infra

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"fivestars/internal/application/usecases"
	"fivestars/internal/infra/adapters/inbound"
	"fivestars/internal/infra/adapters/inbound/controller"
	"fivestars/internal/infra/adapters/outbound/repository/postgres"
	"fivestars/internal/infra/adapters/outbound/repository/postgres/checkins"
	"fivestars/internal/infra/adapters/outbound/repository/postgres/establishments"
	"fivestars/internal/infra/adapters/outbound/repository/postgres/highlights"
	"fivestars/internal/infra/adapters/outbound/repository/postgres/reviewlikes"
	"fivestars/internal/infra/adapters/outbound/repository/postgres/reviews"
	"fivestars/internal/infra/adapters/outbound/repository/postgres/users"
	"fivestars/internal/infra/config"
	"fivestars/internal/infra/policy"

	"github.com/jackc/pgx/v5/pgxpool"
)

// App encapsula toda a aplicação e seu ciclo de vida
type App struct {
	Config *config.Config
	DB     *pgxpool.Pool
	Router http.Handler
	server *http.Server
}

// BuildApp: COMPOSITION ROOT — monta TODAS as dependências
// Ordem crítica:
// 1. Configuração
// 2. Database
// 3. Migrations
// 4. Repositórios (outbound adapters)
// 5. Handlers (inbound adapters)
// 6. Rotas
func BuildApp(ctx context.Context) (*App, error) {
	// ====== 1. LOAD CONFIG ======
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// ====== 2. DATABASE POOL ======
	pool, err := postgres.NewPool(ctx, cfg.DatabasePostgres)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	// ====== 3. REPOSITORIES ======
	userRepo := users.NewUserRepository(pool)
	estabRepo := establishments.NewEstablishmentRepository(pool)
	checkinRepo := checkins.NewCheckinRepository(pool)
	reviewRepo := reviews.NewReviewRepository(pool)
	reviewLikeRepo := reviewlikes.NewReviewLikeRepository(pool)
	highlightRepo := highlights.NewHighlightRepository(pool)

	// ====== 4. USECASES ======
	registerUserUC := usecases.NewRegisterUserUseCase(userRepo, cfg.JWTSecret)
	loginUserUC := usecases.NewLoginUserUseCase(userRepo, cfg.JWTSecret)
	getUserUC := usecases.NewGetUserUseCase(userRepo)
	createEstablishmentUC := usecases.NewCreateEstablishmentUseCase(estabRepo)
	listEstabUC := usecases.NewListEstablishmentsUseCase(estabRepo)
	createCheckinUC := usecases.NewCreateCheckinUseCase(checkinRepo, estabRepo, 100.0)
	listCheckinsUC := usecases.NewListCheckinsUseCase(checkinRepo)
	getEstablishmentDetailUC := usecases.NewGetEstablishmentDetailUseCase(estabRepo, highlightRepo, reviewRepo)
	createReviewUC := usecases.NewCreateReviewUseCase(reviewRepo, checkinRepo)
	getReviewUC := usecases.NewGetReviewUseCase(reviewRepo)
	listReviewsUC := usecases.NewListReviewsByEstablishmentUseCase(reviewRepo)
	likeReviewUC := usecases.NewLikeReviewUseCase(reviewRepo, reviewLikeRepo)
	unlikeReviewUC := usecases.NewUnlikeReviewUseCase(reviewLikeRepo)
	operatorPolicy := policy.NewEstablishmentOwnerOperator(estabRepo)
	createHighlightUC := usecases.NewCreateHighlightUseCase(highlightRepo, reviewRepo, operatorPolicy)
	deleteHighlightUC := usecases.NewDeleteHighlightUseCase(highlightRepo, operatorPolicy)
	getEstablishmentStatsUC := usecases.NewGetEstablishmentStatsUseCase(estabRepo)

	// ====== 5. HANDLERS ======
	healthHandler := controller.NewHealthHandler(pool)
	authHandler := controller.NewAuthHandler(registerUserUC, loginUserUC)
	userHandler := controller.NewUserHandler(getUserUC)
	estabHandler := controller.NewEstablishmentsHandler(createEstablishmentUC, getEstablishmentDetailUC, listEstabUC, getEstablishmentStatsUC, createHighlightUC, deleteHighlightUC)
	checkinsHandler := controller.NewCheckinsHandler(createCheckinUC, listCheckinsUC)
	reviewsHandler := controller.NewReviewsHandler(createReviewUC, getReviewUC, listReviewsUC, likeReviewUC, unlikeReviewUC)

	// ====== 6. ROUTES ======
	controllers := inbound.Handlers{
		Health:         healthHandler,
		Auth:           authHandler,
		User:           userHandler,
		Establishments: estabHandler,
		Checkins:       checkinsHandler,
		Reviews:        reviewsHandler,
	}

	router := inbound.CreateChiRoutes(controllers, cfg.JWTSecret.Secret, cfg.CORS.AllowedOrigins)

	// ====== RETURN APP ======
	return &App{
		Config: cfg,
		DB:     pool,
		Router: router,
	}, nil
}

// Start: Inicia servidor HTTP com suporte a graceful shutdown via context
func (a *App) Start(ctx context.Context) error {
	a.server = &http.Server{
		Addr:         ":" + strconv.Itoa(a.Config.AppPort),
		Handler:      a.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Canal para erros do servidor
	errChan := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Aguardar: erro ou context cancelado
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return a.server.Shutdown(ctx)
	}
}

// Stop: Encerra gracefully (fecha DB, servidor, etc)
func (a *App) Stop(ctx context.Context) error {
	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			return err
		}
	}
	if a.DB != nil {
		a.DB.Close()
	}
	return nil
}
