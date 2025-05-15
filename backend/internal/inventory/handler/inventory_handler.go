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

// GetAll 	godoc
// @Summary List all inventory items
// @Description Get paginated list of all inventory items with optional filters
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" minimum(1) default(1)
// @Param pageSize query int false "Items per page" minimum(1) maximum(100) default(10)
// @Param search query string false "Search term"
// @Param sortBy query string false "Sort field" Enums(name,quantity,created_at)
// @Param sortOrder query string false "Sort order" Enums(asc,desc) default(asc)
// @Success 200 {array} domain.Inventory
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /inventory [get]
func (c *InventoryHandler) GetAll(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	items, err := c.repo.GetAll(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// GetByID 	godoc
// @Summary Get inventory item by ID
// @Description Get inventory item by ID
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inventory ID"
// @Success 200 {object} domain.Inventory
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /inventory/{id} [get]
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

// Create 	godoc
// @Summary Create a new inventory item
// @Description Create a new inventory item
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body domain.Inventory true "Inventory item"
// @Success 201 {object} domain.Inventory
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /inventory [post]
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

// Update 	godoc
// @Summary Update an inventory item
// @Description Update an inventory item
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inventory ID"
// @Param item body domain.Inventory true "Inventory item"
// @Success 200 {object} domain.Inventory
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /inventory/{id} [put]
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

// Delete 	godoc
// @Summary Delete an inventory item
// @Description Delete an inventory item
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inventory ID"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /inventory/{id} [delete]
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

// ImportPets 	godoc
// @Summary Import cars from CSV
// @Description Import pets from CSV file
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param file formData file true "CSV file"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /inventory/import-pets [post]
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

// ListInventoryTypes 	godoc
// @Summary List inventory types
// @Description List avaliable inventory types
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.InventoryType
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /inventory/inventory-types [get]
func (c *InventoryHandler) ListInventoryTypes(ctx *gin.Context) {
	types := []domain.InventoryType{
		domain.InventoryTypeCar,
		domain.InventoryTypeBike,
		domain.InventoryTypeBicycle,
		domain.InventoryTypeScooter,
		domain.InventoryTypePet,
	}
	ctx.JSON(http.StatusOK, types)
}
