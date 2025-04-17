package inventory

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	controller := NewInventoryController(db)

	inventory := router.Group("/inventory")
	{
		inventory.GET("/", controller.GetAll)
		inventory.GET("/:id", controller.GetByID)
		inventory.POST("/", controller.Create)
		inventory.PUT("/:id", controller.Update)
		inventory.DELETE("/:id", controller.Delete)
		inventory.POST("/import-pets", controller.ImportPets)
	}
}
