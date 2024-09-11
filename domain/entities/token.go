package entities

type RefreshToken struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
