package reservation

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	controller := NewReservationController(db)

	reservations := router.Group("/reservations")
	{
		reservations.GET("/", controller.GetAll)
		reservations.GET("/:id", controller.GetByID)
		reservations.POST("/", controller.Create)
		reservations.PUT("/:id", controller.Update)
		reservations.DELETE("/:id", controller.Delete)

		reservations.PUT("/:id/confirm", controller.Confirm)
		reservations.PUT("/:id/cancel", controller.Cancel)
		reservations.PUT("/:id/take-keys", controller.TakeKeys)
		reservations.PUT("/:id/return-keys", controller.ReturnKeys)
		reservations.PUT("/:id/complete", controller.Complete)
		reservations.PUT("/:id/confirm-payment", controller.ConfirmPayment)

		reservations.GET("/resident/:residentId", controller.GetByResident)
		reservations.GET("/space/:space", controller.GetBySpace)
		reservations.GET("/status/:status", controller.GetByStatus)
		reservations.GET("/date-range", controller.GetByDateRange)
		reservations.GET("/upcoming", controller.GetUpcomingReservations)

		reservations.POST("/import-salon", controller.ImportSalonReservations)
	}
}
