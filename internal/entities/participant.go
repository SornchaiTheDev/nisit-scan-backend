package entities

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	Id        uuid.UUID `json:"id"`
	Barcode   string    `json:"barcode"`
	Timestamp time.Time `json:"timestamp"`
}
