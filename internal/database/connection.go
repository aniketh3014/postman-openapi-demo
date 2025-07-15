package database

import (
	"context"
	"database/sql"
	"fmt"
	"postman-api/internal/config"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Database struct {
	*bun.DB
}

func NewConnection(cfg *config.DatabaseConfig) (*Database, error) {
	sqldb, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection %w", err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.DB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database %w", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}
