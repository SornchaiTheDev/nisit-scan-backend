package requests

import "github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"

type GetUsersPaginationParams struct {
	Search    string `json:"search"`
	PageIndex int    `json:"pageIndex" validate:"required,min=0"`
	PageSize  int    `json:"pageSize" validate:"required,min=1"`
}

type ImportUsers struct {
	Users []entities.User `json:"users"`
}
