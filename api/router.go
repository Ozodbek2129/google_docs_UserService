package api

import (
	"google_docs_user/api/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "google_docs_user/api/docs"
)

// @title        Google_docs_user API
// @version      1.0
// @description  This is an API for e-commerce platform.
// @termsOfService http://swagger.io/terms/
// @contact.name  API Support
// @contact.email support@swagger.io
// @BasePath      /
func NewRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.Use(handler.CORSMiddleware())
	
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	user := router.Group("/auth")
	{
		user.POST("/register",h.Register)
		user.POST("/login",h.LoginUser)
		user.GET("/confirm/:email/:code",h.ConfirmationRegister)
		user.GET("/reset_password/:email",h.ResetPassword)
		user.POST("/confirmation_password",h.ConfirmationPassword)
		user.PUT("/update_role/:email/:role",h.UpdateRole)
		user.POST("/products/media/:email",h.UploadMedia)
		user.PUT("/update_password/:email/:old_password/:new_password",h.UpdatePassword)
	}
	return router
}
