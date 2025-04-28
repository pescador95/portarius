package service

import (
	pkgDomain "portarius/internal/package/domain"
	reminderDomain "portarius/internal/reminder/domain"
)

type PackageReminderService struct {
	repo reminderDomain.IReminderRepository
}

func NewPackageReminderService(repo reminderDomain.IReminderRepository) *PackageReminderService {
	return &PackageReminderService{repo: repo}
}

func (s *PackageReminderService) CreateReminderForPackage(pkg *pkgDomain.Package) error {

	reminder := reminderDomain.Reminder{
		Recipient:     pkg.Resident.Phone,
		PackageID:     &pkg.ID,
		ReservationID: nil,
		Channel:       reminderDomain.ReminderChannelWhatsApp,
		Status:        reminderDomain.ReminderStatusPending,
	}
	return s.repo.Create(&reminder)
}
