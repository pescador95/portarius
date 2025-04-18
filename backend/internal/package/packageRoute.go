package pkg

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	controller := NewPackageController(db)

	packages := router.Group("/packages")
	{
		packages.GET("/", controller.GetAll)
		packages.GET("/:id", controller.GetByID)
		packages.POST("/", controller.Create)
		packages.PUT("/:id", controller.Update)
		packages.DELETE("/:id", controller.Delete)
		packages.PUT("/:id/deliver", controller.MarkAsDelivered)
		packages.PUT("/:id/lost", controller.MarkAsLost)
	}
}
