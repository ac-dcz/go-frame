package geerpc

import (
	"errors"
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
