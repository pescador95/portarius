package handler

import (
	"net/http"
	"portarius/internal/inventory/domain"
	"portarius/internal/inventory/interfaces"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	repo          domain.IInventoryRepository
	importService interfaces.ICSVInventoryImporter
}

func NewInventoryHandler(repo domain.IInventoryRepository, importer interfaces.ICSVInventoryImporter) *InventoryHandler {
	return &InventoryHandler{
		repo:          repo,
		importService: importer,
	}
}

func (c *InventoryHandler) GetAll(ctx *gin.Context) {
	items, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (c *InventoryHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	item, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *InventoryHandler) Create(ctx *gin.Context) {
	var item domain.Inventory
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.repo.Create(&item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, item)
}

func (c *InventoryHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var item domain.Inventory
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.ID = uint(id)
	if err := c.repo.Update(&item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *InventoryHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Item excluído com sucesso"})
}

func (c *InventoryHandler) ImportPets(ctx *gin.Context) {
	if err := c.importService.ImportPetsFromCSV(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Pets importados com sucesso",
	})
}
