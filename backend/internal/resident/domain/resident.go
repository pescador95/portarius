package domain

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
	Name         string       `json:"name" gorm:"size:100`
	Document     string       `json:"document" gorm:"size:20`
	Email        string       `json:"email" gorm:"size:100";check:email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'`
	Phone        string       `json:"phone" gorm:"size:15"`
	Apartment    string       `json:"apartment" gorm:"size:2;not null"`
	Block        string       `json:"block" gorm:"size:1;not null"`
	ResidentType ResidentType `json:"resident_type" gorm:"not null;default:'INQUILINO'"`
}
