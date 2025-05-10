package utils

import (
	"strings"
	"testing"
)

func TestGetHDDTypeID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		expectError bool
	}{
		{
			name:        "valid SATA2",
			input:       "SATA2",
			expected:    HDDTypeSATA2,
			expectError: false,
		},
		{
			name:        "valid SAS",
			input:       "SAS",
			expected:    HDDTypeSAS,
			expectError: false,
		},
		{
			name:        "valid SSD",
			input:       "SSD",
			expected:    HDDTypeSSD,
			expectError: false,
		},
		{
			name:        "case insensitive SATA2",
			input:       "sata2",
			expected:    HDDTypeSATA2,
			expectError: false,
		},
		{
			name:        "invalid type",
			input:       "INVALID",
			expected:    0,
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHDDTypeID(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("GetHDDTypeID() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if got != tt.expected {
				t.Errorf("GetHDDTypeID() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetRAMTypeID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		expectError bool
	}{
		{
			name:        "valid DDR3",
			input:       "DDR3",
			expected:    RAMTypeDDR3,
			expectError: false,
		},
		{
			name:        "valid DDR4",
			input:       "DDR4",
			expected:    RAMTypeDDR4,
			expectError: false,
		},
		{
			name:        "case insensitive DDR3",
			input:       "ddr3",
			expected:    RAMTypeDDR3,
			expectError: false,
		},
		{
			name:        "invalid type",
			input:       "INVALID",
			expected:    0,
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRAMTypeID(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("GetRAMTypeID() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if got != tt.expected {
				t.Errorf("GetRAMTypeID() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConvertToGB(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		unit     string
		expected int
	}{
		{
			name:     "convert TB to GB",
			size:     2,
			unit:     "TB",
			expected: 2048,
		},
		{
			name:     "convert GB to GB",
			size:     500,
			unit:     "GB",
			expected: 500,
		},
		{
			name:     "case insensitive TB",
			size:     2,
			unit:     "tb",
			expected: 2048,
		},
		{
			name:     "unknown unit",
			size:     100,
			unit:     "MB",
			expected: 100,
		},
		{
			name:     "zero size",
			size:     0,
			unit:     "TB",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToGB(tt.size, tt.unit)
			if got != tt.expected {
				t.Errorf("ConvertToGB() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetCurrencyID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		expectError bool
	}{
		{
			name:        "valid USD",
			input:       CurrencySymbolUSD,
			expected:    CurrencyUSD,
			expectError: false,
		},
		{
			name:        "valid Euro",
			input:       CurrencySymbolEuro,
			expected:    CurrencyEuro,
			expectError: false,
		},
		{
			name:        "valid SGD",
			input:       CurrencySymbolSGD,
			expected:    CurrencySGD,
			expectError: false,
		},
		{
			name:        "with spaces",
			input:       " " + CurrencySymbolUSD + " ",
			expected:    CurrencyUSD,
			expectError: false,
		},
		{
			name:        "invalid symbol",
			input:       "¥",
			expected:    0,
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCurrencyID(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("GetCurrencyID() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if got != tt.expected {
				t.Errorf("GetCurrencyID() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		target      interface{}
		expectError bool
	}{
		{
			name:  "valid JSON object",
			input: `{"name": "test", "value": 123}`,
			target: &struct {
				Name  string `json:"name"`
				Value int    `json:"value"`
			}{},
			expectError: false,
		},
		{
			name:        "valid JSON array",
			input:       `[1, 2, 3]`,
			target:      &[]int{},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			input:       `{"name": "test", "value": 123`,
			target:      &struct{}{},
			expectError: true,
		},
		{
			name:        "empty JSON",
			input:       "",
			target:      &struct{}{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseJSON(strings.NewReader(tt.input), tt.target)
			if (err != nil) != tt.expectError {
				t.Errorf("ParseJSON() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
