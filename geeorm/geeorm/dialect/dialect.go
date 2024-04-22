package dialect

import "reflect"

type Dialect interface {
	DataTypeOf(value reflect.Type) string
	TableExistsSql(tablename string) string
}

var dialectMap map[string]Dialect

func init() {
	dialectMap = make(map[string]Dialect)
	dialectMap["mysql"] = NewMysqlDialect()
}

func RegistryDialect(driver string, d Dialect) {
	dialectMap[driver] = d
}

func GetDialect(driver string) (Dialect, bool) {
	d, ok := dialectMap[driver]
	return d, ok
}
