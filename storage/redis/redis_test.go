package redis_test

import (
	"context"
	"testing"

	pb "google_docs_user/genproto/user"

	"google_docs_user/storage/redis"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCode(t *testing.T) {
	ctx := context.Background()

	err := redis.StoreCode(ctx, "test@example.com", "123456")
	require.NoError(t, err)
}

func TestGetCode(t *testing.T) {
	ctx := context.Background()

	code, err := redis.GetCode(ctx, "test@example.com")
	require.NoError(t, err)
	assert.Equal(t, "123456", code)
}

func TestDeleteCode(t *testing.T) {
	ctx := context.Background()

	err := redis.DeleteCode(ctx, "test@example.com")
	require.NoError(t, err)
}

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()

	user := &pb.RegisterReq{
		Email:    "sanjarbeka775@gmail.com",
		FirstName: "salombek",
		LastName: "salombekov",
		Password: "password",
		Role: "admin",
		Code: 123456,
	}

	err := redis.RegisterUser(ctx, user)
	require.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()

	user := &pb.RegisterReq{
		Email:    "sanjarbeka775@gmail.com",
		FirstName: "salombek",
		LastName: "salombekov",
		Password: "password",
		Role: "admin",
		Code: 123456,
	}

	gotUser, err := redis.GetUser(ctx, "sanjarbeka775@gmail.com")
	require.NoError(t, err)
	assert.Equal(t, user, gotUser)
}
