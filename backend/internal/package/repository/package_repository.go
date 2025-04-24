package repository

import (
	"portarius/internal/package/domain"

	"gorm.io/gorm"
)

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) domain.IPackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) GetAll() ([]domain.Package, error) {
	var packages []domain.Package
	err := r.db.Find(&packages).Error
	return packages, err
}

func (r *packageRepository) GetByID(id uint) (*domain.Package, error) {
	var pkg domain.Package
	err := r.db.First(&pkg, id).Error
	return &pkg, err
}

func (r *packageRepository) Create(pkg *domain.Package) error {
	return r.db.Create(pkg).Error
}

func (r *packageRepository) Update(pkg *domain.Package) error {
	return r.db.Save(pkg).Error
}

func (r *packageRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Package{}, id).Error
}

func (r *packageRepository) MarkAsDelivered(id uint) error {
	pkg, err := r.GetByID(id)
	if err != nil {
		return err
	}
	pkg.Status = domain.PackageDelivered
	return r.Update(pkg)
}

func (r *packageRepository) MarkAsLost(id uint) error {
	pkg, err := r.GetByID(id)
	if err != nil {
		return err
	}
	pkg.Status = domain.PackageLost
	return r.Update(pkg)
}
