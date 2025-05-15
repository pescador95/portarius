package reservation

import (
	"net/http"
	"portarius/internal/eventbus"
	reminderDomain "portarius/internal/reminder/domain"
	"portarius/internal/reservation/domain"
	"portarius/internal/reservation/interfaces"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	repo          domain.IReservationRepository
	importService interfaces.ICSVReservationImporter
}

func NewReservationHandler(repo domain.IReservationRepository) *ReservationHandler {
	return &ReservationHandler{repo: repo}
}

// GetAll godoc
// @Summary Get all reservations
// @Description Retrieves a paginated list of all reservations
// @Tags Reservations
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Success 200 {array} domain.Reservation
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /reservations [get]
func (c *ReservationHandler) GetAll(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	reservations, err := c.repo.GetAll(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}

// GetByID godoc
// @Summary Get a reservation by ID
// @Description Retrieves a reservation based on its unique ID
// @Tags Reservations
// @Accept json
// @Produce json
// @Param id path int true "Reservation ID"
// @Success 200 {object} domain.Reservation
// @Failure 400
// @Failure 404
// @Router /reservations/{id} [get]
func (c *ReservationHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	reservation, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}
	ctx.JSON(http.StatusOK, reservation)
}

// Create godoc
// @Summary Create a new reservation
// @Description Creates a new reservation with the provided details
// @Tags Reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param reservation body domain.Reservation true "Reservation data"
// @Success 201 {object} domain.Reservation
// @Failure 400
// @Failure 409
// @Failure 500
// @Router /reservations [post]
func (c *ReservationHandler) Create(ctx *gin.Context) {
	var reservation domain.Reservation
	if err := ctx.ShouldBindJSON(&reservation); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.repo.CheckReservationConflict(string(reservation.Space), reservation.StartTime, reservation.EndTime, 0); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	reservation.Status = domain.StatusPending
	reservation.PaymentStatus = domain.PaymentPending

	if err := c.repo.Create(&reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if reservation.Status == domain.StatusPending || reservation.Status == domain.StatusConfirmed {

		eventbus.Publish("ReservationCreated", &eventbus.ReservationCreatedEvent{
			ReservationID: &reservation.ID,
			Channel:       string(reminderDomain.ReminderChannelWhatsApp),
			StartTime:     reservation.StartTime,
		})
	}

	ctx.JSON(http.StatusCreated, reservation)
}

// Update godoc
// @Summary Update a reservation
// @Description Update reservation's time or description
// @Tags Reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Param reservation body domain.Reservation true "Updated data"
// @Success 200 {object} domain.Reservation
// @Failure 400
// @Failure 404
// @Failure 409
// @Failure 500
// @Router /reservations/{id} [put]
func (c *ReservationHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input domain.Reservation

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	if !reservation.StartTime.Equal(input.StartTime) || !reservation.EndTime.Equal(input.EndTime) {
		if err := c.repo.CheckReservationConflict(string(reservation.Space), input.StartTime, input.EndTime, reservation.ID); err != nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		reservation.StartTime = input.StartTime
		reservation.EndTime = input.EndTime
	}

	reservation.Description = input.Description

	if err := c.repo.Update(reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

// Delete godoc
// @Summary Delete a reservation
// @Description Delete a reservation by its ID
// @Tags Reservations
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /reservations/{id} [delete]
func (c *ReservationHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Reserva excluída com sucesso"})
}

// Confirm godoc
// @Summary Confirm a reservation
// @Description Update the reservation status to confirmed
// @Tags Reservations
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} domain.Reservation
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /reservations/{id}/confirm [put]
func (c *ReservationHandler) Confirm(ctx *gin.Context) {
	c.UpdateStatus(ctx, domain.StatusConfirmed)
}

// Cancel godoc
// @Summary Cancel a reservation
// @Description Cancel a reservation by setting its status to cancelled and adding a cancellation reason
// @Tags Reservations
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Param body body domain.Reservation true "Cancellation body (only CancellationReason will be used)"
// @Success 200 {object} domain.Reservation
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /reservations/{id}/cancel [put]
func (c *ReservationHandler) Cancel(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input domain.Reservation

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := c.repo.GetByID(uint(id))

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	reservation.Status = domain.StatusCancelled
	reservation.CancellationReason = input.CancellationReason

	if err := c.repo.Update(reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

// TakeKeys godoc
// @Summary Mark reservation keys as taken
// @Description Marks the keys as taken for a confirmed reservation and updates its status
// @Tags Reservations
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} domain.Reservation
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /reservations/{id}/take-keys [put]
func (c *ReservationHandler) TakeKeys(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	reservation, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	if reservation.Status != domain.StatusConfirmed {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "A reserva precisa estar confirmada para retirar as chaves"})
		return
	}

	now := time.Now()
	reservation.KeysTakenAt = &now
	reservation.Status = domain.StatusKeysTaken

	if err := c.repo.Update(reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

// ReturnKeys godoc
// @Summary Mark reservation keys as returned
// @Description Marks the keys as returned for a reservation and updates its status
// @Tags Reservations
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} domain.Reservation
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /reservations/{id}/return-keys [put]
func (c *ReservationHandler) ReturnKeys(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	reservation, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	if reservation.Status != domain.StatusKeysTaken {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "As chaves precisam ter sido retiradas primeiro"})
		return
	}

	now := time.Now()
	reservation.KeysReturnedAt = &now
	reservation.Status = domain.StatusKeysReturned

	if err := c.repo.Update(reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

// Complete godoc
// @Summary Mark reservation as complete
// @Description Updates the reservation status to "keys returned" (complete)
// @Tags Reservations
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} domain.Reservation
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /reservations/{id}/complete [put]
func (c *ReservationHandler) Complete(ctx *gin.Context) {
	c.UpdateStatus(ctx, domain.StatusKeysReturned)
}

// ConfirmPayment godoc
// @Summary Confirm payment for a reservation
// @Description Confirm the payment of a reservation, updating payment status, amount, date, and reservation status
// @Tags Reservations
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Param body body domain.Reservation true "Reservation payment details (payment_amount)"
// @Success 200 {object} domain.Reservation
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /reservations/{id}/confirm-payment [put]
func (c *ReservationHandler) ConfirmPayment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input domain.Reservation

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	if reservation.PaymentStatus == domain.PaymentPaid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Pagamento já foi confirmado anteriormente"})
		return
	}

	now := time.Now()
	reservation.PaymentStatus = domain.PaymentPaid
	reservation.PaymentAmount = input.PaymentAmount
	reservation.PaymentDate = &now
	reservation.Status = domain.StatusConfirmed

	if err := c.repo.Update(reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

// GetByResident godoc
// @Summary Get reservations by resident ID
// @Description Retrieve all reservations associated with a specific resident
// @Tags Reservations
// @Security BearerAuth
// @Param residentId path int true "Resident ID"
// @Success 200 {array} domain.Reservation
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /reservations/resident/{residentId} [get]
func (c *ReservationHandler) GetByResident(ctx *gin.Context) {
	residentID, err := strconv.ParseUint(ctx.Param("residentId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID do morador inválido"})
		return
	}

	reservations, err := c.repo.GetByResident(uint(residentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

// GetBySpace godoc
// @Summary Get reservations by space type
// @Description Retrieve all reservations for a specific space
// @Tags Reservations
// @Security BearerAuth
// @Param space path string true "Space type"
// @Success 200 {array} domain.Reservation
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /reservations/space/{space} [get]
func (c *ReservationHandler) GetBySpace(ctx *gin.Context) {
	space := ctx.Param("space")

	reservations, err := c.repo.GetBySpace(space)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

// GetByStatus godoc
// @Summary Get reservations by status
// @Description Retrieve all reservations filtered by their status
// @Tags Reservations
// @Security BearerAuth
// @Param status path string true "Reservation status"
// @Success 200 {array} domain.Reservation
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /reservations/status/{status} [get]
func (c *ReservationHandler) GetByStatus(ctx *gin.Context) {
	status := ctx.Param("status")

	reservations, err := c.repo.GetByStatus(status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

// GetByDateRange godoc
// @Summary Get reservations by date range
// @Description Retrieve all reservations between start_date and end_date (format: yyyy-MM-dd)
// @Tags Reservations
// @Security BearerAuth
// @Param start_date query string true "Start date in format yyyy-MM-dd"
// @Param end_date query string true "End date in format yyyy-MM-dd"
// @Success 200 {array} domain.Reservation
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /reservations/date-range [get]
func (c *ReservationHandler) GetByDateRange(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if startDate == "" || endDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data inicial e final são obrigatórias"})
		return
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data inicial inválida"})
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data final inválida"})
		return
	}

	reservations, err := c.repo.FindByDateRange(start, end)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}

// GetUpcomingReservations godoc
// @Summary Get upcoming reservations
// @Description Retrieve all upcoming reservations
// @Tags Reservations
// @Security BearerAuth
// @Success 200 {array} domain.Reservation
// @Failure 401
// @Failure 500
// @Router /reservations/upcoming [get]
func (c *ReservationHandler) GetUpcomingReservations(ctx *gin.Context) {
	reservations, err := c.repo.FindUpcomingReservations()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}

func (c *ReservationHandler) UpdateStatus(ctx *gin.Context, status domain.ReservationStatus) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	reservation, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	reservation.Status = status

	if err := c.repo.Update(reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

// ImportSalonReservations godoc
// @Summary Import salon reservations from CSV
// @Description Imports reservations data from a CSV file for salon spaces
// @Tags Reservations
// @Security BearerAuth
// @Success 200 {object} map[string]string "Success message"
// @Failure 500
// @Router /reservations/import-salon [post]
func (c *ReservationHandler) ImportSalonReservations(ctx *gin.Context) {
	if err := c.importService.ImportReservationsFromCSV(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Reservas de salão importadas com sucesso",
	})
}

// ListReservationStatus godoc
// @Summary List all reservation statuses
// @Description Returns the list of possible reservation statuses
// @Tags Reservations
// @Produce json
// @Success 200 {array} domain.ReservationStatus "List of reservation statuses"
// @Router /reservations/reservationStatus [get]
func (c *ReservationHandler) ListReservationStatus(ctx *gin.Context) {
	reservationStatus := []domain.ReservationStatus{
		domain.StatusPending,
		domain.StatusConfirmed,
		domain.StatusCancelled,
		domain.StatusKeysTaken,
		domain.StatusKeysReturned,
	}

	ctx.JSON(http.StatusOK, reservationStatus)
}

// ListSpaceTypes godoc
// @Summary List all space types
// @Description Returns the list of possible space types for reservations
// @Tags Reservations
// @Produce json
// @Success 200 {array} domain.SpaceType "List of space types"
// @Router /reservations/spaceTypes [get]
func (c *ReservationHandler) ListSpaceTypes(ctx *gin.Context) {
	spaceTypes := []domain.SpaceType{
		domain.Salon1,
		domain.Salon2,
	}

	ctx.JSON(http.StatusOK, spaceTypes)
}

// ListPaymentMethods godoc
// @Summary List all payment methods
// @Description Returns the list of possible payment methods
// @Tags Reservations
// @Produce json
// @Success 200 {array} domain.PaymentMethod "List of payment methods"
// @Router /reservations/paymentMethods [get]
func (c *ReservationHandler) ListPaymentMethods(ctx *gin.Context) {
	paymentMethods := []domain.PaymentMethod{
		domain.PaymentMethodPix,
		domain.PaymentMethodBoleto,
	}

	ctx.JSON(http.StatusOK, paymentMethods)
}

// ListPaymentStatuses godoc
// @Summary List all payment statuses
// @Description Returns the list of possible payment statuses
// @Tags Reservations
// @Produce json
// @Success 200 {array} domain.PaymentStatus "List of payment statuses"
// @Router /reservations/paymentStatuses [get]
func (c *ReservationHandler) ListPaymentStatuses(ctx *gin.Context) {
	paymentStatuses := []domain.PaymentStatus{
		domain.PaymentPending,
		domain.PaymentPaid,
		domain.PaymentRefunded,
	}

	ctx.JSON(http.StatusOK, paymentStatuses)
}
