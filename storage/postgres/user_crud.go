package postgres

import (
	"database/sql"
	"google_docs_user/genproto/user"
	"google_docs_user/models"
	"google_docs_user/pkg/logger"
	"google_docs_user/storage"
	"google_docs_user/storage/redis"
	"log/slog"

	"golang.org/x/net/context"
)

type UserRepository struct {
	Db  *sql.DB
	Log *slog.Logger
}

func NewUserRepository(db *sql.DB) storage.IUserStorage {
	return &UserRepository{Db: db, Log: logger.NewLogger()}
}

func (u *UserRepository) ConfirmationRegister(ctx context.Context, req *user.C) (*user.RegisterRes, error){

	res, err := redis.GetUser(ctx, req.Email)
	if err == nil {
        u.Log.Error("User already exists with email: %s", req.Email)
        return nil, err
    }

	req1 := models.Register{
		Email:     req.Email,
        Password:  req.Password,
        FirstName: req.FirstName,
        LastName:  req.LastName,
		Code: res.Code,
	}


}