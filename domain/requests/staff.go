package requests

type SetStaffRequest struct {
	Email []string `json:"emails" validate:"required"`
}
