package models

import (
	"testing"
)

func TestServerCatalog_TableName(t *testing.T) {
	tests := []struct {
		name     string
		catalog  ServerCatalog
		expected string
	}{
		{
			name:     "default table name",
			catalog:  ServerCatalog{},
			expected: "server_catalog",
		},
		{
			name: "table name with populated struct",
			catalog: ServerCatalog{
				Model:    "HP DL120G7Intel G850",
				RamSize:  4,
				RamType:  1,
				HDDSize:  1000,
				HDDCount: 4,
				HDDType:  2,
				Location: "AmsterdamAMS-01",
				Price:    39.99,
				Currency: 1,
			},
			expected: "server_catalog",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.catalog.TableName(); got != tt.expected {
				t.Errorf("ServerCatalog.TableName() = %v, want %v", got, tt.expected)
			}
		})
	}
}
