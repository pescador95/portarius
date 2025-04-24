package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"portarius/internal/holyday/domain"
	"strconv"
)

func GetHolidaysFromAPI(year int) ([]domain.Holyday, error) {
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

	var holidays []domain.Holyday
	err = json.Unmarshal(body, &holidays)
	if err != nil {
		return nil, err
	}

	return holidays, nil
}
