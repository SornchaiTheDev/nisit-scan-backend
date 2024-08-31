package entities

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	Id        uuid.UUID
	Email     string
	FullName  string
	DeletedAt time.Time
}
