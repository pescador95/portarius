package domain_test

import (
	residentDomain "portarius/internal/resident/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResidentNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    *residentDomain.Resident
		expected *residentDomain.Resident
	}{
		{
			name: "should normalize document, phone and block",
			input: &residentDomain.Resident{
				Document:  "123.456.789-00",
				Phone:     "(11) 99999-9999",
				Apartment: "42",
				Block:     "B",
			},
			expected: &residentDomain.Resident{
				Document:  "12345678900",
				Phone:     "11999999999",
				Apartment: "42",
				Block:     "B",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Normalise()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Document, tt.input.Document)
			assert.Equal(t, tt.expected.Phone, tt.input.Phone)
			assert.Equal(t, tt.expected.Apartment, tt.input.Apartment)
			assert.Equal(t, tt.expected.Block, tt.input.Block)
		})
	}
}
