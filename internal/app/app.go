package app

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/artemiyKew/http-link-shortener/config"
	"github.com/artemiyKew/http-link-shortener/internal/delivery"
	"github.com/artemiyKew/http-link-shortener/internal/repo"
	"github.com/artemiyKew/http-link-shortener/internal/repo/pgdb"
	"github.com/artemiyKew/http-link-shortener/internal/repo/redisdb"
	"github.com/artemiyKew/http-link-shortener/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func Run(configPath string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Init config
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return err
	}

	// Logger
	SetLogrus("")
	if err != nil {
		return err
	}

	// Postgres
	logrus.Info("Initializing postgres...")
	db, err := pgdb.NewDB(cfg.DataBaseURL)
	if err != nil {
		return err
	}
	pg := pgdb.New(db)
	defer pg.DB.Close()

	// Redis
	logrus.Info("Initializing redis...")
	rdb := redisdb.NewRDB()
	if err := rdb.RDB.Ping(ctx).Err(); err != nil {
		return err
	}
	defer rdb.RDB.Close()

	// Repositories
	logrus.Info("Initializing repositories...")
	repositories := repo.NewRepositories(rdb, pg)

	// Services dependencies
	logrus.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos: repositories,
	}
	services := service.NewServices(deps)

	// Router
	logrus.Info("Initializing handlers and router...")
	handler := fiber.New()
	delivery.NewRouter(handler, services)

	// HTTP server
	logrus.Info("Starting http server...")
	logrus.Infof("Server port: %s", cfg.BindAddr)

	ch := make(chan error, 1)

	go func() {
		err = fasthttp.ListenAndServe(cfg.BindAddr, handler.Handler())
		if err != nil {
			ch <- err
		}
		close(ch)
	}()

	select {
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		return handler.ShutdownWithContext(timeout)
	case err := <-ch:
		return err
	}
}
