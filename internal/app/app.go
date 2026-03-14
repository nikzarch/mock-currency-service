package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikzarch/mock-currency-service/internal/config"
	"github.com/nikzarch/mock-currency-service/internal/db"
	"github.com/nikzarch/mock-currency-service/internal/health"
	"github.com/nikzarch/mock-currency-service/internal/valute"
)

func Run() error {
	cfg := config.MustLoad()

	pool, err := db.NewPool()
	if err != nil {
		return err
	}
	defer pool.Close()

	repo := valute.NewPostgresRepository(pool)
	gen := valute.NewGenerator()
	svc := valute.NewValuteService(repo, gen, cfg.ResponseMode)
	h := valute.NewHandler(svc)

	mux := http.NewServeMux()
	mux.Handle("/healthz", health.Handler())
	mux.Handle("/scripts/XML_daily.asp", h)

	server := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("server started on %s", cfg.HTTPPort)
		errCh <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stop:
		log.Printf("shutdown signal: %s", sig)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
