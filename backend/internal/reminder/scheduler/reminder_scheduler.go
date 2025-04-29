package scheduler

import (
	"fmt"
	"portarius/internal/reminder/domain"
	"time"

	cron "github.com/robfig/cron/v3"
)

type ReminderScheduler struct {
	repo domain.IReminderRepository
}

func NewReminderScheduler(repo domain.IReminderRepository) *ReminderScheduler {
	return &ReminderScheduler{repo: repo}
}

func (s *ReminderScheduler) Run() {

	c := cron.New()

	now := time.Now()

	noon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())

	if now.After(noon) {

		go s.StartReminderScheduler()

	}

	c.AddFunc("0 12 * * *", func() {
		s.StartReminderScheduler()
	})

	c.Start()
}

func (s *ReminderScheduler) StartReminderScheduler() {

	now := time.Now()
	windowStart := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())

	reminders, err := s.repo.GetPendingRemindersFromReservationsForToday(windowStart)

	if err != nil {
		fmt.Printf("[ReminderScheduler] Failed to get reminders: %v\n", err)
		return
	}

	for _, r := range reminders {
		r.PublishReservationReminder()
	}

}
