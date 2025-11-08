package server

import (
	"context"
	"fmt"
	"time"

	"regulation/internal/config"
	"regulation/internal/ent"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/DeltaLaboratory/entcache"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/rueidis"
)

func connectDB(ctx context.Context, config *config.Config, redis rueidis.Client) (*ent.Client, error) {
	connectionConfig, err := pgxpool.ParseConfig(config.DB.URI())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URI: %w", err)
	}

	connectionConfig.MaxConns = 75

	pool, err := pgxpool.NewWithConfig(ctx, connectionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	client := stdlib.OpenDBFromPool(pool)

	if err = client.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	entCacheDriver := entcache.NewDriver(
		entsql.OpenDB("postgres", client),
		entcache.TTL(time.Second*30),
		entcache.Levels(
			entcache.NewRedis(redis),
		),
	)

	entClient := ent.NewClient(ent.Driver(entCacheDriver))

	if config.Debug {
		return entClient.Debug(), nil
	}

	return entClient, nil
}
func (s *Server) setupDatabase(ctx context.Context) error {
	var err error

	// Connect to database
	s.db, err = connectDB(ctx, s.config, s.cache)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run database migrations using Ent's auto-migration
	if err := s.db.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
