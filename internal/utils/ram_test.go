package utils

import (
	"reflect"
	"testing"
)

func TestParseRAMValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "single value",
			input:    "16GB",
			expected: []int{16},
		},
		{
			name:     "multiple values",
			input:    "16GB,32GB,64GB",
			expected: []int{16, 32, 64},
		},
		{
			name:     "values with spaces",
			input:    "16GB, 32GB, 64GB",
			expected: []int{16, 32, 64},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []int{},
		},
		{
			name:     "invalid values",
			input:    "abc,def",
			expected: []int{},
		},
		{
			name:     "mixed valid and invalid",
			input:    "16GB,abc,32GB",
			expected: []int{16, 32},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseRAMValues(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseRAMValues = %v, want %v", got, tt.expected)
			}
		})
	}
}
