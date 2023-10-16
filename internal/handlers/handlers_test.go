package handlers

import (
	"bytes"
	"canonicalAuditlog/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockEventService struct {
	mock.Mock
}

func (m *MockEventService) SaveEvent(event models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventService) QueryEvents(filter map[string]string) ([]models.Event, error) {
	args := m.Called(filter)
	return args.Get(0).([]models.Event), args.Error(1)
}

func TestCreateEvents(t *testing.T) {
	mockService := new(MockEventService)
	handler := &Handlers{Service: mockService}

	tests := []struct {
		name           string
		inputEvent     models.Event
		mockReturn     error
		expectedStatus int
	}{
		{
			name: "Valid Event Creation",
			inputEvent: models.Event{
				EventType:    "TestEvent",
				ServiceName:  "TestService",
				CustomerID:   "12345",
				EventDetails: map[string]interface{}{"key": "value"},
			},
			mockReturn:     nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing EventType",
			inputEvent: models.Event{
				ServiceName:  "TestService",
				CustomerID:   "12345",
				EventDetails: map[string]interface{}{"key": "value"},
			},
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventJson, _ := json.Marshal(tt.inputEvent)
			req, _ := http.NewRequest(http.MethodPost, "/events/create", bytes.NewBuffer(eventJson))
			rec := httptest.NewRecorder()

			mockService.On("SaveEvent", tt.inputEvent).Return(tt.mockReturn)
			handler.CreateEvent(rec, req)

			if status := rec.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestQueryEvents(t *testing.T) {
	mockService := new(MockEventService)
	handler := &Handlers{Service: mockService}

	tests := []struct {
		name           string
		inputFilter    map[string]string
		mockReturn     []models.Event
		mockError      error
		expectedStatus int
	}{
		{
			name: "Valid Event Query",
			inputFilter: map[string]string{
				"event_type": "TestEvent",
			},
			mockReturn: []models.Event{
				{
					EventType:    "TestEvent",
					ServiceName:  "TestService",
					CustomerID:   "12345",
					EventDetails: map[string]interface{}{"key": "value"},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty Filter",
			inputFilter:    map[string]string{},
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/events/query", nil)
			q := req.URL.Query()
			for k, v := range tt.inputFilter {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			rec := httptest.NewRecorder()

			mockService.On("QueryEvents", tt.inputFilter).Return(tt.mockReturn, tt.mockError)
			handler.QueryEvents(rec, req)

			if status := rec.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestGetJwt(t *testing.T) {
	mockService := new(MockEventService)
	handler := &Handlers{Service: mockService}

	req, _ := http.NewRequest(http.MethodGet, "/jwt", nil)
	rec := httptest.NewRecorder()
	handler.GetJwt(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	json.Unmarshal(rec.Body.Bytes(), &response)

	if _, ok := response["token"]; !ok {
		t.Errorf("Token not found in response")
	}
}
