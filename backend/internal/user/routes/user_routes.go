package user

import (
	"portarius/internal/user/domain"
	userHandler "portarius/internal/user/handler"
	"portarius/internal/user/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	var (
		repo domain.IUserRepository = repository.NewUserRepository(db)
	)

	handler := userHandler.NewUserHandler(repo)

	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}
}

func RegisterUserProtectedRoutes(router *gin.RouterGroup, db *gorm.DB) {
	var (
		repo domain.IUserRepository = repository.NewUserRepository(db)
	)
	handler := userHandler.NewUserHandler(repo)

	users := router.Group("/users")
	{
		users.GET("/", handler.GetAll)
		users.GET("/:id", handler.GetByID)
		users.PUT("/:id", handler.Update)
		users.DELETE("/:id", handler.Delete)
	}
}
