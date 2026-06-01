package repository

import (
	"cart-api/internal/config"
	"fmt"
	"embed"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

var embedMigrations embed.FS

func InitDB(cfg config.Config) (*sqlx.DB, error){
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Mode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sqlx.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db.DB, "./db/migrations"); err != nil {
		return err
	}

	return nil
}