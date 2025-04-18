package domain

type IResidentRepository interface {
	GetAll() ([]Resident, error)
	GetByID(id uint) (*Resident, error)
	Create(resident *Resident) error
	Update(resident *Resident) error
	Delete(id uint) error
}
