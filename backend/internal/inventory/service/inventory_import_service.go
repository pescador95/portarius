package inventory

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"portarius/internal/inventory/domain"
	residentDomain "portarius/internal/resident/domain"
	"portarius/internal/utils"

	"gorm.io/gorm"
)

type InventoryImportService struct {
	db *gorm.DB
}

func NewInventoryImportService(db *gorm.DB) *InventoryImportService {
	return &InventoryImportService{db: db}
}

func (s *InventoryImportService) ImportPetsFromCSV() error {

	etlPath := filepath.Join("..", "resources", "ETL", "resident_inventory")
	absPath, err := filepath.Abs(etlPath)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}

	log.Printf("Looking for CSV files in directory: %s", absPath)

	files, err := filepath.Glob(filepath.Join(absPath, "*.csv"))
	if err != nil {
		return fmt.Errorf("error finding CSV files: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no CSV files found in directory: %s", absPath)
	}

	log.Printf("Found %d CSV files to process", len(files))
	for _, file := range files {
		log.Printf("Found file: %s", file)
	}

	for _, file := range files {
		log.Printf("Processing file: %s", file)
		if err := s.processCSVFile(file); err != nil {
			return fmt.Errorf("error processing file %s: %v", file, err)
		}
	}

	return nil
}

func (s *InventoryImportService) processCSVFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("error reading header: %v", err)
	}

	recordsProcessed := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error reading record: %v", err)
		}

		if len(record) < 9 || record[7] == "" {
			continue
		}

		var resident residentDomain.Resident
		if err := s.db.Where("document = ?", strings.TrimSpace(utils.KeepOnlyNumbers(record[4]))).First(&resident).Error; err != nil {
			log.Printf("Resident not found with document: %s", utils.KeepOnlyNumbers(record[4]))
			continue
		}

		inventory := domain.Inventory{
			Name:          strings.TrimSpace(record[7]),
			Description:   strings.TrimSpace(record[8]),
			Quantity:      1,
			OwnerID:       &resident.ID,
			InventoryType: domain.InventoryTypePet,
		}

		log.Printf("Processing pet: %+v", inventory)

		var existingInventory domain.Inventory
		result := s.db.Where("owner_id = ? AND name = ? AND inventory_type = ?", resident.ID, inventory.Name, domain.InventoryTypePet).First(&existingInventory)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return fmt.Errorf("error checking existing pet: %v", result.Error)
		}

		if result.Error == gorm.ErrRecordNotFound {
			if err := s.db.Create(&inventory).Error; err != nil {
				log.Printf("Error creating pet: %v", err)
				return fmt.Errorf("error creating pet: %v", err)
			}
			log.Printf("Created new pet for resident ID: %d", resident.ID)
		} else {

			if err := s.db.Model(&existingInventory).Updates(inventory).Error; err != nil {
				log.Printf("Error updating pet: %v", err)
				return fmt.Errorf("error updating pet: %v", err)
			}
			log.Printf("Updated existing pet for resident ID: %d", resident.ID)
		}
		recordsProcessed++
	}

	log.Printf("Successfully processed %d pet records from file: %s", recordsProcessed, filePath)
	return nil
}
