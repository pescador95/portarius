package domain

type IReminderRepository interface {
	Create(pkg *Reminder) error
	Update(pkg *Reminder) error
	Delete(id uint) error
	GetAll() ([]Reminder, error)
	GetByID(id uint) (*Reminder, error)
	GetByReservationID(reservationID uint) (*Reminder, error)
	GetByPackageID(packageID uint) (*Reminder, error)
	GetByStatus(status string) ([]Reminder, error)
	GetByChannel(channel string) ([]Reminder, error)
	GetByRecipient(recipient string) ([]Reminder, error)
	GetByPendingStatus() ([]Reminder, error)
}
