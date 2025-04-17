package resident

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type ResidentImportService struct {
	db *gorm.DB
}

func NewResidentImportService(db *gorm.DB) *ResidentImportService {
	return &ResidentImportService{db: db}
}

func (s *ResidentImportService) ImportResidentsFromCSV() error {

	etlPath := filepath.Join("..", "ETL", "resources", "resident_inventory")
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

func (s *ResidentImportService) processCSVFile(filePath string) error {
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

		if len(record) < 7 || record[3] == "" {
			continue
		}

		residentType := Tenant
		if strings.TrimSpace(strings.ToUpper(record[2])) == "PROPRIETARIO" {
			residentType = Owner
		}

		resident := Resident{
			Name:         strings.TrimSpace(record[3]),
			Document:     strings.TrimSpace(record[4]),
			Email:        strings.TrimSpace(record[6]),
			Phone:        strings.TrimSpace(record[5]),
			Apartment:    strings.TrimSpace(record[1]),
			Block:        strings.TrimSpace(record[0]),
			ResidentType: residentType,
		}

		log.Printf("Processing resident: %+v", resident)

		var existingResident Resident
		result := s.db.Where("document = ?", resident.Document).First(&existingResident)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return fmt.Errorf("error checking existing resident: %v", result.Error)
		}

		if result.Error == gorm.ErrRecordNotFound {
			if err := s.db.Create(&resident).Error; err != nil {
				log.Printf("Error creating resident: %v", err)
				return fmt.Errorf("error creating resident: %v", err)
			}
			log.Printf("Created new resident with document: %s", resident.Document)
		} else {

			if err := s.db.Model(&existingResident).Updates(resident).Error; err != nil {
				log.Printf("Error updating resident: %v", err)
				return fmt.Errorf("error updating resident: %v", err)
			}
			log.Printf("Updated existing resident with document: %s", resident.Document)
		}
		recordsProcessed++
	}

	log.Printf("Successfully processed %d records from file: %s", recordsProcessed, filePath)
	return nil
}
