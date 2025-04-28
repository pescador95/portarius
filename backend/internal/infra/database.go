package infra

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	inventoryDomain "portarius/internal/inventory/domain"

	packageDomain "portarius/internal/package/domain"

	reservationDomain "portarius/internal/reservation/domain"

	residentDomain "portarius/internal/resident/domain"

	userDomain "portarius/internal/user/domain"

	reminderDomain "portarius/internal/reminder/domain"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}
	return db, nil
}

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&inventoryDomain.Inventory{},
		&packageDomain.Package{},
		&residentDomain.Resident{},
		&reservationDomain.Reservation{},
		&userDomain.User{},
		&reminderDomain.Reminder{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
}
