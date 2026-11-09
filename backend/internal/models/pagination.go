package models

type Pagination struct {
	Page int
	Size int
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}
