package utils

import (
	"net/http"
	"testing"
)

func TestPage_Offset(t *testing.T) {
	tests := []struct {
		name     string
		page     Page
		expected int
	}{
		{
			name: "first page",
			page: Page{
				Limit:   10,
				Current: 1,
			},
			expected: 0,
		},
		{
			name: "second page",
			page: Page{
				Limit:   10,
				Current: 2,
			},
			expected: 10,
		},
		{
			name: "third page",
			page: Page{
				Limit:   10,
				Current: 3,
			},
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.page.Offset(); got != tt.expected {
				t.Errorf("Page.Offset() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewPage(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedLimit  int
		expectedPageNo int
	}{
		{
			name:           "default values",
			url:            "/api/servers",
			expectedLimit:  10,
			expectedPageNo: 1,
		},
		{
			name:           "custom values",
			url:            "/api/servers?per_page=20&page_no=2",
			expectedLimit:  20,
			expectedPageNo: 2,
		},
		{
			name:           "invalid values",
			url:            "/api/servers?per_page=0&page_no=0",
			expectedLimit:  10,
			expectedPageNo: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			page := NewPage(req)
			if page.Limit != tt.expectedLimit {
				t.Errorf("NewPage()-Limit = %v, want %v", page.Limit, tt.expectedLimit)
			}
			if page.Current != tt.expectedPageNo {
				t.Errorf("NewPage()-Current = %v, want %v", page.Current, tt.expectedPageNo)
			}
		})
	}
}
