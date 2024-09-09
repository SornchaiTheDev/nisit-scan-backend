package entities

import "github.com/google/uuid"

type Staff struct {
	Id    uuid.UUID `json:"id"`
	Email string `json:"email"`
}
