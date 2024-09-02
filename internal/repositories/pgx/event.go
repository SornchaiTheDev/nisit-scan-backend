package repositories

import (
	"context"

	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/services"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type EventRepoImpl struct {
	q *sqlc.Queries
}

func NewEventRepo(q *sqlc.Queries) services.EventRepository {
	return &EventRepoImpl{
		q: q,
	}
}

func (e *EventRepoImpl) GetAll() ([]*entities.Event, error) {
	events, err := e.q.GetAllEvents(context.Background())

	var parsedEvents []*entities.Event

	for _, event := range events {
		parsedEvent := &entities.Event{
			Id:    event.ID,
			Name:  event.Name,
			Place: event.Place,
			Date:  event.Date.Time,
			Host:  event.Host,
			Owner: event.FullName,
		}

		parsedEvents = append(parsedEvents, parsedEvent)
	}

	return parsedEvents, err
}

func (e *EventRepoImpl) GetById(id uuid.UUID) (*entities.Event, error) {
	event, err := e.q.GetEventById(context.Background(), id)

	parsedEvent := &entities.Event{
		Id:    event.ID,
		Name:  event.Name,
		Place: event.Place,
		Date:  event.Date.Time,
		Host:  event.Host,
		Owner: event.FullName,
	}

	return parsedEvent, err
}

func (e *EventRepoImpl) Create(event *entities.Event, adminId string) error {

	date := pgtype.Date{}
	date.Scan(event.Date)

	parseId, err := libs.ParseUUID(adminId)
	if err != nil {
		return err
	}

	err = e.q.CreateEvent(context.Background(), sqlc.CreateEventParams{
		Name:    event.Name,
		Place:   event.Place,
		Date:    date,
		Host:    event.Host,
		AdminID: *parseId,
	})

	return err
}

func (e *EventRepoImpl) DeleteById(id uuid.UUID) error {
	err := e.q.DeleteEventById(context.Background(), id)

	return err
}

func (e *EventRepoImpl) UpdateById(id uuid.UUID, event *entities.Event) error {
	date := pgtype.Date{}
	date.Scan(event.Date)

	err := e.q.UpdateEventById(context.Background(), sqlc.UpdateEventByIdParams{
		ID:    id,
		Name:  event.Name,
		Place: event.Place,
		Date:  date,
		Host:  event.Host,
	})

	return err
}
