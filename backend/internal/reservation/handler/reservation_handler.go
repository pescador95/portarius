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

func (c *ReservationHandler) Create(ctx *gin.Context) {
	var input struct {
		ResidentID    *uint            `json:"resident_id" binding:"required"`
		Space         domain.SpaceType `json:"space" binding:"required"`
		StartTime     time.Time        `json:"start_time" binding:"required"`
		EndTime       time.Time        `json:"end_time" binding:"required"`
		Description   string           `json:"description"`
		PaymentMethod string           `json:"payment_method" binding:"required"`
		PaymentDate   *time.Time       `json:"payment_date"`
		KeysTakenAt   *time.Time       `json:"keys_taken_at"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.repo.CheckReservationConflict(string(input.Space), input.StartTime, input.EndTime, 0); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	reservation := domain.Reservation{
		ResidentID:    input.ResidentID,
		Space:         input.Space,
		StartTime:     input.StartTime,
		EndTime:       input.EndTime,
		Description:   input.Description,
		Status:        domain.StatusPending,
		PaymentStatus: domain.PaymentPending,
		PaymentMethod: domain.PaymentMethod(input.PaymentMethod),
		PaymentDate:   input.PaymentDate,
	}

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

func (c *ReservationHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input struct {
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		Description string    `json:"description"`
	}

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

func (c *ReservationHandler) Confirm(ctx *gin.Context) {
	c.UpdateStatus(ctx, domain.StatusConfirmed)
}

func (c *ReservationHandler) Cancel(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input struct {
		Reason string `json:"reason"`
	}

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
	reservation.CancellationReason = input.Reason

	if err := c.repo.Update(reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

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

func (c *ReservationHandler) Complete(ctx *gin.Context) {
	c.UpdateStatus(ctx, domain.StatusKeysReturned)
}

func (c *ReservationHandler) ConfirmPayment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input struct {
		Amount float64 `json:"amount" binding:"required"`
	}

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
	reservation.PaymentAmount = input.Amount
	reservation.PaymentDate = &now
	reservation.Status = domain.StatusConfirmed

	if err := c.repo.Update(reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

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

func (c *ReservationHandler) GetBySpace(ctx *gin.Context) {
	space := ctx.Param("space")

	reservations, err := c.repo.GetBySpace(space)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

func (c *ReservationHandler) GetByStatus(ctx *gin.Context) {
	status := ctx.Param("status")

	reservations, err := c.repo.GetByStatus(status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

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

func (c *ReservationHandler) ImportSalonReservations(ctx *gin.Context) {
	if err := c.importService.ImportReservationsFromCSV(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Reservas de salão importad com sucesso",
	})
}

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

func (c *ReservationHandler) ListSpaceTypes(ctx *gin.Context) {
	spaceTypes := []domain.SpaceType{
		domain.Salon1,
		domain.Salon2,
	}

	ctx.JSON(http.StatusOK, spaceTypes)
}

func (c *ReservationHandler) ListPaymentMethods(ctx *gin.Context) {
	paymentMethods := []domain.PaymentMethod{
		domain.PaymentMethodPix,
		domain.PaymentMethodBoleto,
	}

	ctx.JSON(http.StatusOK, paymentMethods)
}

func (c *ReservationHandler) ListPaymentStatuses(ctx *gin.Context) {
	paymentStatuses := []domain.PaymentStatus{
		domain.PaymentPending,
		domain.PaymentPaid,
		domain.PaymentRefunded,
	}

	ctx.JSON(http.StatusOK, paymentStatuses)
}
