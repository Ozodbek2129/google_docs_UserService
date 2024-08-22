package handler


import (
	"google_docs_user/config"
	"google_docs_user/genproto/user"
	"google_docs_user/pkg/logger"
	"log"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	User user.UserServiceClient
	Log  *slog.Logger
}

func NewHandler() *Handler {

	conn, err := grpc.NewClient(config.Load().USER_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
			log.Println("error while connecting authentication service ", err)
	}

	return &Handler{
			User: user.NewUserServiceClient(conn),
			Log:  logger.NewLogger(),
	}
}