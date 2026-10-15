package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ritikmayank/catch-it/config"
	"github.com/ritikmayank/catch-it/db"
	"github.com/ritikmayank/catch-it/handler"
	"github.com/ritikmayank/catch-it/repository"
	"github.com/ritikmayank/catch-it/seed"
	"github.com/ritikmayank/catch-it/web"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal().Msg("DATABASE_URL is required")
	}

	pool, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer pool.Close()

	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Fatal().Err(err).Msg("migrations failed")
	}

	clipRepo    := repository.NewClipRepository(pool)
	ratingRepo  := repository.NewRatingRepository(pool)
	streakRepo  := repository.NewStreakRepository(pool)
	badgeRepo   := repository.NewBadgeRepository(pool)

	seed.RunClips(context.Background(), clipRepo, time.Now().UTC().Truncate(24*time.Hour))

	h := handler.New(clipRepo, ratingRepo, streakRepo, badgeRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health",          handler.Health)
	mux.HandleFunc("GET /api/today",            cors(h.Today))
	mux.HandleFunc("POST /api/clips/{id}/rate", cors(h.Rate))
	mux.HandleFunc("GET /api/leaderboard",      cors(h.Leaderboard))
	mux.HandleFunc("GET /api/me",               cors(h.Me))
	mux.Handle("/", web.Handler())

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Port).Msg("catch-it server started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	log.Info().Msg("server stopped")
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Device-ID")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}
