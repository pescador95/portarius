package resident

import (
	"gorm.io/gorm"
)

type ResidentType string

const (
	Tenant      ResidentType = "INQUILINO"
	Owner       ResidentType = "PROPRIETARIO"
	Krum        ResidentType = "KRUM"
	NotResident ResidentType = "NAO_RESIDENTE"
)

type Resident struct {
	gorm.Model
	Name         string       `json:"name"`
	Document     string       `json:"document"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	Apartment    string       `json:"apartment" gorm:"not null"`
	Block        string       `json:"block" gorm:"not null"`
	ResidentType ResidentType `json:"resident_type" gorm:"not null;default:'INQUILINO'"`
}
