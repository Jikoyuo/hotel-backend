package utils

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
}

type PaginatedResponse struct {
	Meta PaginationMeta `json:"meta"`
	Data interface{}    `json:"data"`
}
