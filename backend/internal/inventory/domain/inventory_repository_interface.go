package domain

type IInventoryRepository interface {
	GetAll() ([]Inventory, error)
	GetByID(id uint) (*Inventory, error)
	Create(inventory *Inventory) error
	Update(inventory *Inventory) error
	Delete(id uint) error
}
