package postgres

import (
	"database/sql"
	"fmt"
	"google_docs_user/config"
	"google_docs_user/storage"
	"log/slog"

	_ "github.com/lib/pq"
)

type postgresStorage struct {
	db  *sql.DB
	log *slog.Logger
}

func NewPostgresStorage(db *sql.DB, log *slog.Logger) storage.IStorage {
	return &postgresStorage{
			db:  db,
			log: log,
	}
}

func ConnectionDb() (*sql.DB, error) {
	conf := config.Load()
	conDb := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			conf.DB_HOST, conf.DB_PORT, conf.DB_USER, conf.DB_NAME, conf.DB_PASSWORD)
	db, err := sql.Open("postgres", conDb)
	if err != nil {
			return nil, err
	}

	if err := db.Ping(); err != nil {
			return nil, err
	}

	return db, nil
}

func (p *postgresStorage) Close() {
	p.db.Close()
}

func (p *postgresStorage) User() storage.IUserStorage {
	return NewUserRepository(p.db)
}