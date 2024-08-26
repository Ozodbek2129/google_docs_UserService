package postgres

import (
	"context"
	"google_docs_user/config"
	pb "google_docs_user/genproto/user"
	"google_docs_user/pkg/logger"
	"log"
	"testing"
)

func TestConnectionDb(t *testing.T) {
	_ = config.Load()
	db, err := ConnectionDb()
	if err != nil {
		log.Fatal("error connection to db: ", err)
	}
	defer db.Close()
	t.Log("success")
}

func TestConfirmationRegister(t *testing.T) {
	_ = config.Load()

	db, err := ConnectionDb()
	if err != nil {
		log.Fatal("error connection to db: ", err)
	}
	defer db.Close()
	storage := NewPostgresStorage(db, logger.NewLogger())
	_, err = storage.User().ConfirmationRegister(context.Background(), &pb.ConfirmationRegisterReq{
			Email: "sanjarbeka775@gmail.com",
			Code:  123456,
	})
	if err != nil {
		t.Error("error while confirmation register: ", err)
	}
	t.Log("success")
}

func TestGetUserByEmail(t *testing.T) {
	_ = config.Load()
	db, err := ConnectionDb()
	if err != nil {
		log.Fatal("error connection to db: ", err)
	}
	defer db.Close()
	storage := NewPostgresStorage(db, logger.NewLogger())
	_, err = storage.User().GetUserByEmail(context.Background(), &pb.GetUSerByEmailReq{
		Email: "ozodbek2129@gmail.com",
	})
	if err != nil {
		t.Error("error while getting user by email: ", err)
	}
	t.Log("success")
}

func TestUpdatePassword(t *testing.T) {
	_ = config.Load()
	db, err := ConnectionDb()
	if err != nil {
		log.Fatal("error connection to db: ", err)
	}
	defer db.Close()
	storage := NewPostgresStorage(db, logger.NewLogger())
	_, err = storage.User().UpdatePassword(context.Background(), &pb.UpdatePasswordReq{

		Email:       "ozodbek2129@gmail.com",
		NewPassword: "123456",
	})
	if err != nil {
		t.Error("error while updating password: ", err)
	}
	t.Log("success")
}

func TestConfirmationPassword(t *testing.T) {
	_ = config.Load()
	db, err := ConnectionDb()
	if err != nil {
		log.Fatal("error connection to db: ", err)
	}
	defer db.Close()
	storage := NewPostgresStorage(db, logger.NewLogger())
	_, err = storage.User().ConfirmationPassword(context.Background(), &pb.ConfirmationReq{
		Email: "ozodbek2129@gmail.com",
		Code:  123456,
	})
	if err != nil {
		t.Error("error while confirmation password: ", err)
	}
	t.Log("success")
}

func TestUpdateUser(t *testing.T) {
	_ = config.Load()
	db, err := ConnectionDb()
	if err != nil {
		log.Fatal("error connection to db: ", err)
	}
	defer db.Close()
	storage := NewPostgresStorage(db, logger.NewLogger())
	_, err = storage.User().UpdateUser(context.Background(), &pb.UpdateUserRequest{
		Id:        "1",
		FirstName: "Ozodbek",
		LastName:  "Raximov",
		Email:     "ozodbek2129@gmail.com",
	})
	if err != nil {
		t.Error("error while updating user: ", err)
	}
	t.Log("success")
}

func TestDeleteUser(t *testing.T) {
	_ = config.Load()
	db, err := ConnectionDb()
	if err != nil {
		log.Fatal("error connection to db: ", err)
	}
	defer db.Close()
	storage := NewPostgresStorage(db, logger.NewLogger())
	_, err = storage.User().DeleteUser(context.Background(), &pb.UserId{
		Id: "1",
	})
	if err != nil {
		t.Error("error while deleting user: ", err)
	}
	t.Log("success")
}

func TestUpdateRole(t *testing.T) {
	_ = config.Load()
	db, err := ConnectionDb()
	if err != nil {
		log.Fatal("error connection to db: ", err)
	}
	defer db.Close()
	storage := NewPostgresStorage(db, logger.NewLogger())
	_, err = storage.User().UpdateRole(context.Background(), &pb.UpdateRoleReq{
		Email: "ozodbek2129@gmail.com",
		Role:  "user",
	})
	if err != nil {
		t.Error("error while updating role: ", err)
	}
	t.Log("success")
}

func TestStoreRefreshToken(t *testing.T) {
	_ = config.Load()
	db, err := ConnectionDb()
	if err != nil {
		log.Fatal("error connection to db: ", err)
	}
	defer db.Close()
	storage := NewPostgresStorage(db, logger.NewLogger())
	_, err = storage.User().StoreRefreshToken(context.Background(), &pb.StoreRefreshTokenReq{
		UserId:  "1",
		Refresh: "123456",
	})
	if err != nil {
		t.Error("error while storing refresh token: ", err)
	}
	t.Log("success")
}
