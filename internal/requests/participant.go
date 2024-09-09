package requests

import "time"

type AddParticipant struct {
	Barcode   string    `json:"barcode"`
	Timestamp time.Time `json:"timestamp"`
}
