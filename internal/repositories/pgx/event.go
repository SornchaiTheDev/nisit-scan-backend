package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type eventRepoImpl struct {
	q *sqlc.Queries
}

func NewEventRepo(q *sqlc.Queries) repositories.EventRepository {
	return &eventRepoImpl{
		q: q,
	}
}

func (e *eventRepoImpl) GetPagination(search string, pageIndex int32, pageSize int32) ([]*entities.Event, error) {
	events, err := e.q.GetAllEvents(context.Background(), sqlc.GetAllEventsParams{
		Name:   fmt.Sprintf("%%%s%%", search),
		Offset: pageIndex,
		Limit:  pageSize,
	})

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

func (e *eventRepoImpl) GetCount(search string) (int64, error) {
	count, err := e.q.GetEventCount(context.Background(), fmt.Sprintf("%%%s%%", search))
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (e *eventRepoImpl) GetById(id uuid.UUID) (*entities.Event, error) {
	event, err := e.q.GetEventById(context.Background(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nerrors.ErrEventNotFound
		}
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

	parseId, err := uuid.Parse(adminId)
	if err != nil {
		return err
	}

	err = e.q.CreateEvent(context.Background(), sqlc.CreateEventParams{
		Name:    event.Name,
		Place:   event.Place,
		Date:    date,
		Host:    event.Host,
		AdminID: parseId,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nerrors.ErrEventAlreadyExists
			}
			if pgErr.Code == "23503" {
				return nerrors.ErrAdminNotFound
			}
		}
	}

	return nil

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

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nerrors.ErrEventAlreadyExists
			}
		}
	}

	return err
}
