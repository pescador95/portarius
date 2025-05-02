package routes

import (
	"portarius/internal/reminder/domain"
	reminderHandler "portarius/internal/reminder/handler"
	"portarius/internal/reminder/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterReminderProtectedRoutes(router *gin.RouterGroup, db *gorm.DB) {
	var (
		repo domain.IReminderRepository = repository.NewReminderRepository(db)
	)

	handler := reminderHandler.NewReminderHandler(repo)

	reminders := router.Group("/reminders")
	{
		reminders.GET("/", handler.GetAll)
		reminders.GET("/:id", handler.GetByID)
		reminders.POST("/", handler.Create)
		reminders.PUT("/:id", handler.Update)
		reminders.DELETE("/:id", handler.Delete)
		reminders.GET("/reservation/:reservationID", handler.GetByReservationID)
		reminders.GET("/package/:packageID", handler.GetByPackageID)
		reminders.GET("/status/:status", handler.GetByStatus)
		reminders.GET("/channel/:channel", handler.GetByChannel)
		reminders.GET("/recipient/:recipient", handler.GetByRecipient)
		reminders.GET("/pending", handler.GetByPendingStatus)
		reminders.GET("/reminderChannel", handler.ListReminderChannel)
		reminders.GET("/reminderStatus", handler.ListReminderStatus)
	}
}
