package entities

import (
	"time"
)

type Participant struct {
	Barcode     string    `json:"barcode"`
	Timestamp   time.Time `json:"timestamp"`
	StudentCode string    `json:"student_code"`
	FullName    string    `json:"full_name"`
	Gmail       string    `json:"gmail"`
	Major       string    `json:"major"`
}
