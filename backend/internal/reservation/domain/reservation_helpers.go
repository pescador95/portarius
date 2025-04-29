package domain

import "strings"

func (r *Reservation) GetLastCharFromSalon() string {
	parts := strings.Split(string(r.Space), "_")
	if len(parts) > 1 {
		suffix := parts[len(parts)-1]
		if len(suffix) > 0 {
			return string(suffix[len(suffix)-1])
		}
	}
	return ""
}
