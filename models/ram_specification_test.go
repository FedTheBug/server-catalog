package models

import (
	"testing"
)

func TestRamSpec_TableName(t *testing.T) {
	tests := []struct {
		name     string
		spec     RamSpec
		expected string
	}{
		{
			name:     "default table name",
			spec:     RamSpec{},
			expected: "ram_spec",
		},
		{
			name: "table name with populated struct",
			spec: RamSpec{
				ID:   1,
				Type: "DDR3",
			},
			expected: "ram_spec",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.TableName(); got != tt.expected {
				t.Errorf("RamSpec.TableName() = %v, want %v", got, tt.expected)
			}
		})
	}
}
