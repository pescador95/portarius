package resident

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"portarius/internal/resident/domain"
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
		log.Printf("Processing file: %s", file)
		if err := s.processCSVFile(file); err != nil {
			log.Printf("Error processing file %s: %v", file, err)
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
			log.Printf("Error reading record: %v", err)
			continue
		}

		resident, err := parseResident(record)
		if err != nil {
			log.Printf("Skipping invalid record: %v", err)
			continue
		}

		if err := s.upsertResident(resident); err != nil {
			log.Printf("Error upserting resident: %v", err)
			continue
		}

		recordsProcessed++
	}

	log.Printf("Successfully processed %d records from file: %s", recordsProcessed, filePath)
	return nil
}

func parseResident(record []string) (domain.Resident, error) {
	if len(record) < 7 || strings.TrimSpace(record[3]) == "" {
		return domain.Resident{}, fmt.Errorf("invalid record or missing name")
	}

	residentType := domain.Tenant
	if strings.TrimSpace(strings.ToUpper(record[2])) == "PROPRIETARIO" {
		residentType = domain.Owner
	}

	resident := domain.Resident{
		Name:         strings.TrimSpace(record[3]),
		Document:     strings.TrimSpace(record[4]),
		Email:        strings.TrimSpace(record[6]),
		Phone:        strings.TrimSpace(record[5]),
		Apartment:    strings.TrimSpace(record[1]),
		Block:        strings.TrimSpace(record[0]),
		ResidentType: residentType,
	}

	resident.Normalise()
	return resident, nil
}

func (s *ResidentImportService) upsertResident(resident domain.Resident) error {
	var existing domain.Resident
	result := s.db.Where("document = ?", resident.Document).First(&existing)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if result.Error == gorm.ErrRecordNotFound {
		if err := s.db.Create(&resident).Error; err != nil {
			return fmt.Errorf("error creating resident: %v", err)
		}
		log.Printf("Created new resident: %s", resident.Document)
	} else {
		if err := s.db.Model(&existing).Updates(resident).Error; err != nil {
			return fmt.Errorf("error updating resident: %v", err)
		}
		log.Printf("Updated resident: %s", resident.Document)
	}

	return nil
}
