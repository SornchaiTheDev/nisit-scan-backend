package requests

type AddParticipant struct {
	Barcode   string    `json:"barcode" validate:"required"`
	Timestamp string `json:"timestamp" validate:"required,timestamp"`
}
