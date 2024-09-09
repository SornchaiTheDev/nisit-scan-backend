package responses

import (
	"time"

	"github.com/google/uuid"
)

type AdminResponse struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"fullName"`
}

type AllAdminResponse struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullName"`
	DeletedAt *time.Time `json:"deletedAt"`
}
