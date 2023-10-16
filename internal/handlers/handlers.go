package handlers

import (
	"canonicalAuditlog/internal/jwtutil"
	"canonicalAuditlog/internal/models"
	"canonicalAuditlog/internal/services"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// Handlers Todo: Split out event and jwt handlers if logic increases
type Handlers struct {
	Service services.EventService
}

func (h *Handlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateEvent(event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.Service.SaveEvent(event)
	if err != nil {
		log.Printf("Error saving event: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) QueryEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawParams := r.URL.Query()
	params := make(map[string]string)

	for key, values := range rawParams {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	events, err := h.Service.QueryEvents(params)

	if err != nil {
		http.Error(w, "Error fetching events", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		return
	}
}

func (h *Handlers) GetJwt(w http.ResponseWriter, _ *http.Request) {
	token, err := jwtutil.GenerateToken()
	if err != nil {
		log.Printf("Error generating token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{
			"error": "Internal Server Error",
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"token": token,
	}
	json.NewEncoder(w).Encode(response)
}

func validateEvent(e models.Event) error {
	if e.EventType == "" || e.ServiceName == "" || e.CustomerID == "" || len(e.EventDetails) == 0 {
		return errors.New("invalid request body: required fields event_type, service_name, customer_id, event_details")
	}

	return nil
}
