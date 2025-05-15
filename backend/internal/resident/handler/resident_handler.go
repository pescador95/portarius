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

// GetAll godoc
// @Summary Get all residents with pagination
// @Description Retrieves a paginated list of residents. Use query params 'page' and 'pageSize' for pagination.
// @Tags Residents
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {array} domain.Resident "List of residents"
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /residents/ [get]
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

// GetByID godoc
// @Summary Get resident by ID
// @Description Retrieves a resident by their ID.
// @Tags Residents
// @Produce json
// @Param id path int true "Resident ID"
// @Success 200 {object} domain.Resident "Resident found"
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /residents/{id} [get]
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

// Create godoc
// @Summary Create a new resident
// @Description Creates a new resident with the provided JSON body.
// @Tags Residents
// @Accept json
// @Produce json
// @Param resident body domain.Resident true "Resident data"
// @Success 201 {object} domain.Resident "Resident created"
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /residents/ [post]
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

// Update godoc
// @Summary Update an existing resident
// @Description Updates the resident with the given ID using the JSON body.
// @Tags Residents
// @Accept json
// @Produce json
// @Param id path int true "Resident ID"
// @Param resident body domain.Resident true "Resident data"
// @Success 200 {object} domain.Resident "Resident updated"
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /residents/{id} [put]
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

// Delete godoc
// @Summary Delete a resident by ID
// @Description Deletes the resident with the specified ID.
// @Tags Residents
// @Produce json
// @Param id path int true "Resident ID"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /residents/{id} [delete]
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

// ImportResidents godoc
// @Summary Import residents from CSV
// @Description Imports residents data from a CSV file (implementation-specific).
// @Tags Residents
// @Produce json
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /residents/import [post]
func (c *ResidentHandler) ImportResidents(ctx *gin.Context) {
	if err := c.importService.ImportResidentsFromCSV(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Moradores importados com sucesso"})
}

// ListResidentType godoc
// @Summary List all resident types
// @Description Returns the list of possible resident types.
// @Tags Residents
// @Produce json
// @Success 200 {array} domain.ResidentType
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /residents/residentType [get]
func (c *ResidentHandler) ListResidentType(ctx *gin.Context) {
	residentTypes := []domain.ResidentType{
		domain.Tenant,
		domain.Owner,
		domain.Krum,
		domain.NotResident,
	}
	ctx.JSON(http.StatusOK, residentTypes)
}
