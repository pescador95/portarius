package domain

type ResidentService struct {
	repo IResidentRepository
}

func NewResidentService(repo IResidentRepository) *ResidentService {
	return &ResidentService{repo: repo}
}

func (s *ResidentService) CreateResident(r *Resident) error {
	return s.repo.Create(r)
}

func (s *ResidentService) GetResidentByID(id uint) (*Resident, error) {
	return s.repo.GetByID(id)
}

func (s *ResidentService) UpdateResident(r *Resident) error {
	return s.repo.Update(r)
}

func (s *ResidentService) DeleteResident(id uint) error {
	return s.repo.Delete(id)
}

func (s *ResidentService) GetAllResidents() ([]Resident, error) {
	return s.repo.GetAll()
}
