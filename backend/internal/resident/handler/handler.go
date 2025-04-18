package handler

import (
	"portarius/internal/resident/domain"
	"portarius/internal/resident/interfaces"
	"portarius/internal/resident/repository"
	resident "portarius/internal/resident/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	var (
		repo     domain.IResidentRepository      = repository.NewResidentRepository(db)
		importer interfaces.ICSVResidentImporter = resident.NewResidentImportService(db)
	)

	controller := NewResidentController(repo, importer)

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
