package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	handler := NewUserHandler(db)

	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}
}

func RegisterUserProtectedRoutes(router *gin.RouterGroup, db *gorm.DB) {
	handler := NewUserHandler(db)

	users := router.Group("/users")
	{
		users.GET("/", handler.GetAll)
		users.GET("/:id", handler.GetByID)
		users.PUT("/:id", handler.Update)
		users.DELETE("/:id", handler.Delete)
	}
}
