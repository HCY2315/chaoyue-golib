package http

import "git.cestong.com.cn/cecf/cecf-golib/pkg/errors"

// 应避免使用属性
type PageQueryParams struct {
	Offset  *int   `form:"offset" json:"offset"`
	Limit   *int   `form:"limit" json:"limit"`
	OrderBy string `form:"orderBy" json:"orderBy"`
	Order   string `form:"order" json:"order"`
	err     error
}

const (
	DefaultOffset  = 0
	DefaultLimit   = 10
	DefaultOrderBy = "id"
	DefaultOrder   = "desc"
)

func (p *PageQueryParams) SetDefault() *PageQueryParams {
	if p.Offset == nil {
		o := DefaultOffset
		p.Offset = &o
	}
	if p.Limit == nil {
		l := DefaultLimit
		p.Limit = &l
	}
	if p.OrderBy == "" {
		p.OrderBy = DefaultOrderBy
	}
	if p.Order == "" {
		p.Order = DefaultOrder
	}
	return p
}

func (p *PageQueryParams) Validate(orderByToCol map[string]string) *PageQueryParams {
	// 避免sql注入
	_, findCol := orderByToCol[p.OrderBy]
	if *p.Limit <= 0 || *p.Offset < 0 || !findCol || (p.Order != "desc" && p.Order != "asc") {
		p.err = errors.ErrBadRequest
	}
	return p
}

func (p *PageQueryParams) Error() error {
	return p.err
}

func (p *PageQueryParams) OrderByColumn(orderByToCol map[string]string) string {
	return orderByToCol[p.OrderBy]
}

func (p *PageQueryParams) IsDesc() bool {
	return p.Order == "desc"
}

func (p *PageQueryParams) GetLimit() int {
	return *p.Limit
}

func (p *PageQueryParams) GetOffset() int {
	return *p.Offset
}

func (p *PageQueryParams) GetOrder() string {
	// sql inject
	if p.Order == "desc" {
		return "desc"
	}
	return "asc"
}

const ReqVOTimeLayout = "2006-01-02 15:04:05"
