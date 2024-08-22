package token

import (
	"errors"
	"google_docs_user/config"
	pb "google_docs_user/genproto/user"
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

func GenerateJWT(user *pb.User) *pb.LoginRes {
        accesstoken := jwt.New(jwt.SigningMethodHS256)
        refreshToken := jwt.New(jwt.SigningMethodHS256)

        accesstClaim := accesstoken.Claims.(jwt.MapClaims)
        accesstClaim["user_id"] = user.Id
        accesstClaim["email"] = user.Email
        accesstClaim["password"] = user.Password
        accesstClaim["role"] = user.Role
        accesstClaim["iat"] = time.Now().Unix()
        accesstClaim["exp"] = time.Now().Add(time.Hour).Unix()

        con := config.Load()
        access, err := accesstoken.SignedString([]byte(con.SIGNING_KEY))
        if err != nil {
                log.Fatalf("Error with generating access token: %s", err)
        }

        refreshClaim := refreshToken.Claims.(jwt.MapClaims)
        refreshClaim["user_id"] = user.Id
        refreshClaim["email"] = user.Email
        refreshClaim["password"] = user.Password
        refreshClaim["role"] = user.Role
        refreshClaim["iat"] = time.Now().Unix()
        refreshClaim["exp"] = time.Now().Add(time.Hour).Unix()

        refresh, err := refreshToken.SignedString([]byte(con.SIGNING_KEY))
        if err != nil {
                log.Fatalf("Error with generating access token: %s", err)
        }

        return &pb.LoginRes{
                Access:  access,
                Refresh: refresh,
        }
}

func RefreshJWT(refreshTokenString string) (*pb.LoginRes, error) {
        con := config.Load()
        token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, errors.New("unexpected signing method")
                }
                return []byte(con.SIGNING_KEY), nil
        })

        if err != nil {
                return nil, err
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
                user := &pb.User{}

                if val, ok := claims["user_id"].(string); ok {
                        user.Id = val
                } else {
                        return nil, errors.New("user_id claim is missing or invalid")
                }

                if val, ok := claims["email"].(string); ok {
                        user.Email = val
                } else {
                        return nil, errors.New("username claim is missing or invalid")
                }

                if val, ok := claims["password"].(string); ok {
                        user.Password = val
                } else {
                        return nil, errors.New("email claim is missing or invalid")
                }

                if val, ok := claims["role"].(string); ok {
                        user.Role = val
                } else {
                        return nil, errors.New("full_name claim is missing or invalid")
                }

                return GenerateJWT(user), nil
        }

        return nil, errors.New("invalid refresh token")
}