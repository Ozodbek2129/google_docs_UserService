package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"google_docs_user/config"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	pb "google_docs_user/genproto/user"
)

func ConnectDB() *redis.Client {
	cfg := config.Load()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RD_HOST,
		Password: cfg.RD_PASSWORD,
		DB:       cfg.RD_NAME,
	})

	return rdb
}

func StoreCode(ctx context.Context, email, code string) error {
	fmt.Println("Storing code for email: " + email + " with code: " + code)
	rdb := ConnectDB()

	err := rdb.Set(ctx, email, code, time.Minute*3).Err()
	if err != nil {
		return errors.Wrap(err, "failed to store code in redis")

	}

	return nil
}

func GetCode(ctx context.Context, email string) (string, error) {
	rdb := ConnectDB()

	code, err := rdb.Get(ctx, email).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("code not found for " + email)
		}
		return "", errors.Wrap(err, "failed to get code from redis")
	}

	return code, nil
}

func DeleteCode(ctx context.Context, email string) error {
	rdb := ConnectDB()

	err := rdb.Del(ctx, email).Err()
	if err != nil {
		return errors.Wrap(err, "failed to delete code from redis")
	}

	return nil
}

func RegisterUser(ctx context.Context, req *pb.RegisterReq) error {
	rdb := ConnectDB()

	res, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "failed to marshal user data")
	}

	err = rdb.Set(ctx, req.Email, res, time.Minute*10).Err()
	if err != nil {
		return errors.Wrap(err, "failed to store user data in redis")
	}

	return nil
}


func GetUser(ctx context.Context, email string) (*pb.RegisterReq, error) {
	rdb := ConnectDB()

    res, err := rdb.Get(ctx, email).Result()
    if err!= nil {
        if err == redis.Nil {
            return nil, errors.New("user not found for " + email)
        }
        return nil, errors.Wrap(err, "failed to get user from redis")
    }

    var user *pb.RegisterReq
    err = json.Unmarshal([]byte(res), &user)
    if err!= nil {
        return nil, errors.Wrap(err, "failed to unmarshal user data")
    }

    return user, nil
}