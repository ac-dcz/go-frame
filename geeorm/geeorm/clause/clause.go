package main

import "strings"

type SqlType int

const (
	SELECT  SqlType = iota // select [fields...] from [tablename]
	WHERE                  // where [desc...]
	ORDERBY                // order by %s [desc/asc]
	INSERT                 // insert into [tablename...](fields...)
	LIMIT                  // limit [num]
	VALUES                 // values (),(),()
)

type Clause struct {
	sql  map[SqlType]string
	vars map[SqlType][]any
}

func NewClause() *Clause {
	return &Clause{
		sql:  make(map[SqlType]string),
		vars: make(map[SqlType][]any),
	}
}

func (c *Clause) Set(typ SqlType, vars ...any) *Clause {
	if gen, ok := defaultGenerator[typ]; ok {
		ts, tv := gen(vars...)
		c.sql[typ] = ts
		c.vars[typ] = tv
	}
	return c
}

func (c *Clause) Build(typs ...SqlType) (string, []any) {
	sqls, vars := make([]string, 0), make([]any, 0)
	for _, typ := range typs {
		sqls = append(sqls, c.sql[typ])
		vars = append(vars, c.vars[typ]...)
	}
	return strings.Join(sqls, " "), vars
}
