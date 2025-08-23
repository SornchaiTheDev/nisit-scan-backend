package requests

type GetUsersPaginationParams struct {
	Search    string `json:"search"`
	PageIndex int32  `json:"pageIndex" validate:"required,min=0"`
	PageSize  int32  `json:"pageSize" validate:"required,min=1"`
}
