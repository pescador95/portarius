package domain

type IPackageRepository interface {
	GetAll(page, pageSize int) ([]Package, error)
	GetByID(id uint) (*Package, error)
	Create(pkg *Package) error
	Update(pkg *Package) error
	Delete(id uint) error
	MarkAsDelivered(id uint) error
	MarkAsLost(id uint) error
}
