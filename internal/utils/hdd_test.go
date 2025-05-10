package utils

import (
	"testing"
)

func TestParseStorageToGB(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		expectError bool
	}{
		{
			name:        "valid GB",
			input:       "500GB",
			expected:    500,
			expectError: false,
		},
		{
			name:        "valid TB",
			input:       "2TB",
			expected:    2048,
			expectError: false,
		},
		{
			name:        "zero value",
			input:       "0",
			expected:    0,
			expectError: false,
		},
		{
			name:        "invalid format",
			input:       "invalid",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid GB number",
			input:       "abcGB",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid TB number",
			input:       "xyzTB",
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
			got, err := ParseStorageToGB(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("ParseStorageToGB() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseStorageToGB() = %v, want %v", got, tt.expected)
			}
		})
	}
}
