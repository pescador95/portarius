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

	controller := NewInventoryHandler(repo, importer)

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
