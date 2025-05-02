package reservation

import (
	"portarius/internal/reservation/domain"
	reservationHandler "portarius/internal/reservation/handler"
	"portarius/internal/reservation/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterReservationRoutes(router *gin.RouterGroup, db *gorm.DB) {
	var (
		repo domain.IReservationRepository = repository.NewReservationRepository(db)
	)

	handler := reservationHandler.NewReservationHandler(repo)

	reservations := router.Group("/reservations")
	{
		reservations.GET("/", handler.GetAll)
		reservations.GET("/:id", handler.GetByID)
		reservations.POST("/", handler.Create)
		reservations.PUT("/:id", handler.Update)
		reservations.DELETE("/:id", handler.Delete)

		reservations.PUT("/:id/confirm", handler.Confirm)
		reservations.PUT("/:id/cancel", handler.Cancel)
		reservations.PUT("/:id/take-keys", handler.TakeKeys)
		reservations.PUT("/:id/return-keys", handler.ReturnKeys)
		reservations.PUT("/:id/complete", handler.Complete)
		reservations.PUT("/:id/confirm-payment", handler.ConfirmPayment)

		reservations.GET("/resident/:residentId", handler.GetByResident)
		reservations.GET("/space/:space", handler.GetBySpace)
		reservations.GET("/status/:status", handler.GetByStatus)
		reservations.GET("/date-range", handler.GetByDateRange)
		reservations.GET("/upcoming", handler.GetUpcomingReservations)

		reservations.POST("/import-salon", handler.ImportSalonReservations)
		reservations.GET("/reservationStatus", handler.ListReservationStatus)
		reservations.GET("/spaceTypes", handler.ListSpaceTypes)
		reservations.GET("/paymentMethods", handler.ListPaymentMethods)
		reservations.GET("/paymentStatuses", handler.ListPaymentStatuses)
	}
}
