package dtos

type PaginationParams struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}
