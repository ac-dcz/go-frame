package geerpc

import (
	"fmt"
	"go/ast"
	"reflect"
	"sync/atomic"
)

type methodType struct {
	method   reflect.Method
	argv     reflect.Type
	replyv   reflect.Type
	numCalls atomic.Int32
}

func (mType *methodType) NewArgv() reflect.Value {
	//所有引用类型都是Ptr吗？
	if mType.argv.Kind() == reflect.Ptr {
		return reflect.New(mType.argv.Elem())
	} else {
		return reflect.New(mType.argv).Elem()
	}
}

func (mType *methodType) NewReplyv() reflect.Value {
	value := reflect.New(mType.replyv.Elem())
	switch value.Kind() {
	case reflect.Slice:
		value.Elem().Set(reflect.MakeSlice(mType.replyv.Elem(), 0, 0))
	case reflect.Map:
		value.Elem().Set(reflect.MakeMap(mType.replyv.Elem()))
	}
	return value
}

func (mType *methodType) NumCalls() int {
	return int(mType.numCalls.Load())
}

type abcService struct {
	name    string
	selfVal reflect.Value
	selfTyp reflect.Type
	methods map[string]*methodType
}

func newAbcService(srvc any) (*abcService, error) {
	service := &abcService{
		selfVal: reflect.ValueOf(srvc),
		selfTyp: reflect.TypeOf(srvc),
		methods: make(map[string]*methodType),
	}
	service.name = reflect.Indirect(service.selfVal).Type().Name()
	if !ast.IsExported(service.name) {
		return nil, fmt.Errorf("service name %s is not exported", service.name)
	}
	service.registryMethod()
	return service, nil
}

func (srvc *abcService) registryMethod() {
	for i := 0; i < srvc.selfTyp.NumMethod(); i++ {
		method := srvc.selfTyp.Method(i)
		if method.Type.NumIn() != 3 || method.Type.NumOut() != 1 {
			continue
		}
		argv, replyv := method.Type.In(1), method.Type.In(2)
		if replyv.Kind() != reflect.Pointer {
			continue
		}
		if method.Type.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		mType := &methodType{
			method:   method,
			argv:     argv,
			replyv:   replyv,
			numCalls: atomic.Int32{},
		}
		srvc.methods[method.Name] = mType
	}
}
