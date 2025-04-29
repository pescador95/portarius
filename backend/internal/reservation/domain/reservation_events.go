package domain

import (
	holydayHandler "portarius/internal/holyday/handler"
	"time"

	"gorm.io/gorm"
)

func (r *Reservation) BeforeSave(tx *gorm.DB) (err error) {

	weekday := r.StartTime.Weekday()
	switch {
	case weekday == time.Friday || weekday == time.Saturday || weekday == time.Sunday || holydayHandler.IsHolyday(r.StartTime):
		r.PaymentAmount = HolydayPaymentAmount
	default:
		r.PaymentAmount = CommonPaymentAmount
	}

	return
}

func GetReminderScheduleDate(reservationDate time.Time, isHolidayFunc func(time.Time) bool) time.Time {
	scheduled := time.Date(
		reservationDate.Year(),
		reservationDate.Month(),
		reservationDate.Day(),
		12, 0, 0, 0,
		reservationDate.Location(),
	)

	for scheduled.Weekday() == time.Saturday || scheduled.Weekday() == time.Sunday || isHolidayFunc(scheduled) {
		scheduled = scheduled.AddDate(0, 0, -1)
	}

	return time.Date(scheduled.Year(), scheduled.Month(), scheduled.Day(), 12, 0, 0, 0, scheduled.Location())
}
