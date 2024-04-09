package geerpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"geerpc/codec"
	"net"
	"sync"
	"sync/atomic"
)

var (
	ErrClientShutDown = errors.New("Client have shutdown already")
)

type Call struct {
	ServiceName string
	Argv        interface{}
	Replyv      interface{} //pointer
	Error       error       // 调用返回值
	Req         int         // 序号
	Done        chan *Call  // callback
}

func (c *Call) done() {
	c.Done <- <-c.Done
}

type RPCClient struct {
	opt     *Option
	cmu     sync.Mutex
	calls   map[int]*Call
	conn    net.Conn
	cc      codec.Codec
	sending sync.Mutex
	req     int
	close   atomic.Bool
}

func (client *RPCClient) Close() error {
	if client.close.Load() {
		return ErrClientShutDown
	}
	client.close.Store(true)
	return nil
}

func (client *RPCClient) IsAvaiable() bool {
	return client.close.Load() == false
}

func (client *RPCClient) registryCall(c *Call) error {
	client.cmu.Lock()
	defer client.cmu.Unlock()
	if client.close.Load() {
		return ErrClientShutDown
	}
	c.Req = client.req
	client.calls[c.Req] = c
	client.req++
	return nil
}

func (client *RPCClient) removeCall(req int) *Call {
	client.cmu.Lock()
	defer client.cmu.Unlock()
	if call, ok := client.calls[req]; ok {
		delete(client.calls, req)
		return call
	}
	return nil
}

func (client *RPCClient) terminalAllCalls(err error) {
	client.cmu.Lock()
	defer client.cmu.Unlock()
	for _, call := range client.calls {
		call.Error = err
		call.done()
	}
}

func (client *RPCClient) run() {
	var err error
	for err == nil {
		header := new(codec.Header)
		err = client.cc.ReadHeader(header)
		if err != nil {
			break
		}
		call := client.removeCall(header.Seq)
		switch {
		case call == nil:
			err = client.cc.ReadBody(nil) //丢弃Body中的内容
			call.done()
		case header.Error != "":
			call.Error = errors.New(header.Error)
			err = client.cc.ReadBody(nil) //丢弃Body中的内容
			call.done()
		default:
			if err = client.cc.ReadBody(call.Replyv); err != nil {
				call.Error = err
			}
			call.done()
		}
	}
	client.terminalAllCalls(err)
}

func parseOptions(opts ...Option) Option {
	if len(opts) == 0 {
		return DefaultOption
	} else if opts[0].SecretKey != MagicNum {
		return DefaultOption
	}
	return opts[0]
}

func NewRPCClient(network, addr string, opts ...Option) (*RPCClient, error) {
	opt := parseOptions(opts...)
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return newRPCClient(conn, opt)
}

func newRPCClient(conn net.Conn, opt Option) (*RPCClient, error) {
	codecFunc := codec.DefaultCodecFuncMap(opt.CodecType)
	if codecFunc == nil {
		return nil, fmt.Errorf("rpc client: invaild codecType %d", opt.CodecType)
	}
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		return nil, fmt.Errorf("rpc client: encode opt error %v", err)
	}
	client := &RPCClient{
		opt:     &opt,
		cmu:     sync.Mutex{},
		calls:   make(map[int]*Call),
		conn:    conn,
		cc:      codecFunc(conn),
		sending: sync.Mutex{},
		req:     0,
		close:   atomic.Bool{},
	}
	client.close.Store(false)
	go client.run()
	return client, nil
}

func (client *RPCClient) sendCall(serviceName string, argv any, replyv any) *Call {
	call := &Call{
		ServiceName: serviceName,
		Argv:        argv,
		Replyv:      replyv,
		Error:       nil,
		Done:        make(chan *Call, 1),
	}
	if err := client.registryCall(call); err != nil {
		call.Error = err
		call.done()
		return call
	}
	header := &codec.Header{
		ServiceMethod: serviceName,
		Error:         "",
		Seq:           call.Req,
	}
	if err := client.cc.Write(header, call.Argv); err != nil {
		client.removeCall(call.Req)
		call.Error = err
		call.done()
		return call
	}
	return call
}

func (client *RPCClient) Do(serviceName string, argv any, replyv any) *Call {
	return <-client.DoChan(serviceName, argv, replyv)
}

func (client *RPCClient) DoChan(serviceName string, argv any, replyv any) <-chan *Call {
	call := client.sendCall(serviceName, argv, replyv)
	return call.Done
}
