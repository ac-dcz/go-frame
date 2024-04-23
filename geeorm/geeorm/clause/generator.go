package main

import (
	"fmt"
	"strings"
)

type generator func(values ...any) (string, []any)

var defaultGenerator = make(map[SqlType]generator)

func init() {
	defaultGenerator[INSERT] = _insert
	defaultGenerator[VALUES] = _values
	defaultGenerator[WHERE] = _where
	defaultGenerator[ORDERBY] = _orderby
	defaultGenerator[LIMIT] = _limit
	defaultGenerator[SELECT] = _select
}

func bindVar(num int) string {
	var temp []string
	for i := 0; i < num; i++ {
		temp = append(temp, "?")
	}
	return strings.Join(temp, ",")
}

func _insert(values ...any) (string, []any) {
	//insert into [tablename] ([Fields...])
	tablename := values[0].(string)
	fields := values[1].([]string) // not values[1:]
	return fmt.Sprintf("insert into %s (%s)", tablename, strings.Join(fields, ",")), []any{}
}

func _values(values ...any) (string, []any) {
	//values (?),(?)
	sqls, vars := make([]string, 0), make([]any, 0)
	for _, val := range values {
		value := val.([]any)
		v := bindVar(len(value))
		sqls = append(sqls, fmt.Sprintf("(%s)", v))
		vars = append(vars, value...)
	}
	sql := "values " + strings.Join(sqls, ",")
	return sql, vars
}

func _where(values ...any) (string, []any) {
	//where [desc(?)]
	desc := values[0].(string)
	vars := values[1].([]any)
	return fmt.Sprintf("where %s", desc), vars
}

func _orderby(values ...any) (string, []any) {
	//order by field desc/asc
	var temp []string
	for _, val := range values {
		temp = append(temp, val.(string))
	}
	return fmt.Sprintf("order by %s", strings.Join(temp, " ")), []any{}
}

func _limit(values ...any) (string, []any) {
	//limit [?]
	return "limit ?", values
}

func _select(values ...any) (string, []any) {
	//select [fields...] from [tablename]
	fields := values[0].([]string)
	return fmt.Sprintf("select %s from %s", strings.Join(fields, " "), values[1]), []any{}
}
