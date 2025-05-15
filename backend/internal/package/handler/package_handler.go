package handler

import (
	"net/http"
	"portarius/internal/eventbus"
	"portarius/internal/package/domain"
	reminderDomain "portarius/internal/reminder/domain"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PackageHandler struct {
	repo domain.IPackageRepository
}

func NewPackageHandler(repo domain.IPackageRepository) *PackageHandler {
	return &PackageHandler{repo: repo}
}

// GetAll godoc
// @summary List all package items
// @Description Get paginated list of all package items with optional filters
// @Tags Package
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" minimum(1) default(1)
// @Param pageSize query int false "Items per page" minimum(1) maximum(100) default(10)
// @Param search query string false "Search term"
// @Param sortBy query string false "Sort field" Enums(description,status,quantity,created_at)
// @Param sortOrder query string false "Sort order" Enums(asc,desc) default(asc)
// @Success 200 {array} domain.Package
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /package [get]
func (c *PackageHandler) GetAll(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	packages, err := c.repo.GetAll(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, packages)
}

// GetByID 	godoc
// @Summary Get package item by ID
// @Description Get package item by ID
// @Tags Package
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Package ID"
// @Success 200 {object} domain.Package
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /package/{id} [get]
func (c *PackageHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	pkg, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Encomenda não encontrada"})
		return
	}
	ctx.JSON(http.StatusOK, pkg)
}

// Create 	godoc
// @Summary Create a new package item
// @Description Create a new package item
// @Tags Package
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body domain.Package true "Package item"
// @Success 201 {object} domain.Package
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /package [post]
func (c *PackageHandler) Create(ctx *gin.Context) {
	var pkg domain.Package
	if err := ctx.ShouldBindJSON(&pkg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg.ReceivedAt = time.Now()
	if err := c.repo.Create(&pkg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if pkg.Status == domain.PackagePending {

		eventbus.Publish("PackageCreated", &eventbus.PackageCreatedEvent{
			PackageID: &pkg.ID,
			Channel:   string(reminderDomain.ReminderChannelWhatsApp),
		})
	}

	ctx.JSON(http.StatusCreated, pkg)
}

// Update 	godoc
// @Summary Update an package item
// @Description Update an package item
// @Tags Package
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Package ID"
// @Param item body domain.Package true "Package item"
// @Success 200 {object} domain.Package
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /package/{id} [put]
func (c *PackageHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var pkg domain.Package
	if err := ctx.ShouldBindJSON(&pkg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg.ID = uint(id)
	if err := c.repo.Update(&pkg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pkg)
}

// Delete 	godoc
// @Summary Delete an package item
// @Description Delete an package item
// @Tags Package
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Package ID"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /package/{id} [delete]
func (c *PackageHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Encomenda excluída com sucesso"})
}

// MarkAsDelivered 	godoc
// @Summary Mark a package as delivered
// @Description Change the package status
// @Tags Package
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Package ID"
// @Success 200 {object} domain.Package
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /package/{id}/deliver [put]
func (c *PackageHandler) MarkAsDelivered(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	pkg, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Encomenda não encontrada"})
		return
	}

	pkg.Status = domain.PackageDelivered
	pkg.DeliveredAt = time.Now()
	if err := c.repo.Update(pkg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pkg)
}

// MarkAsLost 	godoc
// @Summary Mark a package as lost
// @Description Change the package status
// @Tags Package
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Package ID"
// @Success 200 {object} domain.Package
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /package/{id}/lost [put]
func (c *PackageHandler) MarkAsLost(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	pkg, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Encomenda não encontrada"})
		return
	}

	pkg.Status = domain.PackageLost
	pkg.DeliveredAt = time.Now()
	if err := c.repo.Update(pkg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pkg)
}

// ListPackageStatus 	godoc
// @Summary List package status
// @Description List avaliable package status
// @Tags Package
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.PackageStatus
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /package/status [get]
func (c *PackageHandler) ListPackageStatus(ctx *gin.Context) {
	status := []domain.PackageStatus{
		domain.PackagePending,
		domain.PackageDelivered,
		domain.PackageLost,
	}

	ctx.JSON(http.StatusOK, status)
}
