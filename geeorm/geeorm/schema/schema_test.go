package schema

import (
	"fmt"
	"testing"
)

type User struct {
	Id    int     `geeorm:"primary key not null"`
	Name  string  `geeorm:"not null"`
	Socre float32 `geeorm:"not null"`
}

func TestSchemaParse(t *testing.T) {
	schema, err := NewSchema(&User{})
	if err != nil {
		t.Fatal(err)
	}
	if schema.Name != "User" {
		t.Fatalf("name != User %s", schema.Name)
	}
	for i := range schema.FieldNames {
		name := schema.FieldNames[i]
		field := schema.Fields[i]
		fmt.Println(field.Name, field.Tag, field.Typ)
		if field.Name != name || schema.FieldMap[name] != field {
			t.Fatal()
		}
	}
}
