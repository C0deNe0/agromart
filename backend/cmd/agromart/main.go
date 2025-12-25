package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/C0deNe0/agromart/internal/config"
	"github.com/C0deNe0/agromart/internal/database"
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/C0deNe0/agromart/internal/logger"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/C0deNe0/agromart/internal/router"
	"github.com/C0deNe0/agromart/internal/server"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/go-playground/validator/v10"
)

const DefaultContextTimeout = 30

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		if err := validator.New().Struct(cfg); err != nil {
			panic("failed to validate the config: " + err.Error())
		}
		panic("failed to load the config: " + err.Error())
	}

	// googleOAuth := utils.NewGoogleOAuth(
	// 	cfg.OAuth.GoogleClientID,
	// 	cfg.OAuth.GoogleClientSecret,
	// 	cfg.OAuth.GoogleRedirectURI,
	// )

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

	repos := repository.NewRepositories(srv.DB.Pool)

	//INFRASTRUCTRE
	tokenManager := utils.NewTokenManager(cfg.Primary.Secret, cfg.Primary.Access)
	googleOAuth := utils.NewGoogleOAuth(
		cfg.OAuth.GoogleClientID,
		cfg.OAuth.GoogleClientSecret,
		cfg.OAuth.GoogleRedirectURI,
	)

	services := service.NewServices(repos, tokenManager, googleOAuth)
	handlers := handler.NewHandlers(services)
	r := router.NewRouter(&handlers, tokenManager)

	srv.SetupHTTPServer(r)

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
