package repository

import (
	"context"
	"testing"

	"github.com/server-catalog/internal/dto"
	"github.com/server-catalog/internal/utils"
	"github.com/server-catalog/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.ServerCatalog{}, &models.HDDSpec{})
	assert.NoError(t, err)

	return db
}

func TestServerCatalog_Upload(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServerCatalog(db)
	ctx := context.Background()

	// tests ....
	servers := []models.ServerCatalog{
		{
			Model:    "Dell R210-II",
			RamSize:  16,
			RamType:  1,
			HDDSize:  500,
			HDDCount: 2,
			HDDType:  1,
			Location: "Amsterdam",
			Price:    35.99,
			Currency: 1,
		},
		{
			Model:    "HP DL120G7",
			RamSize:  32,
			RamType:  2,
			HDDSize:  1000,
			HDDCount: 4,
			HDDType:  2,
			Location: "Singapore",
			Price:    45.99,
			Currency: 1,
		},
	}

	err := repo.Upload(ctx, servers)
	assert.NoError(t, err)

	var count int64
	db.Model(&models.ServerCatalog{}).Count(&count)
	assert.Equal(t, int64(2), count)

	var result models.ServerCatalog
	db.First(&result, "model = ?", "Dell R210-II")
	assert.Equal(t, "Dell R210-II", result.Model)
	assert.Equal(t, 16, result.RamSize)
	assert.Equal(t, "Amsterdam", result.Location)
}

func TestServerCatalog_GetLocations(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServerCatalog(db)
	ctx := context.Background()

	// Insert test data
	testData := []models.ServerCatalog{
		{Location: "Amsterdam"},
		{Location: "Singapore"},
		{Location: "Amsterdam"},
		{Location: "London"},
	}
	db.Create(&testData)

	locations, err := repo.GetLocations(ctx)
	assert.NoError(t, err)
	assert.Len(t, locations, 3)
	assert.Contains(t, locations, "Amsterdam")
	assert.Contains(t, locations, "Singapore")
	assert.Contains(t, locations, "London")
}

func TestServerCatalog_GetHDDTypes(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServerCatalog(db)
	ctx := context.Background()

	testData := []models.HDDSpec{
		{Type: "SATA2"},
		{Type: "SAS"},
		{Type: "SATA2"},
		{Type: "SSD"},
	}
	db.Create(&testData)

	// Test GetHDDTypes
	types, err := repo.GetHDDTypes(ctx)
	assert.NoError(t, err)
	assert.Len(t, types, 3)
	assert.Contains(t, types, "SATA2")
	assert.Contains(t, types, "SAS")
	assert.Contains(t, types, "SSD")
}

func TestServerCatalog_GetServers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServerCatalog(db)
	ctx := context.Background()

	testData := []models.ServerCatalog{
		{
			Model:    "Server 1",
			RamSize:  16,
			HDDSize:  500,
			HDDCount: 2,
			HDDType:  1, // SATA2
			Location: "Amsterdam",
			Price:    35.99,
		},
		{
			Model:    "Server 2",
			RamSize:  32,
			HDDSize:  1000,
			HDDCount: 4,
			HDDType:  2, // SAS
			Location: "Singapore",
			Price:    45.99,
		},
	}
	db.Create(&testData)

	tests := []struct {
		name          string
		ctr           *dto.ListServersCtr
		expectedCount int
	}{
		{
			name: "filter by RAM",
			ctr: &dto.ListServersCtr{
				RAM: []int{16},
				Page: &utils.Page{
					Limit:   10,
					Current: 1,
				},
			},
			expectedCount: 1,
		},
		{
			name: "filter by HDD type",
			ctr: &dto.ListServersCtr{
				HDD: &[]int{1}[0], // SATA2
				Page: &utils.Page{
					Limit:   10,
					Current: 1,
				},
			},
			expectedCount: 1,
		},
		{
			name: "filter by location",
			ctr: &dto.ListServersCtr{
				Location: &[]string{"Amsterdam"}[0],
				Page: &utils.Page{
					Limit:   10,
					Current: 1,
				},
			},
			expectedCount: 1,
		},
		{
			name: "filter by storage range",
			ctr: &dto.ListServersCtr{
				StorageMin: &[]int{1000}[0],
				StorageMax: &[]int{5000}[0],
				Page: &utils.Page{
					Limit:   1,
					Current: 1,
				},
			},
			expectedCount: 1,
		},
		{
			name: "pagination test",
			ctr: &dto.ListServersCtr{
				Page: &utils.Page{
					Limit:   1,
					Current: 1,
				},
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			servers, err := repo.GetServers(ctx, tt.ctr)
			assert.NoError(t, err)
			assert.Len(t, servers, tt.expectedCount)

			assert.Equal(t, 2, tt.ctr.Page.Total)
		})
	}
}

func TestServerCatalog_GetServers_Error(t *testing.T) {
	db := setupTestDB(t)
	_ = NewServerCatalog(db)
	ctx := context.Background()

	invalidDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	assert.NoError(t, err)

	invalidRepo := NewServerCatalog(invalidDB)

	// Test with invalid table name to trigger error
	invalidDB.Exec("DROP TABLE IF EXISTS server_catalogs")

	_, err = invalidRepo.GetServers(ctx, &dto.ListServersCtr{
		Page: &utils.Page{
			Limit:   10,
			Current: 1,
		},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch count of servers")
}
