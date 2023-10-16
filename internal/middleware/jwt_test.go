package middleware

import (
	"canonicalAuditlog/internal/jwtutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWTMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Valid JWT Token",
			authHeader:     "Bearer " + generateValidToken(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing Authorization Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Authorization Format",
			authHeader:     "Invalid " + generateValidToken(),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid JWT Token",
			authHeader:     "Bearer invalidToken12345",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := JWTMiddleware(nextHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func generateValidToken() string {
	token, _ := jwtutil.GenerateToken()
	return token
}
