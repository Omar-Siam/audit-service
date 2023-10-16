package services

import (
	"canonicalAuditlog/internal/models"
)

type EventService interface {
	SaveEvent(event models.Event) error
	QueryEvents(filter map[string]string) ([]models.Event, error)
}
