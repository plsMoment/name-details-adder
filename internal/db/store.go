package db

import (
	"context"
	"fmt"
	"name-details-adder/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	connPool *pgxpool.Pool
}

func New(cfg *config.Config) (*Storage, error) {
	scope := "internal.db.store.New"
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", scope, err)
	}
	return &Storage{connPool: pool}, nil
}

func (s *Storage) Close() {
	s.connPool.Close()
}
