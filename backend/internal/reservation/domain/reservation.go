package domain

import (
	residentDomain "portarius/internal/resident/domain"
	"time"

	"gorm.io/gorm"
)

type SpaceType string

const (
	Salon1 SpaceType = "SALAO_1"
	Salon2 SpaceType = "SALAO_2"
)

type ReservationStatus string

const (
	StatusPending      ReservationStatus = "PENDENTE"
	StatusConfirmed    ReservationStatus = "CONFIRMADA"
	StatusCancelled    ReservationStatus = "CANCELADA"
	StatusKeysTaken    ReservationStatus = "CHAVES_RETIRADAS"
	StatusKeysReturned ReservationStatus = "CHAVES_DEVOLVIDAS"
)

type PaymentMethod string

const (
	PaymentMethodPix    PaymentMethod = "PIX"
	PaymentMethodBoleto PaymentMethod = "BOLETO"
)

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "PAGAMENTO_PENDENTE"
	PaymentPaid     PaymentStatus = "PAGO"
	PaymentRefunded PaymentStatus = "REEMBOLSADO"
)

const (
	CommonPaymentAmount  = 45.00
	HolydayPaymentAmount = 70.00
)

type Reservation struct {
	gorm.Model
	ResidentID    uint                     `json:"resident_id"`
	Resident      *residentDomain.Resident `json:"resident" gorm:"foreignKey:ResidentID"`
	Space         SpaceType                `json:"space" gorm:"not null;type:varchar(10);default:'SALAO_1'"`
	StartTime     time.Time                `json:"start_time" gorm:"not null"`
	EndTime       time.Time                `json:"end_time"`
	Status        ReservationStatus        `json:"status" gorm:"type:varchar(20);not null;default:'PENDENTE'"`
	PaymentStatus PaymentStatus            `json:"payment_status" gorm:"type:varchar(20);not null;default:'PAGAMENTO_PENDENTE'"`
	Description   string                   `json:"description"`

	KeysTakenAt    *time.Time `json:"keys_taken_at" gorm:"type:timestamp"`
	KeysReturnedAt *time.Time `json:"keys_returned_at" gorm:"type:timestamp"`

	PaymentAmount float64    `json:"payment_amount" gorm:"type:decimal(10,2);default:0"`
	PaymentDate   *time.Time `json:"payment_date" gorm:"type:timestamp"`

	PaymentMethod PaymentMethod `json:"payment_method" gorm:"type:varchar(20);not null;"`

	CancellationReason string `json:"cancellation_reason" gorm:"type:text"`
}
