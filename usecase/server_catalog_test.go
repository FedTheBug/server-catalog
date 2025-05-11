package usecase

import (
	"bytes"
	"context"
	"errors"
	"github.com/server-catalog/internal/dto"
	"github.com/server-catalog/internal/utils"
	"github.com/server-catalog/models"
	_ "github.com/server-catalog/repository"
	"github.com/xuri/excelize/v2"
	_ "mime/multipart"
	"testing"
)

type mockFile struct {
	*bytes.Reader
}

func (m *mockFile) Close() error {
	return nil
}

type mockCatalogRepository struct {
	uploadFunc       func(ctx context.Context, catalogs []models.ServerCatalog) error
	getLocationsFunc func(ctx context.Context) ([]string, error)
	getHDDTypesFunc  func(ctx context.Context) ([]string, error)
	getServersFunc   func(ctx context.Context, ctr *dto.ListServersCtr) ([]models.ServerCatalog, error)
}

func (m *mockCatalogRepository) Upload(ctx context.Context, catalogs []models.ServerCatalog) error {
	return m.uploadFunc(ctx, catalogs)
}

func (m *mockCatalogRepository) GetLocations(ctx context.Context) ([]string, error) {
	return m.getLocationsFunc(ctx)
}

func (m *mockCatalogRepository) GetHDDTypes(ctx context.Context) ([]string, error) {
	return m.getHDDTypesFunc(ctx)
}

func (m *mockCatalogRepository) GetServers(ctx context.Context, ctr *dto.ListServersCtr) ([]models.ServerCatalog, error) {
	return m.getServersFunc(ctx, ctr)
}

func createTestExcelFile(data [][]string) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"

	for i, row := range data {
		for j, cell := range row {
			cellName, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				return nil, err
			}
			f.SetCellValue(sheet, cellName, cell)
		}
	}

	return f.WriteToBuffer()
}

