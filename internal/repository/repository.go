package repository

import "canonicalAuditlog/internal/models"

type EventRepository interface {
	Insert(event models.Event) error
	Find(filter map[string]string) ([]models.Event, error)
}
