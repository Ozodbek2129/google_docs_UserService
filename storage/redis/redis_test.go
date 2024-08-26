package redis_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	pb "google_docs_user/genproto/user"

	"google_docs_user/storage/redis"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCode(t *testing.T) {
	ctx := context.Background()
	_, mock := redismock.NewClientMock()

	mock.ExpectSet("test@example.com", "123456", time.Minute*3).SetVal("OK")

	err := redis.StoreCode(ctx, "test@example.com", "123456")
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetCode(t *testing.T) {
	ctx := context.Background()

	_, mock := redismock.NewClientMock()

	mock.ExpectGet("test@example.com").SetVal("123456")

	code, err := redis.GetCode(ctx, "test@example.com")
	require.NoError(t, err)
	assert.Equal(t, "123456", code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCode(t *testing.T) {
	ctx := context.Background()

	_, mock := redismock.NewClientMock()

	mock.ExpectDel("test@example.com").SetVal(1)

	err := redis.DeleteCode(ctx, "test@example.com")
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()

	_, mock := redismock.NewClientMock()

	user := &pb.RegisterReq{
		Email:    "test@example.com",
		Password: "password",
	}

	data, _ := json.Marshal(user)
	mock.ExpectSet("test@example.com", data, time.Minute*10).SetVal("OK")

	err := redis.RegisterUser(ctx, user)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()

	_, mock := redismock.NewClientMock()

	user := &pb.RegisterReq{
		Email:    "test@example.com",
		Password: "password",
	}

	data, _ := json.Marshal(user)
	mock.ExpectGet("test@example.com").SetVal(string(data))

	gotUser, err := redis.GetUser(ctx, "test@example.com")
	require.NoError(t, err)
	assert.Equal(t, user, gotUser)
	assert.NoError(t, mock.ExpectationsWereMet())
}
