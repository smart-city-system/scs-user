package types

type PaginateResponse[T any] struct {
	Data       []T `json:"data"`
	Pagination `json:"pagination"`
}
type Pagination struct {
	TotalPages int `json:"total_pages"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
