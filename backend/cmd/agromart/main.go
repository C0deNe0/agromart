package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/C0deNe0/agromart/internal/config"
	"github.com/C0deNe0/agromart/internal/database"
	"github.com/C0deNe0/agromart/internal/logger"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/C0deNe0/agromart/internal/server"
	"github.com/C0deNe0/agromart/internal/service"
)

const DefaultContextTimeout = 30

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load the config: " + err.Error())
	}

	//new custom logger
	log := logger.New(cfg.Primary.Env)

	//database migrations part
	if cfg.Primary.Env != "local" {
		if err := database.Migrate(context.Background(), log, cfg); err != nil {
			panic("failed to migrate the database: " + err.Error())
		}
	}

	//starting the server
	srv, err := server.New(cfg, log)
	if err != nil {
		panic("failed to create the server: " + err.Error())
	}

	repos := repository.NewRepository(srv.DB.Pool)
	services, err := service.NewService(srv, repos)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create services")
	}

	handlers := handler.NewHandler(srv, services)

	r := router.NewRouter(srv, handlers, services)

	srv.SetupHTTPServer(handlers)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// start server
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("could not start the server")
		}
	}()

	// wait for signal
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)
	defer cancel()
	// stop server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	// stop signal
	stop()

	log.Info().Msg("server stopped properly")
}
