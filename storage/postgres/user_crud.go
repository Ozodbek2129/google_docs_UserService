package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	pb "google_docs_user/genproto/user"
	"google_docs_user/pkg/logger"
	"google_docs_user/storage"
	"google_docs_user/storage/redis"
	"log/slog"
	"strings"
	"time"

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

func (u *UserRepository) StoreRefreshToken(ctx context.Context, req *pb.StoreRefreshTokenReq) (*pb.StoreRefReshTokenRes, error) {

	query := `INSERT INTO refresh_token (user_id, token) VALUES ($1, $2) returning id`

	var id string
	err := u.Db.QueryRow(query, req.UserId, req.Refresh).Scan(&id)
	if err != nil {
		u.Log.Error("Error storing refresh token", "error", err)
		return nil, errors.ErrUnsupported
	}

	return &pb.StoreRefReshTokenRes{
		Message: req.UserId,
	}, nil
}

func (u *UserRepository) ConfirmationRegister(ctx context.Context, req *pb.ConfirmationRegisterReq) (*pb.ConfirmationRegisterRes, error) {

	res, err := redis.GetUser(ctx, req.Email)
	if err != nil {
		u.Log.Error("User already exists", "email", req.Email)
		return nil, errors.New("user already exists")
	}

	if res.Code != req.Code {
		u.Log.Error("Invalid confirmation code", "email", req.Email)
		return nil, errors.New("invalid confirmation code")
	}

	id := uuid.NewString()
	var createdAt, updatedAt string

	query := `INSERT INTO users (id, email, password, first_name, last_name) VALUES ($1, $2, $3, $4, $5) returning created_at, updated_at`
	err = u.Db.QueryRow(query, id, req.Email, res.Password, res.FirstName, res.LastName).Scan(&createdAt, &updatedAt)
	if err != nil {
		u.Log.Error("Error inserting user", "err", err)
		return nil, err
	}

	return &pb.ConfirmationRegisterRes{
		User: &pb.User{
			Id:        id,
			Email:     req.Email,
			FirstName: res.FirstName,
			LastName:  res.LastName,
			Password:  res.Password,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
	}, nil
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, req *pb.GetUSerByEmailReq) (*pb.GetUserResponse, error) {
	var user pb.User
	var createdAt, updatedAt string

	query := `SELECT id, email, first_name, last_name, password, role, created_at, updated_at FROM users WHERE email = $1 and deleted_at = 0`

	err := u.Db.QueryRow(query, req.Email).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.Role, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		u.Log.Error("No user found", "email", req.Email)
		return nil, errors.New("no user found")
	} else if err != nil {
		u.Log.Error("Error getting user by email", "err", err)
		return nil, err
	}

	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return &pb.GetUserResponse{
		User: &user,
	}, nil
}

func (u *UserRepository) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordReq) (*pb.UpdatePasswordRes, error) {

	query := `UPDATE users SET password = $1 WHERE email = $2 AND deleted_at = 0`

	result, err := u.Db.ExecContext(ctx, query, req.NewPassword, req.Email)
	if err != nil {
		u.Log.Error("Error updating password", "err", err)
		return nil, errors.New("error updating password")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		u.Log.Error("Error retrieving rows affected", "err", err)
		return nil, errors.ErrUnsupported
	}

	if rowsAffected == 0 {
		u.Log.Error("No user found", "email", req.Email)
		return nil, errors.New("no user found")
	}

	return &pb.UpdatePasswordRes{
		Message: "Your password has been changed.",
	}, nil
}

func (u *UserRepository) ConfirmationPassword(ctx context.Context, req *pb.ConfirmationReq) (*pb.ConfirmationResponse, error) {

	query := `UPDATE users SET password = $1 WHERE email = $2 AND deleted_at = 0`
	_,err := u.Db.ExecContext(ctx,query, req.NewPassword, req.Email)
	if err != nil {
		u.Log.Error("No user found with email ", "email", req.Email)
		return nil, errors.New("no user found")
	}

	return &pb.ConfirmationResponse{
		Message: "Your Password changed.",
	}, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserRespose, error) {
	var existingUser struct {
		Email     string
		FirstName string
		LastName  string
	}

	query := `SELECT email, first_name, last_name FROM users WHERE id = $1 AND deleted_at = 0`
	err := u.Db.QueryRowContext(ctx, query, req.Id).Scan(&existingUser.Email, &existingUser.FirstName, &existingUser.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %s not found", req.Id)
		}
		return nil, fmt.Errorf("failed to fetch existing user: %v", err)
	}

	if req.Email == "" || !strings.Contains(req.Email, "@") || !strings.HasSuffix(req.Email, "@gmail.com") {
		req.Email = existingUser.Email
	}

	if req.FirstName == "" {
		req.FirstName = existingUser.FirstName
	}

	if req.LastName == "" {
		req.LastName = existingUser.LastName
	}

	updateQuery := `UPDATE users SET email = $1, first_name = $2, last_name = $3, updated_at = $4 WHERE id = $5 AND deleted_at = 0`
	_, err = u.Db.ExecContext(ctx, updateQuery, req.Email, req.FirstName, req.LastName, time.Now(), req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	response := &pb.UpdateUserRespose{
		Message: "User updated successfully",
	}

	return response, nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, req *pb.UserId) (*pb.DeleteUserr, error) {
	query := `UPDATE users SET deleted_at = $1 WHERE id = $2 AND deleted_at = 0`
	res, err := u.Db.ExecContext(ctx, query, time.Now(), req.Id)
	if err != nil {
		u.Log.Error("Error deleting ", "user", err)
		return nil, errors.ErrUnsupported
	}

	count, err := res.RowsAffected()
	if err != nil {
		u.Log.Error("Error getting rows affected ", "err", err)
		return nil, errors.ErrUnsupported
	}

	if count == 0 {
		u.Log.Error("No user found ", "ID", req.Id)
		return nil, errors.New("no user found")
	}

	return &pb.DeleteUserr{
		Message: "User deleted successfully",
	}, nil
}

func (u *UserRepository) UpdateRole(ctx context.Context, req *pb.UpdateRoleReq) (*pb.UpdateRoleRes, error) {
	query := `UPDATE users SET role = $1 WHERE email = $2 AND deleted_at = 0`
	res, err := u.Db.ExecContext(ctx, query, req.Role, req.Email)
	if err != nil {
		u.Log.Error("Error updating role ", "err", err)
		return nil, errors.ErrUnsupported
	}

	count, err := res.RowsAffected()
	if err != nil {
		u.Log.Error("Error getting rows affected ", "err", err)
		return nil, errors.ErrUnsupported
	}

	if count == 0 {
		u.Log.Error("No user not found ", "email", req.Email)
		return nil, errors.New("no user found")
	}

	return &pb.UpdateRoleRes{
		Message: "Role updated successfully",
	}, nil
}

func (u *UserRepository) ProfileImage(ctx context.Context, req *pb.ImageReq) (*pb.ImageRes, error) {
	query := `
		UPDATE users
		SET image = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE email = $2 AND deleted_at = 0`

	_, err := u.Db.ExecContext(ctx, query, req.Image, req.Email)
	if err != nil {
		u.Log.Error("No user not found ", "email", req.Email)
		return nil, errors.New("no user found")
	}
	fmt.Println("Image updated successfully")
	return &pb.ImageRes{
		Message: "Image uploaded successfully",
	}, nil
}
