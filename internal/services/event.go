package services

import (
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/requests"
	"github.com/google/uuid"
)

type EventRepository interface {
	GetAll() ([]*entities.Event, error)
	GetById(id uuid.UUID) (*entities.Event, error)
	Create(e *entities.Event, adminId string) error
	DeleteById(id uuid.UUID) error
	UpdateById(id uuid.UUID, e *entities.Event) error
}

type EventService struct {
	repo EventRepository
}

func NewEventService(repo EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func parseRequestToEntity(r *requests.EventRequest) (*entities.Event, error) {

	date, err := time.Parse("02/01/2006", r.Date)
	if err != nil {
		return nil, err
	}

	event := &entities.Event{
		Name:  r.Name,
		Place: r.Place,
		Date:  date,
		Host:  r.Host,
	}

	return event, nil
}

func (s *EventService) GetAll() ([]*entities.Event, error) {
	return s.repo.GetAll()
}

func (s *EventService) GetById(id string) (*entities.Event, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetById(parsedId)
}

func (s *EventService) Create(r *requests.EventRequest, adminId string) error {
	event, err := parseRequestToEntity(r)
	if err != nil {
		return err
	}

	return s.repo.Create(event, adminId)
}

func (s *EventService) DeleteById(id string) error {
	parsedId, err := libs.ParseUUID(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteById(*parsedId)
}

func (s *EventService) UpdateById(id string, r *requests.EventRequest) error {
	parsedId, err := libs.ParseUUID(id)
	if err != nil {
		return err
	}

	event, err := parseRequestToEntity(r)
	if err != nil {
		return err
	}

	return s.repo.UpdateById(*parsedId, event)
}
