package requests

type EventRequest struct {
	Name  string `json:"name"`
	Place string `json:"place"`
	Date  string `json:"date"`
	Host  string `json:"host"`
}
