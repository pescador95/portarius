package handler

import (
	"portarius/internal/inventory/domain"
	"portarius/internal/inventory/interfaces"
	"portarius/internal/inventory/repository"
	inventoryService "portarius/internal/inventory/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterInventoryRoutes(router *gin.RouterGroup, db *gorm.DB) {

	var (
		repo     domain.IInventoryRepository      = repository.NewInventoryRepository(db)
		importer interfaces.ICSVInventoryImporter = inventoryService.NewInventoryImportService(db)
	)

	handler := NewInventoryHandler(repo, importer)

	inventory := router.Group("/inventory")
	{
		inventory.GET("/", handler.GetAll)
		inventory.GET("/:id", handler.GetByID)
		inventory.POST("/", handler.Create)
		inventory.PUT("/:id", handler.Update)
		inventory.DELETE("/:id", handler.Delete)
		inventory.POST("/import-pets", handler.ImportPets)
	}
}
