package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/C0deNe0/agromart/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

const (
	DatabasePingTimeout = 10 * time.Second
)

type Database struct {
	Pool *pgxpool.Pool
	log  *zerolog.Logger
}

func New(cfg *config.Config, log *zerolog.Logger) (*Database, error) {
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

	//parse the pgx  config
	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the pgx pool config: %w", err)

	}

	

	//creating the pool
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	db := &Database{
		Pool: pool,
		log:  log,
	}

	//pinging the database to ensure connection is established
	ctx, cancel := context.WithTimeout(context.Background(), DatabasePingTimeout)
	defer cancel()

	if err := db.Pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func (db *Database) Close() error {
	db.log.Info().Msg("closing database connection pool")
	db.Pool.Close()
	return nil
}
