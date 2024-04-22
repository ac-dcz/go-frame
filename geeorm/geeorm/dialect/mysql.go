package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type MysqlDialect struct {
}

func NewMysqlDialect() *MysqlDialect {
	return &MysqlDialect{}
}

func (d *MysqlDialect) DataTypeOf(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := reflect.New(typ).Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Name(), typ.Kind()))
}

func (d *MysqlDialect) TableExistsSql(tableName string) string {
	return fmt.Sprintf("show tables like '%s'", tableName)
}
