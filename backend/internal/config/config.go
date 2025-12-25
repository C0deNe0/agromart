package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary  Primary        `koanf:"primary" validate:"required"`
	Server   Server         `koanf:"server" validate:"required"`
	Database DatabaseConfig `koanf:"database" validate:"required"`
	OAuth    OAuthConfig    `koanf:"oauth" validate:"required"`
}

type Primary struct {
	Env    string `koanf:"env" validate:"required"`
	Secret string `koanf:"secret" validate:"required"`
	Access string `koanf:"access" validate:"required"`
}

type Server struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"dive,url"`
}

type DatabaseConfig struct {
	Host            string `koanf:"host" validate:"required"`
	Port            int    `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

type OAuthConfig struct {
	GoogleClientID     string `koanf:"google_client_id" validate:"required"`
	GoogleClientSecret string `koanf:"google_client_secret" validate:"required"`
	GoogleRedirectURI  string `koanf:"google_redirect_uri" validate:"required"`
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	k := koanf.New(".")
	if err := k.Load(env.Provider("AGROMART_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "AGROMART_"))
	}), nil); err != nil {
		logger.Fatal().Err(err).Msg("failed to load config from environment variables")
	}

	mainConfig := &Config{}

	if err := k.Unmarshal("", mainConfig); err != nil {
		logger.Fatal().Err(err).Msg("failed to unmarshal config into struct")
	}

	validate := validator.New()
	if err := validate.Struct(mainConfig); err != nil {
		logger.Fatal().Err(err).Msg("config validation failed")
	}

	return mainConfig, nil
}
