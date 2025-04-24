package domain

import (
	"time"
)

type Holyday struct {
	Date time.Time `json:"date"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}
