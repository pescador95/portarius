package resident

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResidentController struct {
	db            *gorm.DB
	importService *ResidentImportService
}

func NewResidentController(db *gorm.DB) *ResidentController {
	return &ResidentController{
		db:            db,
		importService: NewResidentImportService(db),
	}
}

func (c *ResidentController) GetAll(ctx *gin.Context) {
	var residents []Resident
	if err := c.db.Find(&residents).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, residents)
}

func (c *ResidentController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var resident Resident
	if err := c.db.First(&resident, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Morador não encontrado"})
		return
	}
	ctx.JSON(http.StatusOK, resident)
}

func (c *ResidentController) Create(ctx *gin.Context) {
	var resident Resident
	if err := ctx.ShouldBindJSON(&resident); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.db.Create(&resident).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, resident)
}

func (c *ResidentController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var resident Resident
	if err := ctx.ShouldBindJSON(&resident); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resident.ID = uint(id)
	if err := c.db.Save(&resident).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resident)
}

func (c *ResidentController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.db.Delete(&Resident{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Morador excluído com sucesso"})
}

func (c *ResidentController) ImportResidents(ctx *gin.Context) {
	if err := c.importService.ImportResidentsFromCSV(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Residents imported successfully",
	})
}
