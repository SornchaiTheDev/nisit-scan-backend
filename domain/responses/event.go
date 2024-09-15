package responses

import (
	"time"

	"github.com/google/uuid"
)

type EventResponse struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Place             string    `json:"place"`
	Date              time.Time `json:"date"`
	Host              string    `json:"host"`
	Owner             string    `json:"owner"`
	ParticipantsCount int64     `json:"participants_count"`
}
