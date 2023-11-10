package app

import (
	"context"

	"github.com/artemiyKew/http-link-shortener/config"
	"github.com/artemiyKew/http-link-shortener/internal/delivery"
	"github.com/artemiyKew/http-link-shortener/internal/repo"
	"github.com/artemiyKew/http-link-shortener/internal/repo/redisdb"
	"github.com/artemiyKew/http-link-shortener/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func Run(ctx context.Context, configPath string) error {
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

	// Redis
	logrus.Info("Initializing redis...")
	rdb := redisdb.NewRDB()
	if err := rdb.RDB.Ping(ctx).Err(); err != nil {
		return err
	}
	defer rdb.RDB.Close()

	// Repositories
	logrus.Info("Initializing repositories...")
	repositories := repo.NewRepositories(rdb)

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

	return fasthttp.ListenAndServe(cfg.BindAddr, handler.Handler())
}
