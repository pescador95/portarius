package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	controller := NewUserController(db)

	auth := router.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
	}
}

func RegisterProtectedRoutes(router *gin.RouterGroup, db *gorm.DB) {
	controller := NewUserController(db)

	users := router.Group("/users")
	{
		users.GET("/", controller.GetAll)
		users.GET("/:id", controller.GetByID)
		users.PUT("/:id", controller.Update)
		users.DELETE("/:id", controller.Delete)
	}
}
