package reservation

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"strings"
	"time"

	"gorm.io/gorm"

	domainReservation "portarius/internal/reservation/domain"
)

var monthMap = map[string]int{
	"jan": 1,
	"fev": 2,
	"mar": 3,
	"abr": 4,
	"mai": 5,
	"jun": 6,
	"jul": 7,
	"ago": 8,
	"set": 9,
	"out": 10,
	"nov": 11,
	"dez": 12,
}

var residentMap = map[string]uint{
	"A01": 1,
	"A02": 3,
	"A03": 5,
	"A04": 7,
	"A05": 9,
	"A06": 10,
	"A07": 12,
	"A08": 15,
	"A09": 16,
	"A10": 17,
	"A11": 18,
	"A13": 20,
	"A14": 22,
	"A15": 23,
	"A16": 26,
	"A17": 28,
	"A18": 30,
	"A20": 31,
	"B01": 32,
	"B03": 33,
	"B05": 34,
	"B08": 36,
	"B09": 38,
	"B10": 39,
	"B11": 40,
	"B12": 42,
	"B13": 43,
	"B14": 45,
	"B15": 47,
	"B16": 48,
	"B17": 49,
	"B18": 50,
	"B19": 51,
	"B20": 53,
	"C01": 54,
	"C02": 55,
	"C04": 56,
	"C07": 58,
	"C08": 60,
	"C10": 62,
	"C11": 63,
	"C12": 64,
	"C13": 66,
	"C15": 68,
	"C16": 71,
	"C17": 72,
	"C19": 74,
	"C20": 76,
	"D01": 78,
	"D02": 79,
	"D03": 80,
	"D04": 81,
	"D05": 82,
	"D06": 84,
	"D08": 85,
	"D09": 87,
	"D10": 88,
	"D11": 90,
	"D13": 92,
	"D15": 93,
	"D16": 94,
	"D17": 95,
	"D18": 96,
	"D20": 98,
	"E02": 99,
	"E04": 101,
	"E05": 102,
	"E06": 6,
	"E07": 103,
	"E08": 105,
	"E09": 106,
	"E10": 108,
	"E11": 110,
	"E12": 111,
	"E13": 112,
	"E15": 113,
	"E16": 114,
	"E17": 115,
	"E18": 116,
	"E19": 117,
	"E20": 118,
}

type ReservationImportService struct {
	db *gorm.DB
}

func NewReservationImportService(db *gorm.DB) *ReservationImportService {
	return &ReservationImportService{db: db}
}

func normalizeUnit(unit string) string {
	unit = strings.TrimSpace(strings.ToUpper(unit))
	if len(unit) < 2 {
		return ""
	}

	block := unit[0:1]
	number := strings.TrimLeft(unit[1:], "0")
	if number == "" {
		number = "0"
	}
	return fmt.Sprintf("%s%02s", block, number)
}

func (s *ReservationImportService) ImportSalonReservationsFromCSV() error {

	etlPath := filepath.Join("..", "resources", "ETL", "salon")
	absPath, err := filepath.Abs(etlPath)
	if err != nil {
		return fmt.Errorf("erro ao obter caminho absoluto: %v", err)
	}

	files, err := filepath.Glob(filepath.Join(absPath, "*.csv"))
	if err != nil {
		return fmt.Errorf("erro ao buscar arquivos CSV: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("nenhum arquivo CSV encontrado em: %s", absPath)
	}

	for _, file := range files {
		if err := s.processCSVFile(file); err != nil {
			return fmt.Errorf("erro ao processar arquivo %s: %v", file, err)
		}
	}

	return nil
}

func (s *ReservationImportService) processCSVFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("erro ao ler cabeçalho: %v", err)
	}

	var processados, ignorados int
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("erro ao ler registro: %v", err)
		}

		if len(record) < 4 {
			ignorados++
			continue
		}

		data := strings.TrimSpace(record[0])
		salao := strings.TrimSpace(record[1])
		unidade := normalizeUnit(record[2])
		formaPagamento := strings.TrimSpace(record[3])

		if unidade == "" {
			log.Printf("Unidade inválida: %s", record[2])
			ignorados++
			continue
		}

		if strings.Contains("FGH", unidade[0:1]) {
			log.Printf("Bloco não cadastrado: %s", unidade)
			ignorados++
			continue
		}

		residentID, ok := residentMap[unidade]
		if !ok {
			log.Printf("Residente não encontrado para unidade: %s", unidade)
			ignorados++
			continue
		}

		var date time.Time

		date, err = time.Parse(time.RFC3339, data)
		if err != nil {

			date, err = time.Parse("2006-01-02", data)
			if err != nil {
				log.Printf("Erro ao converter data %s: %v", data, err)
				ignorados++
				continue
			}
		}

		date, err = time.Parse("2006-01-02", data[:10])

		fmt.Println(data[:10])

		reserva := &domainReservation.Reservation{
			ResidentID:    &residentID,
			Space:         domainReservation.Salon1,
			StartTime:     time.Date(date.Year(), date.Month(), date.Day(), 8, 0, 0, 0, time.UTC),
			EndTime:       time.Date(date.Year(), date.Month(), date.Day(), 20, 0, 0, 0, time.UTC),
			Status:        domainReservation.StatusConfirmed,
			PaymentStatus: domainReservation.PaymentPaid,
			PaymentMethod: domainReservation.PaymentMethodBoleto,
		}

		if strings.Contains(strings.ToUpper(salao), "2") {
			reserva.Space = domainReservation.Salon2
		}

		if strings.Contains(strings.ToUpper(formaPagamento), "PIX") {
			reserva.PaymentMethod = domainReservation.PaymentMethodPix
		}

		paymentDate := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
		reserva.PaymentDate = &paymentDate

		var existingReservation domainReservation.Reservation
		result := s.db.Where("resident_id = ? AND start_time = ? AND end_time = ?",
			reserva.ResidentID, reserva.StartTime, reserva.EndTime).First(&existingReservation)

		if result.Error == nil {

			if err := s.db.Model(&existingReservation).Updates(reserva).Error; err != nil {
				return fmt.Errorf("erro ao atualizar reserva: %v", err)
			}
		} else if result.Error == gorm.ErrRecordNotFound {

			if err := s.db.Create(reserva).Error; err != nil {
				return fmt.Errorf("erro ao criar reserva: %v", err)
			}
		} else {
			return fmt.Errorf("erro ao verificar reserva existente: %v", result.Error)
		}

		processados++
	}

	log.Printf("Arquivo %s processado. Registros processados: %d, ignorados: %d",
		filepath.Base(filePath), processados, ignorados)
	return nil
}
