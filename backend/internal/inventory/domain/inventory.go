package domain

import (
	"portarius/internal/resident/domain"
	"time"

	"gorm.io/gorm"
)

type InventoryType string

const (
	InventoryTypeCar     InventoryType = "CARRO"
	InventoryTypeBike    InventoryType = "MOTO"
	InventoryTypeBicycle InventoryType = "BICICLETA"
	InventoryTypeScooter InventoryType = "SCOOTER"
	InventoryTypePet     InventoryType = "PET"
)

// Inventory represents an inventory item
// swagger:model
type Inventory struct {
	gorm.Model    `swaggerignore:"true"`
	Name          string           `json:"name" gorm:"not null"`
	Description   string           `json:"description"`
	Quantity      int              `json:"quantity" gorm:"not null"`
	OwnerID       *uint            `json:"owner_id" gorm:"not null"`
	Owner         *domain.Resident `json:"owner" gorm:"foreignKey:OwnerID" swaggerignore:"true"`
	LastUpdated   time.Time        `json:"last_updated"`
	InventoryType InventoryType    `json:"inventory_type"`
}
