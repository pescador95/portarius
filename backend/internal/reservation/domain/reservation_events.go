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
