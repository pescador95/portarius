package pkg

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PackageController struct {
	db *gorm.DB
}

func NewPackageController(db *gorm.DB) *PackageController {
	return &PackageController{db: db}
}

func (c *PackageController) GetAll(ctx *gin.Context) {
	var packages []Package
	if err := c.db.Find(&packages).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, packages)
}

func (c *PackageController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var pkg Package
	if err := c.db.First(&pkg, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Encomenda não encontrada"})
		return
	}
	ctx.JSON(http.StatusOK, pkg)
}

func (c *PackageController) Create(ctx *gin.Context) {
	var pkg Package
	if err := ctx.ShouldBindJSON(&pkg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg.ReceivedAt = time.Now()
	if err := c.db.Create(&pkg).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, pkg)
}

func (c *PackageController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var pkg Package
	if err := ctx.ShouldBindJSON(&pkg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg.ID = uint(id)
	if err := c.db.Save(&pkg).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pkg)
}

func (c *PackageController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.db.Delete(&Package{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Encomenda excluída com sucesso"})
}

func (c *PackageController) MarkAsDelivered(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var pkg Package
	if err := c.db.First(&pkg, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Encomenda não encontrada"})
		return
	}

	pkg.Status = PackageDelivered
	pkg.DeliveredAt = time.Now()
	if err := c.db.Save(&pkg).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pkg)
}

func (c *PackageController) MarkAsLost(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var pkg Package
	if err := c.db.First(&pkg, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Encomenda não encontrada"})
		return
	}

	pkg.Status = PackageLost
	pkg.DeliveredAt = time.Now()
	if err := c.db.Save(&pkg).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pkg)
}
