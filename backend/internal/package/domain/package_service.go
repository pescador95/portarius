package domain

type PackageService struct {
	repo IPackageRepository
}

func NewPackageService(repo IPackageRepository) *PackageService {
	return &PackageService{repo: repo}
}

func (s *PackageService) CreatePackage(p *Package) error {
	return s.repo.Create(p)
}

func (s *PackageService) GetPackageByID(id uint) (*Package, error) {
	return s.repo.GetByID(id)
}

func (s *PackageService) UpdatePackage(p *Package) error {
	return s.repo.Update(p)
}

func (s *PackageService) DeletePackage(id uint) error {
	return s.repo.Delete(id)
}

func (s *PackageService) GetAllPackages() ([]Package, error) {
	return s.repo.GetAll()
}

func (s *PackageService) MarkAsDelivered(id uint) error {
	return s.repo.MarkAsDelivered(id)
}

func (s *PackageService) MarkAsLost(id uint) error {
	return s.repo.MarkAsLost(id)
}
