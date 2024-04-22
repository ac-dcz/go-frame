package schema

import (
	"fmt"
	"reflect"
)

type Field struct {
	Name string
	Typ  reflect.Type
	Tag  string
}

type Schema struct {
	Model      any
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldMap   map[string]*Field
}

func NewSchema(obj any) (*Schema, error) {
	typ := reflect.Indirect(reflect.ValueOf(obj)).Type()
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("must be struct")
	}
	schema := &Schema{
		Model:    obj,
		Name:     typ.Name(),
		FieldMap: make(map[string]*Field),
	}
	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		field := &Field{
			Name: p.Name,
			Typ:  p.Type,
		}
		if tag, ok := p.Tag.Lookup("geeorm"); ok {
			field.Tag = tag
		}
		schema.FieldNames = append(schema.FieldNames, field.Name)
		schema.Fields = append(schema.Fields, field)
		schema.FieldMap[field.Name] = field
	}
	return schema, nil
}
