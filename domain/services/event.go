package services

import (
	"errors"
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EventService interface {
	GetAll() ([]*entities.Event, error)
	GetById(id string) (*entities.Event, error)
	Create(e *requests.EventRequest, adminId string) error
	DeleteById(id string) error
	UpdateById(id string, r *requests.EventRequest) error
}

type eventService struct {
	repo repositories.EventRepository
}

func NewEventService(repo repositories.EventRepository) EventService {
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
			return nerrors.ErrEventNotFound
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
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = s.isEventExist(&parsedId)
	if err != nil {
		return err
	}

	return s.repo.DeleteById(parsedId)
}

func (s *eventService) UpdateById(id string, r *requests.EventRequest) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = s.isEventExist(&parsedId)
	if err != nil {
		return err
	}

	event, err := parseRequestToEntity(r)
	if err != nil {
		return err
	}

	return s.repo.UpdateById(parsedId, event)
}
