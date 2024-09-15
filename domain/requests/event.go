package requests

type EventRequest struct {
	Name  string `json:"name" validate:"required,min=1"`
	Place string `json:"place" validate:"required,min=1"`
	Date  string `json:"date" validate:"required,date"`
	Host  string `json:"host" validate:"required,min=1"`
}
