package domain

import (
	residentDomain "portarius/internal/resident/domain"
	"time"

	"gorm.io/gorm"
)

type PackageStatus string

const (
	PackagePending   PackageStatus = "PENDENTE"
	PackageDelivered PackageStatus = "ENTREGUE"
	PackageLost      PackageStatus = "EXTRAVIADO"
)

// Package represents a package
// swagger:model
type Package struct {
	gorm.Model    `swaggerignore:"true"`
	Quantity      int                      `json:"quantity" gorm:"not null";default:1`
	ResidentID    *uint                    `json:"resident_id" gorm:"not null"`
	Resident      *residentDomain.Resident `json:"resident" gorm:"foreignKey:ResidentID" swaggerignore:"true"`
	Description   string                   `json:"description"`
	Status        PackageStatus            `json:"status" gorm:"not null;default:'PENDENTE'"`
	DeliveredToID *uint                    `json:"delivered_to_id"`
	DeliveredTo   *residentDomain.Resident `json:"delivered_to" gorm:"foreignKey:DeliveredToID" swaggerignore:"true"`
	ReceivedAt    time.Time                `json:"received_at"`
	DeliveredAt   time.Time                `json:"delivered_at"`
}
