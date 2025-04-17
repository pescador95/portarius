package resident

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	controller := NewResidentController(db)

	residents := router.Group("/residents")
	{
		residents.GET("/", controller.GetAll)
		residents.GET("/:id", controller.GetByID)
		residents.POST("/", controller.Create)
		residents.PUT("/:id", controller.Update)
		residents.DELETE("/:id", controller.Delete)
		residents.POST("/import", controller.ImportResidents)
	}
}
