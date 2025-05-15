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

// GetAll godoc
// @Summary List all reminders
// @Description Get paginated list of all reminders with optional filters
// @Tags Reminders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" minimum(1) default(1)
// @Param pageSize query int false "Items per page" minimum(1) maximum(100) default(10)
// @Param search query string false "Search by recipient"
// @Param sortBy query string false "Sort field" Enums(recipient,scheduled_at,sent_at,status,channel,created_at)
// @Param sortOrder query string false "Sort order" Enums(asc,desc) default(asc)
// @Success 200 {array} domain.Reminder
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /reminders [get]
func (h *ReminderHandler) GetAll(c *gin.Context) {
	reminders, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reminders)
}

// GetByID godoc
// @Summary Get reminder by ID
// @Description Get a single reminder by its ID
// @Tags Reminders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reminder ID"
// @Success 200 {object} domain.Reminder
// @Failure 400
// @Failure 404
// @Router /reminders/{id} [get]
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

// Create godoc
// @Summary Create a new reminder
// @Description Create a new reminder with recipient, scheduled time, etc.
// @Tags Reminders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param reminder body domain.Reminder true "Reminder to create"
// @Success 201 {object} domain.Reminder
// @Failure 400
// @Failure 500
// @Router /reminders [post]
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

// Update godoc
// @Summary Update an existing reminder
// @Description Update reminder fields like scheduled time, status, etc.
// @Tags Reminders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reminder ID"
// @Param reminder body domain.Reminder true "Updated reminder"
// @Success 200 {object} domain.Reminder
// @Failure 400
// @Failure 500
// @Router /reminders/{id} [put]
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

// Delete godoc
// @Summary Delete a reminder
// @Description Delete a reminder by its ID
// @Tags Reminders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reminder ID"
// @Success 204 {string} string "No Content"
// @Failure 400
// @Failure 500
// @Router /reminders/{id} [delete]
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

// GetByReservationID godoc
// @Summary Get reminder by reservation ID
// @Description Get a reminder linked to a specific reservation
// @Tags Reminders
// @Produce json
// @Security BearerAuth
// @Param reservationID path int true "Reservation ID"
// @Success 200 {object} domain.Reminder
// @Failure 400
// @Failure 404
// @Router /reminders/reservation/{reservationID} [get]
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

// GetByPackageID godoc
// @Summary Get reminder by package ID
// @Description Get a reminder linked to a specific package
// @Tags Reminders
// @Produce json
// @Security BearerAuth
// @Param packageID path int true "Package ID"
// @Success 200 {object} domain.Reminder
// @Failure 400
// @Failure 404
// @Router /reminders/package/{packageID} [get]
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

// GetByStatus godoc
// @Summary Get reminders by status
// @Description Get a list of reminders filtered by status
// @Tags Reminders
// @Produce json
// @Security BearerAuth
// @Param status path string true "Reminder status" Enums(PENDING,SENT,FAILED,CANCELLED)
// @Success 200 {array} domain.Reminder
// @Failure 500
// @Router /reminders/status/{status} [get]
func (h *ReminderHandler) GetByStatus(c *gin.Context) {
	status := c.Param("status")

	reminders, err := h.repo.GetByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

// GetByChannel godoc
// @Summary Get reminders by channel
// @Description Get a list of reminders filtered by channel
// @Tags Reminders
// @Produce json
// @Security BearerAuth
// @Param channel path string true "Reminder channel" Enums(WHATSAPP,EMAIL,SMS,TELEGRAM,INSTAGRAM,FACEBOOK,DISCORD)
// @Success 200 {array} domain.Reminder
// @Failure 500
// @Router /reminders/channel/{channel} [get]
func (h *ReminderHandler) GetByChannel(c *gin.Context) {
	channel := c.Param("channel")

	reminders, err := h.repo.GetByChannel(channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

// GetByRecipient godoc
// @Summary Get reminders by recipient
// @Description Get all reminders sent to a specific recipient
// @Tags Reminders
// @Produce json
// @Security BearerAuth
// @Param recipient path string true "Recipient"
// @Success 200 {array} domain.Reminder
// @Failure 500
// @Router /reminders/recipient/{recipient} [get]
func (h *ReminderHandler) GetByRecipient(c *gin.Context) {
	recipient := c.Param("recipient")

	reminders, err := h.repo.GetByRecipient(recipient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

// GetByPendingStatus godoc
// @Summary Get reminders with status PENDING
// @Description Get all reminders that are still pending
// @Tags Reminders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Reminder
// @Failure 500
// @Router /reminders/pending [get]
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

// ListReminderChannel godoc
// @Summary List all reminder channels
// @Description Returns all supported reminder channels
// @Tags Reminders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} string
// @Router /reminders/reminderChannel [get]
func (h *ReminderHandler) ListReminderChannel(c *gin.Context) {
	reminderChannels := []domain.ReminderChannel{
		domain.ReminderChannelWhatsApp,
		domain.ReminderChannelEmail,
		domain.ReminderChannelSMS,
		domain.ReminderChannelTelegram,
		domain.ReminderChannelInstagram,
		domain.ReminderChannelFacebook,
		domain.ReminderChannelDiscord,
	}

	c.JSON(http.StatusOK, reminderChannels)
}

// ListReminderStatus godoc
// @Summary List all reminder statuses
// @Description Returns all possible reminder statuses
// @Tags Reminders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} string
// @Router /reminders/reminderStatus [get]
func (h *ReminderHandler) ListReminderStatus(c *gin.Context) {
	reminderStatuses := []domain.ReminderStatus{
		domain.ReminderStatusPending,
		domain.ReminderStatusSent,
		domain.ReminderStatusFailed,
		domain.ReminderStatusCancelled,
	}

	c.JSON(http.StatusOK, reminderStatuses)
}
