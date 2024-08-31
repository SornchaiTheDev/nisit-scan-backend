package requests

type AdminRequest struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}
