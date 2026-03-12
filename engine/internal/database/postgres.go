package database

import (
	"context"
	"log"

	"github.com/Prayas-35/ragkit/engine/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
)

var DB *pgxpool.Pool

func Connect() {
	cfg := config.LoadConfig()
	dsn := cfg.DatabaseUri
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Failed to parse database config:", err)
	}
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxvec.RegisterTypes(ctx, conn)
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	DB = dbpool
	log.Println("✅ Connected to Neon PostgreSQL")

	if cfg.DB_SYNC {
		log.Println("🔄 Running database migrations...")
		if err := runMigrations(dsn); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
	}
}

func runMigrations(dsn string) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("✅ Migrations applied successfully")
	return nil
}
