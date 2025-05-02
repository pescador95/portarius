package handler

import (
	"net/http"
	"portarius/internal/resident/domain"
	"portarius/internal/resident/interfaces"

	"strconv"

	"github.com/gin-gonic/gin"
)

type ResidentHandler struct {
	repo          domain.IResidentRepository
	importService interfaces.ICSVResidentImporter
}

func NewResidentHandler(repo domain.IResidentRepository, importer interfaces.ICSVResidentImporter) *ResidentHandler {
	return &ResidentHandler{
		repo:          repo,
		importService: importer,
	}
}

func (c *ResidentHandler) GetAll(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	residents, err := c.repo.GetAll(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, residents)
}

func (c *ResidentHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	resident, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Morador não encontrado"})
		return
	}
	ctx.JSON(http.StatusOK, resident)
}

func (c *ResidentHandler) Create(ctx *gin.Context) {
	var resident domain.Resident
	if err := ctx.ShouldBindJSON(&resident); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.repo.Create(&resident); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, resident)
}

func (c *ResidentHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var resident domain.Resident
	if err := ctx.ShouldBindJSON(&resident); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resident.ID = uint(id)
	if err := c.repo.Update(&resident); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resident)
}

func (c *ResidentHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Morador excluído com sucesso"})
}

func (c *ResidentHandler) ImportResidents(ctx *gin.Context) {
	if err := c.importService.ImportResidentsFromCSV(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Moradores importados com sucesso"})
}
