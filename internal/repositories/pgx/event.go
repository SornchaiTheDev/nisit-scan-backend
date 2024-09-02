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

type eventRepoImpl struct {
	q *sqlc.Queries
}

func NewEventRepo(q *sqlc.Queries) services.EventRepository {
	return &eventRepoImpl{
		q: q,
	}
}

func (e *eventRepoImpl) GetAll() ([]*entities.Event, error) {
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

func (e *eventRepoImpl) GetById(id uuid.UUID) (*entities.Event, error) {
	event, err := e.q.GetEventById(context.Background(), id)
	if err != nil {
		return nil, err
	}

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

func (e *eventRepoImpl) Create(event *entities.Event, adminId string) error {

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

func (e *eventRepoImpl) DeleteById(id uuid.UUID) error {
	err := e.q.DeleteEventById(context.Background(), id)

	return err
}

func (e *eventRepoImpl) UpdateById(id uuid.UUID, event *entities.Event) error {
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
