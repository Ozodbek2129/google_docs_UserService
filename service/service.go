package service

import (
	"context"
	"database/sql"
	pb "google_docs_user/genproto/user"
	"google_docs_user/storage"
	"google_docs_user/storage/postgres"
	"google_docs_user/storage/redis"
	"log/slog"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	User storage.IStorage

	Logger *slog.Logger
}

func NewUserService(db *sql.DB, Logger *slog.Logger) *UserService {
	return &UserService{
		User:   postgres.NewPostgresStorage(db, Logger),
		Logger: Logger,
	}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	err := redis.RegisterUser(ctx, req)
	if err != nil {
		s.Logger.Error("failed to create user", "error", err)
		return nil, err
	}
	return &pb.RegisterRes{
		Message: "Sizning pochtangizga xabar yubordik.",
	}, nil
}

func (s *UserService) StoreRefreshToken(ctx context.Context, req *pb.StoreRefreshTokenReq) (*pb.StoreRefReshTokenRes, error) {
	res, err := s.User.User().StoreRefreshToken(ctx, req)
	if err != nil {
		s.Logger.Error("failed to login user", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) ConfirmationRegister(ctx context.Context, req *pb.ConfirmationRegisterReq) (*pb.ConfirmationRegisterRes, error) {
	res, err := s.User.User().ConfirmationRegister(ctx, req)
	if err != nil {
		s.Logger.Error("failed to confirmation register", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, req *pb.GetUSerByEmailReq) (*pb.GetUserResponse, error) {
	res, err := s.User.User().GetUserByEmail(ctx, req)
	if err != nil {
		s.Logger.Error("failed to get user by email", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordReq) (*pb.UpdatePasswordRes, error) {
	res, err := s.User.User().UpdatePassword(ctx, req)
	if err != nil {
		s.Logger.Error("failed to update password", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) ConfirmationPassword(ctx context.Context, req *pb.ConfirmationReq) (*pb.ConfirmationResponse, error) {
	res, err := s.User.User().ConfirmationPassword(ctx, req)
	if err != nil {
		s.Logger.Error("failed to confirmation password", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserRespose, error) {
	res, err := s.User.User().UpdateUser(ctx, req)
	if err != nil {
		s.Logger.Error("failed to update user", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.UserId) (*pb.DeleteUserr, error) {
	res, err := s.User.User().DeleteUser(ctx, req)
	if err != nil {
		s.Logger.Error("failed to delete user", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) UpdateRole(ctx context.Context,req *pb.UpdateRoleReq)(*pb.UpdateRoleRes,error){
	res, err := s.User.User().UpdateRole(ctx, req)
	if err != nil {
		s.Logger.Error("failed to UpdateRole user", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *UserService) ProfileImage(ctx context.Context,req *pb.ImageReq)(*pb.ImageRes,error){
	res, err := s.User.User().ProfileImage(ctx, req)
	if err != nil {
		s.Logger.Error("failed to ProfileImage user", "error", err)
		return nil, err
	}
	return res, nil
}