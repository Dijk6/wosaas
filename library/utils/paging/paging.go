// Package paging 分页工具。
package paging

import "reflect"

const (
	pageDefault  = 1
	limitDefault = 10
)

type Paging struct {
	page  int
	limit int
}

// NewPaging 创建分页实例。
// 一个参数: limit；两个参数: page, limit。
func NewPaging(params ...int) *Paging {
	p := Paging{}
	switch len(params) {
	case 1:
		p.SetLimit(params[0])
	case 2:
		p.SetPage(params[0])
		p.SetLimit(params[1])
	}
	p.init()
	return &p
}

func (p *Paging) init() {
	if p.page <= 0 {
		p.page = pageDefault
	}
	if p.limit <= 0 {
		p.limit = limitDefault
	}
}

func (p *Paging) SetPage(page int)   { p.page = page }
func (p *Paging) IncPage()           { p.page++ }
func (p *Paging) GetPage() int       { return p.page }
func (p *Paging) SetLimit(limit int) { p.limit = limit }
func (p *Paging) GetLimit() int      { return p.limit }
func (p *Paging) GetOffset() int     { return p.page*p.limit - p.limit }

// Paginate 对任意切片分页。
func (p *Paging) Paginate(data any) any {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil
	}

	total := v.Len()
	if total == 0 {
		return reflect.MakeSlice(v.Type(), 0, 0).Interface()
	}

	totalPages := total / p.limit
	if total%p.limit != 0 {
		totalPages++
	}
	if p.page > totalPages {
		p.page = totalPages
	}

	start := (p.page - 1) * p.limit
	end := start + p.limit
	if end > total {
		end = total
	}
	return v.Slice(start, end).Interface()
}
