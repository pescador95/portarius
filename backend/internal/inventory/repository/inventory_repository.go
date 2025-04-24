package repository

import (
	"portarius/internal/inventory/domain"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) domain.IInventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetAll() ([]domain.Inventory, error) {
	var inventories []domain.Inventory
	err := r.db.Find(&inventories).Error
	return inventories, err
}

func (r *inventoryRepository) GetByID(id uint) (*domain.Inventory, error) {
	var inventory domain.Inventory
	err := r.db.First(&inventory, id).Error
	return &inventory, err
}

func (r *inventoryRepository) Create(inventory *domain.Inventory) error {
	return r.db.Create(inventory).Error
}

func (r *inventoryRepository) Update(inventory *domain.Inventory) error {
	return r.db.Save(inventory).Error
}

func (r *inventoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Inventory{}, id).Error
}
