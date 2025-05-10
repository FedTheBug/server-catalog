package models

import (
	"testing"
)

func TestHDDSpec_TableName(t *testing.T) {
	tests := []struct {
		name     string
		spec     HDDSpec
		expected string
	}{
		{
			name:     "default table name",
			spec:     HDDSpec{},
			expected: "hdd_spec",
		},
		{
			name: "table name with populated struct",
			spec: HDDSpec{
				ID:   1,
				Type: "SATA2",
			},
			expected: "hdd_spec",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.TableName(); got != tt.expected {
				t.Errorf("HDDSpec.TableName() = %v, want %v", got, tt.expected)
			}
		})
	}
}
