package domain

type InventoryService struct {
	repo IInventoryRepository
}

func NewInventoryService(repo IInventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) CreateInventory(i *Inventory) error {
	return s.repo.Create(i)
}

func (s *InventoryService) GetInventoryByID(id uint) (*Inventory, error) {
	return s.repo.GetByID(id)
}

func (s *InventoryService) UpdateInventory(i *Inventory) error {
	return s.repo.Update(i)
}

func (s *InventoryService) DeleteInventory(id uint) error {
	return s.repo.Delete(id)
}

func (s *InventoryService) GetAllInventory() ([]Inventory, error) {
	return s.repo.GetAll()
}
