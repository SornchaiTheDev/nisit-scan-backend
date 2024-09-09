package requests

type AdminRequest struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}

type GetAdminsPaginationParams struct {
	Search    string `json:"search"`
	PageIndex int32  `json:"pageIndex"`
	PageSize  int32  `json:"pageSize"`
}
