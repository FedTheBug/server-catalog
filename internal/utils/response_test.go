package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponse_Render(t *testing.T) {
	tests := []struct {
		name           string
		response       Response
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Success",
			response: Response{
				Status:  http.StatusOK,
				Message: "Success",
				Data:    map[string]interface{}{"key": "value"},
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Success",
				"data":    map[string]interface{}{"key": "value"},
			},
		},
		{
			name: "error response",
			response: Response{
				Status:  http.StatusBadRequest,
				Message: "Error occurred",
				Error:   "Invalid input",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "Error occurred",
				"error":   "Invalid input",
			},
		},
		{
			name: "empty response",
			response: Response{
				Status:  http.StatusOK,
				Message: "",
				Data:    nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "",
				"data":    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := tt.response.Render(w)
			if err != nil {
				t.Errorf("Render() error = %v", err)
				return
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			err = json.NewDecoder(w.Body).Decode(&response)
			if err != nil {
				t.Errorf("failed to parse response body: %v", err)
				return
			}

			if message, ok := response["message"].(string); !ok {
				t.Errorf("message field is not a string: %v", response["message"])
			} else if message != tt.expectedBody["message"] {
				t.Errorf("expected message %v; got %v", tt.expectedBody["message"], message)
			}

			if tt.expectedBody["error"] != nil {
				if err, ok := response["error"].(string); !ok {
					t.Errorf("error field is not a string: %v", response["error"])
				} else if err != tt.expectedBody["error"] {
					t.Errorf("expected error %v; got %v", tt.expectedBody["error"], err)
				}
			}

			if tt.expectedBody["data"] != nil {
				if data, ok := response["data"].(map[string]interface{}); !ok {
					t.Errorf("data field is not a map: %v", response["data"])
				} else if expectedData, ok := tt.expectedBody["data"].(map[string]interface{}); ok {
					if data["key"] != expectedData["key"] {
						t.Errorf("expected data %v; got %v", tt.expectedBody["data"], data)
					}
				}
			}
		})
	}
}
