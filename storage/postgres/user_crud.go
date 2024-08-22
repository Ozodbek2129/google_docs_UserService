package postgres

import (
	"database/sql"
	"google_docs_user/pkg/logger"
	"google_docs_user/storage"
	"log/slog"
)

type UserRepository struct {
	Db  *sql.DB
	Log *slog.Logger
}

func NewUserRepository(db *sql.DB) storage.IUserStorage {
	return &UserRepository{Db: db, Log: logger.NewLogger()}
}

func 