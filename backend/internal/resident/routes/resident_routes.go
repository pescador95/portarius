package handler

import (
	"portarius/internal/resident/domain"
	residentHandler "portarius/internal/resident/handler"
	"portarius/internal/resident/interfaces"
	"portarius/internal/resident/repository"
	residentService "portarius/internal/resident/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ResidentRegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	var (
		repo     domain.IResidentRepository      = repository.NewResidentRepository(db)
		importer interfaces.ICSVResidentImporter = residentService.NewResidentImportService(db)
	)

	handler := residentHandler.NewResidentHandler(repo, importer)

	residents := router.Group("/residents")
	{
		residents.GET("/", handler.GetAll)
		residents.GET("/:id", handler.GetByID)
		residents.POST("/", handler.Create)
		residents.PUT("/:id", handler.Update)
		residents.DELETE("/:id", handler.Delete)
		residents.POST("/import", handler.ImportResidents)
	}
}
