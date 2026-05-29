package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/Konshtantin/go-rest/internal/config"
	"github.com/Konshtantin/go-rest/internal/handler"
	"github.com/Konshtantin/go-rest/internal/repository"
	"github.com/Konshtantin/go-rest/internal/service"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DBConn)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer redisClient.Close()
	
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Warn("redis unavailable, running without cache", "error", err)
	} else {
		log.Info("connected to redis")
	}
	
	if err := pool.Ping(context.Background()); err != nil {
		log.Error("database ping failed", "error", err)
		os.Exit(1)
	}
	log.Info("connected to database")

	pgxUserRepo := repository.NewPostgresUserRepo(pool)
	userRepo := repository.NewCachedUserRepo(pgxUserRepo, redisClient)

	postRepo := repository.NewPostgresPostRepo(pool)

	userSvc := service.NewUserService(userRepo)
	postSvc := service.NewPostService(postRepo)

	userHandler := handler.NewUserHandler(userSvc, log)
	postHandler := handler.NewPostHandler(postSvc, log)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.List)
		r.Post("/", userHandler.Create)
		r.Get("/{id}", userHandler.GetByID)
		r.Delete("/{id}", userHandler.Delete)
	})


	r.Route("/posts", func(r chi.Router) {
		r.Get("/", postHandler.List)
		r.Post("/", postHandler.Create)
		r.Get("/{id}", postHandler.GetByID)
		r.Delete("/{id}", postHandler.Delete)
	})




	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	go func() {
		log.Info("server started", "port", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	<-quit

	log.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown error", "error", err)
	}
	
	log.Info("server stopped")
}