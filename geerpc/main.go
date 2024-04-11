package main

import (
	"fmt"
	"reflect"
)

func TestM() error {
	return nil
}

func main() {
	method := reflect.TypeOf(TestM)
	if method.Out(0) == reflect.TypeOf((*error)(nil)).Elem() {
		fmt.Println("error")
	}
}
