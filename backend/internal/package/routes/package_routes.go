package handler

import (
	"portarius/internal/package/domain"
	packageHandler "portarius/internal/package/handler"
	"portarius/internal/package/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterPackageRoutes(router *gin.RouterGroup, db *gorm.DB) {
	var (
		repo domain.IPackageRepository = repository.NewPackageRepository(db)
	)

	handler := packageHandler.NewPackageHandler(repo)

	packages := router.Group("/packages")
	{
		packages.GET("/", handler.GetAll)
		packages.GET("/:id", handler.GetByID)
		packages.POST("/", handler.Create)
		packages.PUT("/:id", handler.Update)
		packages.DELETE("/:id", handler.Delete)
		packages.PUT("/:id/deliver", handler.MarkAsDelivered)
		packages.PUT("/:id/lost", handler.MarkAsLost)
	}
}
