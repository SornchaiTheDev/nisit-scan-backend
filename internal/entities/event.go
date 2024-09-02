package entities

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Place string    `json:"place"`
	Date  time.Time `json:"date"`
	Host  string    `json:"host"`
	Owner string    `json:"owner"`
}