func TestServerCatalog_UploadCatalog(t *testing.T) {
	tests := []struct {
		name          string
		excelData     [][]string
		mockUpload    func(ctx context.Context, catalogs []models.ServerCatalog) error
		expectedError error
	}{
		{
			name: "valid catalog upload",
			excelData: [][]string{
				{"Model", "RAM", "HDD", "Location", "Price"},
				{"Dell R210-II", "16GB DDR3", "2x500GBSATA2", "AmsterdamAMS-01", "$35.99"},
			},
			mockUpload: func(ctx context.Context, catalogs []models.ServerCatalog) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "invalid header",
			excelData: [][]string{
				{"Invalid", "Header", "Format", "Location", "Price"},
				{"Dell R210-II", "16GB DDR3", "2x500GBSATA2", "AmsterdamAMS-01", "$35.99"},
			},
			mockUpload: func(ctx context.Context, catalogs []models.ServerCatalog) error {
				return nil
			},
			expectedError: errors.New("usecase:server_catalog:invalid XLSX columns"),
		},
		{
			name: "repository error",
			excelData: [][]string{
				{"Model", "RAM", "HDD", "Location", "Price"},
				{"Dell R210-II", "16GB DDR3", "2x500GBSATA2", "AmsterdamAMS-01", "$35.99"},
			},
			mockUpload: func(ctx context.Context, catalogs []models.ServerCatalog) error {
				return errors.New("database error")
			},
			expectedError: errors.New("usecase:server_catalog:: failed to upload failed to upload data into the database"),
		},
		{
			name: "too many rows",
			excelData: func() [][]string {
				
				data := make([][]string, 1002)
				data[0] = []string{"Model", "RAM", "HDD", "Location", "Price"}
				for i := 1; i < 1002; i++ {
					data[i] = []string{"Dell R210-II", "16GB DDR3", "2x500GBSATA2", "AmsterdamAMS-01", "$35.99"}
				}
				return data
			}(),
			mockUpload: func(ctx context.Context, catalogs []models.ServerCatalog) error {
				return nil
			},
			expectedError: errors.New("usecase:server_catalog:maximum number of rows exceeded (1000)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			excelBuffer, err := createTestExcelFile(tt.excelData)
			if err != nil {
				t.Fatalf("Failed to create test Excel file: %v", err)
			}

			// mock repository
			mockRepo := &mockCatalogRepository{
				uploadFunc: tt.mockUpload,
			}

			uc := New(mockRepo)

			ctr := &dto.UploadCatalogCtr{
				File: &mockFile{bytes.NewReader(excelBuffer.Bytes())},
			}

			// Test upload
			err = uc.UploadCatalog(context.Background(), ctr)
			if (err != nil && tt.expectedError == nil) ||
				(err == nil && tt.expectedError != nil) ||
				(err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("UploadCatalog() error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

func TestServerCatalog_GetLocations(t *testing.T) {
	tests := []struct {
		name          string
		mockLocations func(ctx context.Context) ([]string, error)
		expected      []string
		expectedError error
	}{
		{
			name: "successful locations retrieval",
			mockLocations: func(ctx context.Context) ([]string, error) {
				return []string{"Amsterdam", "Singapore", "London"}, nil
			},
			expected:      []string{"Amsterdam", "Singapore", "London"},
			expectedError: nil,
		},
		{
			name: "repository error",
			mockLocations: func(ctx context.Context) ([]string, error) {
				return nil, errors.New("database error")
			},
			expected:      nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockCatalogRepository{
				getLocationsFunc: tt.mockLocations,
			}

			uc := New(mockRepo)
			locations, err := uc.GetLocations(context.Background())

			if (err != nil && tt.expectedError == nil) ||
				(err == nil && tt.expectedError != nil) ||
				(err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("GetLocations() error = %v, want %v", err, tt.expectedError)
			}

			if err == nil && !compareStringSlices(locations, tt.expected) {
				t.Errorf("GetLocations() = %v, want %v", locations, tt.expected)
			}
		})
	}
}

func TestServerCatalog_GetHDDTypes(t *testing.T) {
	tests := []struct {
		name          string
		mockHDDTypes  func(ctx context.Context) ([]string, error)
		expected      []string
		expectedError error
	}{
		{
			name: "successful HDD types retrieval",
			mockHDDTypes: func(ctx context.Context) ([]string, error) {
				return []string{"SATA2", "SAS", "SSD"}, nil
			},
			expected:      []string{"SATA2", "SAS", "SSD"},
			expectedError: nil,
		},
		{
			name: "repository error",
			mockHDDTypes: func(ctx context.Context) ([]string, error) {
				return nil, errors.New("database error")
			},
			expected:      nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockCatalogRepository{
				getHDDTypesFunc: tt.mockHDDTypes,
			}

			uc := New(mockRepo)
			hddTypes, err := uc.GetHDDTypes(context.Background())

			if (err != nil && tt.expectedError == nil) ||
				(err == nil && tt.expectedError != nil) ||
				(err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("GetHDDTypes() error = %v, want %v", err, tt.expectedError)
			}

			if err == nil && !compareStringSlices(hddTypes, tt.expected) {
				t.Errorf("GetHDDTypes() = %v, want %v", hddTypes, tt.expected)
			}
		})
	}
}

func TestServerCatalog_GetListOfServers(t *testing.T) {
	amsterdam := "Amsterdam"
	invalid := "Invalid"
	hddType := 1 // SATA2

	tests := []struct {
		name          string
		ctr           *dto.ListServersCtr
		mockServers   func(ctx context.Context, ctr *dto.ListServersCtr) ([]models.ServerCatalog, error)
		expectedError error
	}{
		{
			name: "successful server list retrieval",
			ctr: &dto.ListServersCtr{
				Location: &amsterdam,
				HDD:      &hddType,
				RAM:      []int{16},
			},
			mockServers: func(ctx context.Context, ctr *dto.ListServersCtr) ([]models.ServerCatalog, error) {
				return []models.ServerCatalog{
					{
						Model:    "Dell R210-II",
						RamSize:  16,
						RamType:  utils.RAMTypeDDR3,
						HDDSize:  500,
						HDDCount: 2,
						HDDType:  utils.HDDTypeSATA2,
						Location: "Amsterdam",
						Price:    35.99,
						Currency: utils.CurrencyUSD,
					},
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "no servers found",
			ctr: &dto.ListServersCtr{
				Location: &invalid,
			},
			mockServers: func(ctx context.Context, ctr *dto.ListServersCtr) ([]models.ServerCatalog, error) {
				return []models.ServerCatalog{}, nil
			},
			expectedError: utils.ErrServerNotFound,
		},
		{
			name: "repository error",
			ctr: &dto.ListServersCtr{
				Location: &amsterdam,
			},
			mockServers: func(ctx context.Context, ctr *dto.ListServersCtr) ([]models.ServerCatalog, error) {
				return nil, errors.New("database error")
			},
			expectedError: errors.New("usecase:server_catalog:: failed to get list of servers database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockCatalogRepository{
				getServersFunc: tt.mockServers,
			}

			uc := New(mockRepo)
			servers, err := uc.GetListOfServers(context.Background(), tt.ctr)

			if (err != nil && tt.expectedError == nil) ||
				(err == nil && tt.expectedError != nil) ||
				(err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("GetListOfServers() error = %v, want %v", err, tt.expectedError)
			}

			if err == nil && len(servers) == 0 {
				t.Error("GetListOfServers() returned empty list when no error was expected")
			}
		})
	}
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
