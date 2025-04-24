package handler

import (
	"log"
	"portarius/internal/holyday/domain"
	"portarius/internal/holyday/service"
	"time"
)

func IsHolyday(date time.Time) bool {
	holidays := GetHolidaysMock(date.Year())

	dateStr := date.Format("2006-01-02")

	for _, holiday := range holidays {
		holidayStr := holiday.Date.Format("2006-01-02")
		if holidayStr == dateStr {
			return true
		}
	}

	return false
}

func GetHolidays(year int) []domain.Holyday {
	holidays, err := service.GetHolidaysFromAPI(year)
	if err != nil {
		log.Printf("Erro ao obter feriados da API: %v", err)
		return []domain.Holyday{}
	}

	holidays = append(holidays, domain.Holyday{
		Date: time.Date(year, time.November, 14, 0, 0, 0, 0, time.UTC),
		Name: "Aniversário de Cascavel",
		Type: "local",
	})

	return holidays
}

func GetHolidaysMock(year int) []domain.Holyday {
	holidays := []domain.Holyday{
		{
			Date: time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC),
			Name: "Confraternização mundial",
			Type: "national",
		},
		{
			Date: time.Date(year, time.March, 4, 0, 0, 0, 0, time.UTC),
			Name: "Carnaval",
			Type: "national",
		},
		{
			Date: time.Date(year, time.April, 18, 0, 0, 0, 0, time.UTC),
			Name: "Sexta-feira Santa",
			Type: "national",
		},
		{
			Date: time.Date(year, time.April, 20, 0, 0, 0, 0, time.UTC),
			Name: "Páscoa",
			Type: "national",
		},
		{
			Date: time.Date(year, time.April, 21, 0, 0, 0, 0, time.UTC),
			Name: "Tiradentes",
			Type: "national",
		},
		{
			Date: time.Date(year, time.May, 1, 0, 0, 0, 0, time.UTC),
			Name: "Dia do trabalho",
			Type: "national",
		},
		{
			Date: time.Date(year, time.June, 19, 0, 0, 0, 0, time.UTC),
			Name: "Corpus Christi",
			Type: "national",
		},
		{
			Date: time.Date(year, time.September, 7, 0, 0, 0, 0, time.UTC),
			Name: "Independência do Brasil",
			Type: "national",
		},
		{
			Date: time.Date(year, time.October, 12, 0, 0, 0, 0, time.UTC),
			Name: "Nossa Senhora Aparecida",
			Type: "national",
		},
		{
			Date: time.Date(year, time.November, 2, 0, 0, 0, 0, time.UTC),
			Name: "Finados",
			Type: "national",
		},
		{
			Date: time.Date(year, time.November, 15, 0, 0, 0, 0, time.UTC),
			Name: "Proclamação da República",
			Type: "national",
		},
		{
			Date: time.Date(year, time.November, 20, 0, 0, 0, 0, time.UTC),
			Name: "Dia da consciência negra",
			Type: "national",
		},
		{
			Date: time.Date(year, time.December, 25, 0, 0, 0, 0, time.UTC),
			Name: "Natal",
			Type: "national",
		},
		{
			Date: time.Date(year, time.November, 14, 0, 0, 0, 0, time.UTC),
			Name: "Aniversário de Cascavel",
			Type: "local",
		},
	}

	return holidays
}
