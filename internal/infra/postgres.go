package infra

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Heiji57/ETB-BE/config"
)

func NewPostgres(cfg *config.Config) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode,
	)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to postgres: %v", err))
	}
	if err := pool.Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("failed to ping postgres: %v", err))
	}
	return pool
}
