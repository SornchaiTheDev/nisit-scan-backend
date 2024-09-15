package requests

type AdminRequest struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"fullName" validate:"required,min=1"`
}

type GetAdminsPaginationParams struct {
	Search    string `json:"search"`
	PageIndex int32  `json:"pageIndex" validate:"required,min=0"`
	PageSize  int32  `json:"pageSize" validate:"required,min=1"`
}
