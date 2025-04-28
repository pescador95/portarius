package handler

import (
	"net/http"
	"portarius/internal/reminder/domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReminderHandler struct {
	repo domain.IReminderRepository
}

func NewReminderHandler(repo domain.IReminderRepository) *ReminderHandler {
	return &ReminderHandler{
		repo: repo,
	}
}
func (h *ReminderHandler) GetAll(c *gin.Context) {
	reminders, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reminders)
}

func (h *ReminderHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	reminder, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found"})
		return
	}

	c.JSON(http.StatusOK, reminder)
}

func (h *ReminderHandler) Create(c *gin.Context) {
	var reminder domain.Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(&reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

func (h *ReminderHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var reminder domain.Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminder.ID = uint(id)

	if err := h.repo.Update(&reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminder)
}

func (h *ReminderHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ReminderHandler) GetByReservationID(c *gin.Context) {
	idParam := c.Param("reservationID")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	reminder, err := h.repo.GetByReservationID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found for reservation"})
		return
	}

	c.JSON(http.StatusOK, reminder)
}

func (h *ReminderHandler) GetByPackageID(c *gin.Context) {
	idParam := c.Param("packageID")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid package id"})
		return
	}

	reminder, err := h.repo.GetByPackageID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found for package"})
		return
	}

	c.JSON(http.StatusOK, reminder)
}

func (h *ReminderHandler) GetByStatus(c *gin.Context) {
	status := c.Param("status")

	reminders, err := h.repo.GetByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *ReminderHandler) GetByChannel(c *gin.Context) {
	channel := c.Param("channel")

	reminders, err := h.repo.GetByChannel(channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *ReminderHandler) GetByRecipient(c *gin.Context) {
	recipient := c.Param("recipient")

	reminders, err := h.repo.GetByRecipient(recipient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *ReminderHandler) GetByPendingStatus(c *gin.Context) {
	reminders, err := h.repo.GetByPendingStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *ReminderHandler) GetPendingRemindersFromReservations(c *gin.Context) {
	reminders, err := h.repo.GetPendingRemindersFromReservations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *ReminderHandler) GetPendingRemindersFromPackages(c *gin.Context) {
	reminders, err := h.repo.GetPendingRemindersFromPackages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}
