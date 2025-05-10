package transformer

import (
	"github.com/server-catalog/internal/dto"
	"github.com/server-catalog/internal/utils"
	"github.com/server-catalog/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransformServerList(t *testing.T) {
	tests := []struct {
		name     string
		input    []models.ServerCatalog
		expected []dto.ListServerResp
	}{
		{
			name: "transform DDR3 server with GB storage",
			input: []models.ServerCatalog{
				{
					Model:    "Dell R210",
					RamSize:  16,
					RamType:  utils.RAMTypeDDR3,
					HDDSize:  500,
					HDDCount: 2,
					HDDType:  utils.HDDTypeSATA2,
					Location: "AmsterdamAMS-01",
					Price:    99.99,
					Currency: 2,
				},
			},
			expected: []dto.ListServerResp{
				{
					Model:    "Dell R210",
					Ram:      "16GBDDR3",
					HDD:      "2x500GBSATA2",
					Location: "AmsterdamAMS-01",
					Price:    "€99.99",
				},
			},
		},
		{
			name: "transform DDR4 server with TB storage",
			input: []models.ServerCatalog{
				{
					Model:    "Dell R730XD",
					RamSize:  128,
					RamType:  utils.RAMTypeDDR4,
					HDDSize:  2048, // 2TB
					HDDCount: 4,
					HDDType:  utils.HDDTypeSSD,
					Location: "SingaporeSIN-11",
					Price:    565.99,
					Currency: 3,
				},
			},
			expected: []dto.ListServerResp{
				{
					Model:    "Dell R730XD",
					Ram:      "128GBDDR4",
					HDD:      "4x2TBSSD",
					Location: "SingaporeSIN-11",
					Price:    "S$565.99",
				},
			},
		},
		{
			name: "transform server with USD currency",
			input: []models.ServerCatalog{
				{
					Model:    "HP DL380",
					RamSize:  64,
					RamType:  utils.RAMTypeDDR3,
					HDDSize:  1024, // 1TB
					HDDCount: 8,
					HDDType:  utils.HDDTypeSAS,
					Location: "Washington D.C.WDC-01",
					Price:    199.99,
					Currency: 1,
				},
			},
			expected: []dto.ListServerResp{
				{
					Model:    "HP DL380",
					Ram:      "64GBDDR3",
					HDD:      "8x1TBSAS",
					Location: "Washington D.C.WDC-01",
					Price:    "$199.99",
				},
			},
		},
		{
			name: "transform multiple servers",
			input: []models.ServerCatalog{
				{
					Model:    "Dell R730XD2x Intel Xeon E5-2620v4",
					RamSize:  32,
					RamType:  utils.RAMTypeDDR4,
					HDDSize:  500,
					HDDCount: 2,
					HDDType:  utils.HDDTypeSATA2,
					Location: "AmsterdamAMS-01",
					Price:    149.99,
					Currency: 2,
				},
				{
					Model:    "Dell R730XD2x Intel Xeon E5-2620v4",
					RamSize:  64,
					RamType:  utils.RAMTypeDDR3,
					HDDSize:  2048,
					HDDCount: 4,
					HDDType:  utils.HDDTypeSSD,
					Location: "SingaporeSIN-11",
					Price:    299.99,
					Currency: 3,
				},
			},
			expected: []dto.ListServerResp{
				{
					Model:    "Dell R730XD2x Intel Xeon E5-2620v4",
					Ram:      "32GBDDR4",
					HDD:      "2x500GBSATA2",
					Location: "AmsterdamAMS-01",
					Price:    "€149.99",
				},
				{
					Model:    "Dell R730XD2x Intel Xeon E5-2620v4",
					Ram:      "64GBDDR3",
					HDD:      "4x2TBSSD",
					Location: "SingaporeSIN-11",
					Price:    "S$299.99",
				},
			},
		},
		{
			name: "transform server with unknown RAM type",
			input: []models.ServerCatalog{
				{
					Model:    "Dell R730XD2x Intel Xeon E5-2650v3",
					RamSize:  16,
					RamType:  999, // Unknown type
					HDDSize:  500,
					HDDCount: 2,
					HDDType:  utils.HDDTypeSATA2,
					Location: "AmsterdamAMS-01",
					Price:    99.99,
					Currency: 2,
				},
			},
			expected: []dto.ListServerResp{
				{
					Model:    "Dell R730XD2x Intel Xeon E5-2650v3",
					Ram:      "16GB",
					HDD:      "2x500GBSATA2",
					Location: "AmsterdamAMS-01",
					Price:    "€99.99",
				},
			},
		},
		{
			name: "transform server with unknown HDD type",
			input: []models.ServerCatalog{
				{
					Model:    "IBM X3650M42x Intel Xeon E5-2620",
					RamSize:  16,
					RamType:  utils.RAMTypeDDR3,
					HDDSize:  500,
					HDDCount: 2,
					HDDType:  999, // Unknown type
					Location: "AmsterdamAMS-01",
					Price:    99.99,
					Currency: 2,
				},
			},
			expected: []dto.ListServerResp{
				{
					Model:    "IBM X3650M42x Intel Xeon E5-2620",
					Ram:      "16GBDDR3",
					HDD:      "2x500GB",
					Location: "AmsterdamAMS-01",
					Price:    "€99.99",
				},
			},
		},
		{
			name: "transform server with unknown currency",
			input: []models.ServerCatalog{
				{
					Model:    "Dell R730XD2x Intel Xeon E5-2650v3",
					RamSize:  16,
					RamType:  utils.RAMTypeDDR3,
					HDDSize:  500,
					HDDCount: 2,
					HDDType:  utils.HDDTypeSATA2,
					Location: "AmsterdamAMS-01",
					Price:    99.99,
					Currency: 999,
				},
			},
			expected: []dto.ListServerResp{
				{
					Model:    "Dell R730XD2x Intel Xeon E5-2650v3",
					Ram:      "16GBDDR3",
					HDD:      "2x500GBSATA2",
					Location: "AmsterdamAMS-01",
					Price:    "99.99",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TransformServerList(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
