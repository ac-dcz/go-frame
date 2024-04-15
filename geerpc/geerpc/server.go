package geerpc

import (
	"encoding/json"
	"fmt"
	"geerpc/codec"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
)

const MagicNum = 0xbeffffeb

type Option struct {
	SecretKey uint
	CodecType codec.CodecType
}

var DefaultOption = Option{
	SecretKey: MagicNum,
	CodecType: codec.GobType,
}

var (
	ErrInvaildRequest = struct{}{}
)

type Request struct {
	H      *codec.Header
	Argv   reflect.Value
	Replyv reflect.Value
	Srvc   *abcService
	Mtype  *methodType
}

type Server struct {
	services sync.Map
}

var defaultService = Server{
	services: sync.Map{},
}

func NewServer() *Server {
	return &Server{
		services: sync.Map{},
	}
}

func Accept(lis net.Listener) error {
	return defaultService.Accept(lis)
}

func (server *Server) Accept(lis net.Listener) error {
	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		go server.handleConn(conn)
	}
}

func RegistryService(srvc any) error {
	return defaultService.RegistryService(srvc)
}

func (server *Server) RegistryService(srvc any) error {
	service, err := newAbcService(srvc)
	if err != nil {
		return err
	}
	server.services.Store(service.Name, service)
	return nil
}

func (server *Server) handleConn(conn net.Conn) {
	opt := &Option{}
	if err := json.NewDecoder(conn).Decode(opt); err != nil {
		log.Printf("rpc server: option decode error %v\n", err)
		return
	}
	if opt.SecretKey != MagicNum {
		log.Printf("rpc server: secret key invaild\n")
		return
	}
	codecFunc, ok := codec.DefaultCodecFuncMap(opt.CodecType)
	if !ok {
		log.Printf("rpc server: not found codec function\n")
		return
	}
	server.handleCodec(codecFunc(conn))
}

func (server *Server) handleCodec(cc codec.Codec) {
	sending := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			log.Printf("rpc server: %v\n", err)
			if req == nil {
				break
			}
			req.H.Error = err.Error()
			server.handleResponse(cc, req.H, ErrInvaildRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
}

func (server *Server) readRequest(cc codec.Codec) (*Request, error) {
	req := &Request{
		H: &codec.Header{},
	}
	if err := cc.ReadHeader(req.H); err != nil {
		return nil, err
	}
	srvc, mtype, err := server.findService(req.H.ServiceMethod)
	if err != nil {
		req.H.Error = err.Error()
		return req, err
	}
	req.Srvc, req.Mtype = srvc, mtype
	req.Argv, req.Replyv = mtype.NewArgv(), mtype.NewReplyv()
	inte := req.Argv.Interface()
	if req.Argv.Kind() != reflect.Pointer {
		inte = req.Argv.Addr().Interface()
	}
	if err := cc.ReadBody(inte); err != nil {
		req.H.Error = err.Error()
		return req, err
	}
	return req, nil
}

func (server *Server) findService(servcieMethod string) (*abcService, *methodType, error) {
	items := strings.Split(servcieMethod, ".")
	if len(items) != 2 {
		return nil, nil, fmt.Errorf("invaild serviceMethod %s", servcieMethod)
	}
	sName, mName := items[0], items[1]
	obj, ok := server.services.Load(sName)
	if !ok {
		return nil, nil, fmt.Errorf("invaild serviceMethod %s", servcieMethod)
	}
	srvc := obj.(*abcService)
	mtype, ok := srvc.methods[mName]
	if !ok {
		return nil, nil, fmt.Errorf("invaild serviceMethod %s", servcieMethod)
	}
	return srvc, mtype, nil
}

func (server *Server) handleRequest(cc codec.Codec, req *Request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	err := req.Srvc.Call(req.Mtype, req.Argv, req.Replyv)
	if err != nil {
		req.H.Error = err.Error()
	}
	server.handleResponse(cc, req.H, req.Replyv.Interface(), sending)
}

func (server *Server) handleResponse(cc codec.Codec, h *codec.Header, body any, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Printf("rpc server: send reponse error %v\n", err)
	}
}
