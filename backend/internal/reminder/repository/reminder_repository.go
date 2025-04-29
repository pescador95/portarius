package repository

import (
	"portarius/internal/reminder/domain"
	"time"

	"gorm.io/gorm"
)

type reminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) domain.IReminderRepository {
	return &reminderRepository{db: db}
}

func (r *reminderRepository) GetAll() ([]domain.Reminder, error) {
	var reminders []domain.Reminder
	err := r.db.Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetByID(id uint) (*domain.Reminder, error) {
	var reminder domain.Reminder
	err := r.db.First(&reminder, id).Error
	return &reminder, err
}

func (r *reminderRepository) Create(reminder *domain.Reminder) error {
	return r.db.Create(reminder).Error
}

func (r *reminderRepository) Update(reminder *domain.Reminder) error {
	return r.db.Save(reminder).Error
}

func (r *reminderRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Reminder{}, id).Error
}

func (r *reminderRepository) GetByReservationID(reservationID uint) (*domain.Reminder, error) {
	var reminder domain.Reminder
	err := r.db.Where("reservation_id = ?", reservationID).First(&reminder).Error
	return &reminder, err
}

func (r *reminderRepository) GetByPackageID(packageID uint) (*domain.Reminder, error) {
	var reminder domain.Reminder
	err := r.db.Where("package_id = ?", packageID).First(&reminder).Error
	return &reminder, err
}

func (r *reminderRepository) GetByStatus(status string) ([]domain.Reminder, error) {
	var reminders []domain.Reminder
	err := r.db.Where("status = ?", status).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetByChannel(channel string) ([]domain.Reminder, error) {
	var reminders []domain.Reminder
	err := r.db.Where("channel = ?", channel).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetByRecipient(recipient string) ([]domain.Reminder, error) {
	var reminders []domain.Reminder
	err := r.db.Where("recipient = ?", recipient).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetByPendingStatus() ([]domain.Reminder, error) {
	var reminders []domain.Reminder
	err := r.db.Where("status IN ?", []domain.ReminderStatus{
		domain.ReminderStatusPending,
		domain.ReminderStatusFailed}).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetPendingRemindersFromReservations() ([]domain.Reminder, error) {
	var reminders []domain.Reminder
	err := r.db.Where("status IN ? AND reservation_id IS NOT NULL", []domain.ReminderStatus{
		domain.ReminderStatusPending,
		domain.ReminderStatusFailed}).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetPendingRemindersFromPackages() ([]domain.Reminder, error) {
	var reminders []domain.Reminder
	err := r.db.Where("status IN ? AND package_id IS NOT NULL", []domain.ReminderStatus{
		domain.ReminderStatusPending,
		domain.ReminderStatusFailed}).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetPendingRemindersFromReservationsForToday(now time.Time) ([]domain.Reminder, error) {
	var reminders []domain.Reminder

	startOfWindow := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	endOfWindow := startOfWindow.Add(time.Minute)

	err := r.db.
		Where("status IN ? AND reservation_id IS NOT NULL AND scheduled_at >= ? AND scheduled_at < ?",
			[]domain.ReminderStatus{
				domain.ReminderStatusPending,
				domain.ReminderStatusFailed,
			},
			startOfWindow,
			endOfWindow,
		).
		Find(&reminders).Error

	return reminders, err
}
