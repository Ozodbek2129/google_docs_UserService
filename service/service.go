package service

import (
	"context"
	"database/sql"
	pb "google_docs_user/genproto/user"
	"google_docs_user/storage"
	"google_docs_user/storage/postgres"
	"log/slog"
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

func (s *UserService) CreateUser(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	res, err := s.User.User().CreateUser(ctx, req)
	if err != nil {
		s.Logger.Error("failed to create user", err)
		return nil, err
	}
	return res, nil	
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	res, err := s.User.User().Login(ctx, req)
	if err != nil {
		s.Logger.Error("failed to login user", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) ConfirmationRegister(ctx context.Context, req *pb.ConfirmationReq) (*pb.ConfirmationRes, error) {
	res, err := s.User.User().ConfirmationRegister(ctx, req)
	if err != nil {
		s.Logger.Error("failed to confirmation register", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, req *pb.GetUSerByEmailReq) (*pb.GetUserResponse, error) {
	res, err := s.User.User().GetUserByEmail(ctx, req)
	if err != nil {
		s.Logger.Error("failed to get user by email", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordReq) (*pb.UpdatePasswordRes, error) {
	res, err := s.User.User().UpdatePassword(ctx, req)
	if err != nil {
		s.Logger.Error("failed to update password", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *pb.ResetPasswordReq) (*pb.ResetPasswordRes, error) {
	res, err := s.User.User().ResetPassword(ctx, req)
	if err != nil {
		s.Logger.Error("failed to reset password", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) ConfirmationPassword(ctx context.Context, req *pb.ConfirmationReset) (*pb.ConfirmationResponse, error) {
	res, err := s.User.User().ConfirmationPassword(ctx, req)
	if err != nil {
		s.Logger.Error("failed to confirmation password", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserRespose, error) {
	res, err := s.User.User().UpdateUser(ctx, req)
	if err != nil {
		s.Logger.Error("failed to update user", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.UserId) (*pb.DeleteUserr, error) {
	res, err := s.User.User().DeleteUser(ctx, req)
	if err != nil {
		s.Logger.Error("failed to delete user", err)
		return nil, err
	}
	return res, nil
}