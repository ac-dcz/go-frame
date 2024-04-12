package geerpc

import (
	"encoding/json"
	"geerpc/codec"
	"log"
	"net"
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

}
