package holyday

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Holyday struct {
	Date time.Time `json:"date"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

func IsHolyday(date time.Time) bool {
	holidays := getHolidaysMock(date.Year())

	dateStr := date.Format("2006-01-02")

	for _, holiday := range holidays {
		holidayStr := holiday.Date.Format("2006-01-02")
		if holidayStr == dateStr {
			return true
		}
	}

	return false
}

func getHolidays(year int) []Holyday {
	holidays, err := getHolidaysFromAPI(year)
	if err != nil {
		log.Printf("Erro ao obter feriados da API: %v", err)
		return []Holyday{}
	}

	holidays = append(holidays, Holyday{
		Date: time.Date(year, time.November, 14, 0, 0, 0, 0, time.UTC),
		Name: "Aniversário de Cascavel",
		Type: "local",
	})

	return holidays
}

func getHolidaysFromAPI(year int) ([]Holyday, error) {
	url := "https://brasilapi.com.br/api/feriados/v1/" + strconv.Itoa(year)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na requisição: status %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var holidays []Holyday
	err = json.Unmarshal(body, &holidays)
	if err != nil {
		return nil, err
	}

	return holidays, nil
}

func getHolidaysMock(year int) []Holyday {
	holidays := []Holyday{
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
