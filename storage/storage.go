package storage

import (
	pb "google_docs_user/genproto/user"
	"context"
)

type IStorage interface {
	User() IUserStorage
	Close()
}

type IUserStorage interface {
	CreateUser(context.Context, *pb.RegisterReq) (*pb.RegisterRes, error)
	StoreRefreshToken(context.Context, *pb.StoreRefreshTokenReq) (*pb.StoreRefReshTokenRes, error)
	ConfirmationRegister(context.Context, *pb.ConfirmationReq) (*pb.ConfirmationRes,error)
	GetUserByEmail(context.Context, *pb.GetUSerByEmailReq) (*pb.GetUserResponse, error)
	UpdatePassword(context.Context, *pb.UpdatePasswordReq) (*pb.UpdatePasswordRes,error)
	ResetPassword(context.Context,*pb.ResetPasswordReq) (*pb.ResetPasswordRes,error)
	ConfirmationPassword(context.Context,*pb.ConfirmationReq)(*pb.ConfirmationResponse,error)
	UpdateUser(context.Context, *pb.UpdateUserRequest) (*pb.UpdateUserRespose,error)
	DeleteUser(context.Context, *pb.UserId) (*pb.DeleteUserr,error)
}