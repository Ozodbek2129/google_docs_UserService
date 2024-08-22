package service

import (
	"database/sql"
	"google_docs_user/storage"
	"google_docs_user/storage/postgres"
	"log/slog"
	pb "google_docs_user/genproto/user"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	User   storage.IStorage
	Logger *slog.Logger
}

func NewUserService(db *sql.DB, Logger *slog.Logger) *UserService {
	return &UserService{
		User:   postgres.NewPostgresStorage(db, Logger),
		Logger: Logger,
	}
}
