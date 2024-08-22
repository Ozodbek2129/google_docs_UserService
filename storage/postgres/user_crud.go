package postgres

import (
	"database/sql"
	"errors"
	"google_docs_user/genproto/user"
	"google_docs_user/pkg/logger"
	"google_docs_user/storage"
	"google_docs_user/storage/redis"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type UserRepository struct {
	Db  *sql.DB
	Log *slog.Logger
}

func NewUserRepository(db *sql.DB) storage.IUserStorage {
	return &UserRepository{Db: db, Log: logger.NewLogger()}
}

func (u *UserRepository) ConfirmationRegister(ctx context.Context, req *user.ConfirmationReq) (*user.RegisterRes, error) {

	res, err := redis.GetUser(ctx, req.Email)
	if err == nil {
		u.Log.Error("User already exists with email: %s", req.Email)
		return nil, err
	}

	if res.Code != req.Code {
		u.Log.Error("Invalid confirmation code for email: %s", req.Email)
		return nil, errors.New("Invalid confirmation code")
	}

	id := uuid.NewString()

	query := `INSERT INTO users (id, email, password, first_name, last_name) VALUES ($1, $2, $3, $4, $5)`

	_, err = u.Db.ExecContext(ctx, query, id, res.Email, res.Password, res.FirstName, res.LastName)
	if err!= nil {
        u.Log.Error("Failed to insert user into database: %v", err)
        return nil, err
    }

	
}
