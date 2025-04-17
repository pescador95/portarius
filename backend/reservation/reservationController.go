package reservation

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReservationController struct {
	db            *gorm.DB
	importService *ReservationImportService
}

func NewReservationController(db *gorm.DB) *ReservationController {
	return &ReservationController{
		db:            db,
		importService: NewReservationImportService(db),
	}
}

func (c *ReservationController) GetAll(ctx *gin.Context) {
	var reservations []Reservation
	if err := c.db.Preload("Resident").Find(&reservations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}

func (c *ReservationController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var reservation Reservation
	if err := c.db.Preload("Resident").First(&reservation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}
	ctx.JSON(http.StatusOK, reservation)
}

func (c *ReservationController) Create(ctx *gin.Context) {
	var input struct {
		ResidentID  uint      `json:"resident_id" binding:"required"`
		Space       SpaceType `json:"space" binding:"required"`
		StartTime   time.Time `json:"start_time" binding:"required"`
		EndTime     time.Time `json:"end_time" binding:"required"`
		Description string    `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.checkReservationConflict(input.Space, input.StartTime, input.EndTime, 0); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	reservation := Reservation{
		ResidentID:    input.ResidentID,
		Space:         input.Space,
		StartTime:     input.StartTime,
		EndTime:       input.EndTime,
		Description:   input.Description,
		Status:        StatusPending,
		PaymentStatus: PaymentPending,
	}

	if err := c.db.Create(&reservation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, reservation)
}

func (c *ReservationController) Update(ctx *gin.Context) {
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

	var reservation Reservation
	if err := c.db.First(&reservation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	if !reservation.StartTime.Equal(input.StartTime) || !reservation.EndTime.Equal(input.EndTime) {
		if err := c.checkReservationConflict(reservation.Space, input.StartTime, input.EndTime, reservation.ID); err != nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		reservation.StartTime = input.StartTime
		reservation.EndTime = input.EndTime
	}

	reservation.Description = input.Description

	if err := c.db.Save(&reservation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

func (c *ReservationController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.db.Delete(&Reservation{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Reserva excluída com sucesso"})
}

func (c *ReservationController) Confirm(ctx *gin.Context) {
	c.updateStatus(ctx, StatusConfirmed)
}

func (c *ReservationController) Cancel(ctx *gin.Context) {
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

	var reservation Reservation
	if err := c.db.First(&reservation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	reservation.Status = StatusCancelled
	reservation.CancellationReason = input.Reason

	if err := c.db.Save(&reservation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

func (c *ReservationController) TakeKeys(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var reservation Reservation
	if err := c.db.First(&reservation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	if reservation.Status != StatusConfirmed {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "A reserva precisa estar confirmada para retirar as chaves"})
		return
	}

	now := time.Now()
	reservation.KeysTakenAt = &now
	reservation.Status = StatusKeysTaken

	if err := c.db.Save(&reservation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

func (c *ReservationController) ReturnKeys(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var reservation Reservation
	if err := c.db.First(&reservation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	if reservation.Status != StatusKeysTaken {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "As chaves precisam ter sido retiradas primeiro"})
		return
	}

	now := time.Now()
	reservation.KeysReturnedAt = &now
	reservation.Status = StatusKeysReturned

	if err := c.db.Save(&reservation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

func (c *ReservationController) Complete(ctx *gin.Context) {
	c.updateStatus(ctx, StatusKeysReturned)
}

func (c *ReservationController) ConfirmPayment(ctx *gin.Context) {
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

	var reservation Reservation
	if err := c.db.First(&reservation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	if reservation.PaymentStatus == PaymentPaid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Pagamento já foi confirmado anteriormente"})
		return
	}

	now := time.Now()
	reservation.PaymentStatus = PaymentPaid
	reservation.PaymentAmount = input.Amount
	reservation.PaymentDate = &now
	reservation.Status = StatusConfirmed

	if err := c.db.Save(&reservation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

func (c *ReservationController) GetByResident(ctx *gin.Context) {
	residentID, err := strconv.ParseUint(ctx.Param("residentId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID do morador inválido"})
		return
	}

	var reservations []Reservation
	if err := c.db.Where("resident_id = ?", residentID).Find(&reservations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}

func (c *ReservationController) GetBySpace(ctx *gin.Context) {
	space := ctx.Param("space")

	var reservations []Reservation
	if err := c.db.Where("space = ?", space).Find(&reservations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}

func (c *ReservationController) GetByStatus(ctx *gin.Context) {
	status := ctx.Param("status")

	var reservations []Reservation
	if err := c.db.Where("status = ?", status).Find(&reservations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}

func (c *ReservationController) GetByDateRange(ctx *gin.Context) {
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

	var reservations []Reservation
	if err := c.db.Where("start_time BETWEEN ? AND ?", start, end).Find(&reservations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}

func (c *ReservationController) GetUpcomingReservations(ctx *gin.Context) {
	var reservations []Reservation
	if err := c.db.Where("start_time > ? AND status != ?", time.Now(), StatusCancelled).
		Order("start_time ASC").
		Find(&reservations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}

func (c *ReservationController) updateStatus(ctx *gin.Context, status ReservationStatus) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var reservation Reservation
	if err := c.db.First(&reservation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	reservation.Status = status

	if err := c.db.Save(&reservation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

func (c *ReservationController) checkReservationConflict(space SpaceType, startTime, endTime time.Time, excludeID uint) error {
	var conflictingReservation Reservation
	query := c.db.Where(
		"space = ? AND status NOT IN (?) AND ((start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?) OR (start_time <= ? AND end_time >= ?))",
		space,
		[]ReservationStatus{StatusCancelled, StatusKeysReturned},
		startTime,
		endTime,
		startTime,
		endTime,
		startTime,
		endTime,
	)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.First(&conflictingReservation).Error; err == nil {
		return errors.New("já existe uma reserva para este salão no horário selecionado")
	}

	return nil
}

func (c *ReservationController) ImportSalonReservations(ctx *gin.Context) {
	if err := c.importService.ImportSalonReservationsFromCSV(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Salon reservations imported successfully",
	})
}
