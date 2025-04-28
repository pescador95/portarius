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

func (c *PackageHandler) GetAll(ctx *gin.Context) {
	packages, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, packages)
}

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

	eventbus.Publish("PackageCreated", eventbus.PackageCreatedEvent{
		PackageID: pkg.ID,
		Channel:   string(reminderDomain.ReminderChannelWhatsApp),
	})

	ctx.JSON(http.StatusCreated, pkg)
}

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
