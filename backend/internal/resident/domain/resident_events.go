package domain

import (
	"portarius/utils"

	"gorm.io/gorm"
)

func (r *Resident) BeforeSave(tx *gorm.DB) error {
	return r.Normalise()
}

func (r *Resident) BeforeCreate(tx *gorm.DB) error {
	return r.Normalise()
}

func (r *Resident) BeforeUpdate(tx *gorm.DB) error {
	return r.Normalise()
}

func (r *Resident) Normalise() error {
	r.Document = utils.KeepOnlyNumbers(r.Document)
	r.Phone = utils.KeepOnlyNumbers(r.Phone)
	r.Apartment = utils.KeepOnlyNumbers(r.Apartment)
	r.Block = utils.GetFirstLetter(r.Block)
	return nil
}
