package inventory

import (
	"portarius/resident"
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

type Inventory struct {
	gorm.Model
	Name          string            `json:"name" gorm:"not null"`
	Description   string            `json:"description"`
	Quantity      int               `json:"quantity" gorm:"not null"`
	OwnerID       uint              `json:"owner_id" gorm:"not null"`
	Owner         resident.Resident `json:"owner" gorm:"foreignKey:OwnerID"`
	LastUpdated   time.Time         `json:"last_updated"`
	InventoryType InventoryType     `json:"inventory_type"`
}
