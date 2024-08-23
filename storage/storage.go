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
	StoreRefreshToken(context.Context, *pb.StoreRefreshTokenReq) (*pb.StoreRefReshTokenRes, error)
	ConfirmationRegister(context.Context, *pb.ConfirmationRegisterReq) (*pb.ConfirmationRegisterRes,error)
	GetUserByEmail(context.Context, *pb.GetUSerByEmailReq) (*pb.GetUserResponse, error)
	UpdatePassword(context.Context, *pb.UpdatePasswordReq) (*pb.UpdatePasswordRes,error)
	ConfirmationPassword(context.Context,*pb.ConfirmationReq)(*pb.ConfirmationResponse,error)
	UpdateUser(context.Context, *pb.UpdateUserRequest) (*pb.UpdateUserRespose,error)
	DeleteUser(context.Context, *pb.UserId) (*pb.DeleteUserr,error)
	UpdateRole(context.Context, *pb.UpdateRoleReq) (*pb.UpdateRoleRes,error)
	ProfileImage(ctx context.Context,req *pb.ImageReq)(*pb.ImageRes,error)
}