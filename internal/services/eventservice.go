package services

import (
	"canonicalAuditlog/internal/models"
	"canonicalAuditlog/internal/repository"
	"time"
)

type eventServiceImpl struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventServiceImpl{repo: repo}
}

func (s *eventServiceImpl) SaveEvent(event models.Event) error {
	event.Timestamp = time.Now()
	return s.repo.Insert(event)
}

func (s *eventServiceImpl) QueryEvents(filter map[string]string) ([]models.Event, error) {
	return s.repo.Find(filter)
}
