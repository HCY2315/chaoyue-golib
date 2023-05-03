package gorm

import (
	"fmt"
	"gorm.io/gorm/clause"
)

type Clauses struct {
	selectStmt clause.Select
	cls        []clause.Expression
	where      []clause.Expression
	limit      clause.Limit
}

func NewClauses(cap int) *Clauses {
	ret := Clauses{
		cls:   make([]clause.Expression, 0, cap),
		where: make([]clause.Expression, 0, cap),
	}
	return &ret
}

func (c *Clauses) Like(column, value string) *Clauses {
	c.where = append(c.where, clause.Like{
		Column: column,
		Value:  "%" + value + "%",
	})
	return c
}

func (c *Clauses) PageAndOrder(offset int, limit *int, orderByColumn string, isDesc bool) *Clauses {
	c.cls = append(c.cls, clause.OrderBy{
		Columns: []clause.OrderByColumn{
			{
				Column: clause.Column{
					Name: orderByColumn,
				},
				Desc: isDesc,
			},
		},
	})
	c.limit = clause.Limit{Limit: limit, Offset: offset}
	return c
}

func (c *Clauses) Range(columnName string, lower, upper interface{}) *Clauses {
	return c.GreatThan(columnName, lower).
		LittleThan(columnName, upper)
}

func (c *Clauses) LittleThan(column string, upper interface{}) *Clauses {
	if upper == nil {
		return c
	}
	c.where = append(c.where, clause.Lt{
		Column: column,
		Value:  upper,
	})
	return c
}

func (c *Clauses) GreatThan(column string, lower interface{}) *Clauses {
	if lower == nil {
		return c
	}
	c.where = append(c.where, clause.Gt{
		Column: column,
		Value:  lower,
	})
	return c
}

func (c Clauses) Export() []clause.Expression {
	return append(c.cls, clause.Where{Exprs: c.where}, c.selectStmt, c.limit)
}

func (c *Clauses) Equal(col string, v interface{}) *Clauses {
	c.where = append(c.where, clause.Eq{
		Column: col,
		Value:  v,
	})
	return c
}

func (c *Clauses) NotEqual(col string, v interface{}) *Clauses {
	c.where = append(c.where, clause.Neq{
		Column: col,
		Value:  v,
	})
	return c
}

func (c *Clauses) IsNull(col string) *Clauses {
	c.where = append(c.where, clause.Expr{
		SQL: col + " is NULL",
	})
	return c
}

func (c *Clauses) Exists(tableName string) *Clauses {
	if tableName == "" {
		return c.IsNotNull("deleted_at")
	}
	return c.IsNull(ColNameInTable(tableName, "deleted_at"))
}

func (c *Clauses) IsNotNull(col string) *Clauses {
	c.where = append(c.where, clause.Expr{
		SQL: col + " is not NULl",
	})
	return c
}

type tableProp [2]string

func TableProp(tableName, propName string) tableProp {
	return [2]string{tableName, propName}
}

type SelectProperties struct {
	distinct bool
	props    []tableProp
}

func (s *SelectProperties) ToSelect() clause.Select {
	cs := make([]clause.Column, 0, len(s.props))
	for _, p := range s.props {
		var col clause.Column
		col.Table, col.Name = p[0], p[1]
		cs = append(cs, col)
	}
	return clause.Select{
		Distinct:   s.distinct,
		Columns:    cs,
		Expression: nil,
	}
}

func (s *SelectProperties) AddTableProp(tableName, prop string) *SelectProperties {
	s.props = append(s.props, TableProp(tableName, prop))
	return s
}

func (s *SelectProperties) Distinct() *SelectProperties {
	s.distinct = true
	return s
}

func TableAllProperties(tableName string) *SelectProperties {
	var sp SelectProperties
	sp.AddTableProp(tableName, "*")
	return &sp
}

func TableProperties(tableName string, ps ...string) *SelectProperties {
	var sp SelectProperties
	for _, p := range ps {
		sp.AddTableProp(tableName, p)
	}
	return &sp
}

func (c *Clauses) Select(properties SelectProperties) *Clauses {
	c.selectStmt = properties.ToSelect()
	return c
}

func (c *Clauses) From(tableName string) *Clauses {
	c.cls = append(c.cls, clause.From{
		Tables: []clause.Table{
			{Name: tableName},
		},
	})
	return c
}

func tableFromName(tn string) clause.Table {
	return clause.Table{
		Name: tn,
	}
}

func ColNameInTable(tableName, colName string) string {
	return fmt.Sprintf("%s.%s", tableName, colName)
}

func onOfJoin(cfg JoinCfg) clause.Where {
	return clause.Where{
		Exprs: []clause.Expression{
			cfg,
		},
	}
}

func (c *Clauses) FromJoin(cfgs ...JoinCfg) *Clauses {
	joins := make([]clause.Join, 0, len(cfgs))
	for _, cfg := range cfgs {
		joins = append(joins, clause.Join{
			Type:  clause.InnerJoin,
			Table: tableFromName(cfg.toTN),
			ON:    onOfJoin(cfg),
		})
	}
	c.cls = append(c.cls, clause.From{
		Tables: []clause.Table{
			{
				Name: cfgs[0].fromTN,
			},
		},
		Joins: joins,
	})
	return c
}

type JoinCfg struct {
	fromTN   string
	fromProp string
	toTN     string
	toProp   string
}

func (j JoinCfg) Build(builder clause.Builder) {
	builder.WriteString(fmt.Sprintf(" `%s`.`%s`=`%s`.`%s` ",
		j.fromTN, j.fromProp,
		j.toTN, j.toProp))
}

func JoinBy(t1, p1, t2, p2 string) JoinCfg {
	return JoinCfg{
		fromTN:   t1,
		fromProp: p1,
		toTN:     t2,
		toProp:   p2,
	}
}

func (c *Clauses) In(col string, values interface{}) *Clauses {
	v := values.([]interface{})
	c.where = append(c.where, clause.IN{
		Column: clause.Column{
			Name: col,
		},
		Values: v,
	})
	return c
}

func (c *Clauses) ReplaceOffset(offset int) {
	c.limit.Offset = offset
}

func (c *Clauses) CountClauses() Clauses {
	co := *c
	co.selectStmt.Columns = []clause.Column{
		{
			Name: "count(*)",
		},
	}
	co.ReplaceOffset(0)
	return co
}
