package services

import (
	"errors"
	"time"

	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/requests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EventRepository interface {
	GetAll() ([]*entities.Event, error)
	GetById(id uuid.UUID) (*entities.Event, error)
	Create(e *entities.Event, adminId string) error
	DeleteById(id uuid.UUID) error
	UpdateById(id uuid.UUID, e *entities.Event) error
}

type eventService struct {
	repo EventRepository
}

func NewEventService(repo EventRepository) *eventService {
	return &eventService{
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

func (s *eventService) isEventExist(id *uuid.UUID) error {
	_, err := s.repo.GetById(*id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrEventNotFound
		}
	}
	return nil
}

func (s *eventService) GetAll() ([]*entities.Event, error) {
	return s.repo.GetAll()
}

func (s *eventService) GetById(id string) (*entities.Event, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetById(parsedId)
}

func (s *eventService) Create(r *requests.EventRequest, adminId string) error {
	event, err := parseRequestToEntity(r)
	if err != nil {
		return err
	}

	return s.repo.Create(event, adminId)
}

func (s *eventService) DeleteById(id string) error {
	parsedId, err := libs.ParseUUID(id)
	if err != nil {
		return err
	}

	err = s.isEventExist(parsedId)
	if err != nil {
		return err
	}

	return s.repo.DeleteById(*parsedId)
}

func (s *eventService) UpdateById(id string, r *requests.EventRequest) error {
	parsedId, err := libs.ParseUUID(id)
	if err != nil {
		return err
	}

	err = s.isEventExist(parsedId)
	if err != nil {
		return err
	}

	event, err := parseRequestToEntity(r)
	if err != nil {
		return err
	}

	return s.repo.UpdateById(*parsedId, event)
}
