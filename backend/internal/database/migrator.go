package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"strconv"

	"github.com/C0deNe0/agromart/internal/config"
	"github.com/jackc/pgx/v5"

	tern "github.com/jackc/tern/v2/migrate"
	"github.com/rs/zerolog"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(ctx context.Context, logger *zerolog.Logger, cfg *config.Config) error {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

	encodedPassword := url.QueryEscape(cfg.Database.Password)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		encodedPassword,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database for migration: %w", err)
	}
	defer conn.Close(ctx)

	m, err := tern.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	subtree, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations subdirectory: %w", err)
	}

	if err := m.LoadMigrations(subtree); err != nil {
		return fmt.Errorf("failed to load migrations from embed FS: %w", err)
	}

	from, err := m.GetCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	if err := m.Migrate(ctx); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if from == int32(len(m.Migrations)) {
		logger.Info().Msg("database is already up to date, no migrations applied")
	} else {
		logger.Info().Msgf("database migrated from version %d to %d", from, len(m.Migrations))
	}
	return nil
}
