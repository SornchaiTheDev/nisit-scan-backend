package entities

import (
	"time"
)

type Participant struct {
	Barcode   string    `json:"barcode"`
	Timestamp time.Time `json:"timestamp"`
}
