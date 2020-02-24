package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"google.golang.org/grpc"
)

// Request in event stream
type Request struct {
	ID string
}

var closeSend = &rpc.StreamingMessage{}

// eventStream handler for azure functions
type eventStream struct {
	host                 string
	port                 string
	writer               chan *rpc.StreamingMessage
	writerCh             chan chan *rpc.StreamingMessage
	reader               chan *rpc.StreamingMessage
	readerCh             chan chan *rpc.StreamingMessage
	client               rpc.FunctionRpc_EventStreamClient
	clientCancel         context.CancelFunc
	lock                 sync.Mutex
	conn                 *grpc.ClientConn
	grpcMaxMessageLength int
}

func (e *eventStream) closeConn() {
	e.lock.Lock()
	if e.clientCancel != nil {
		e.clientCancel()
		e.clientCancel = nil
	}
	if e.conn != nil {
		e.conn.Close()
		e.conn = nil
	}
	e.lock.Unlock()
	wr, ok := <-e.writerCh
	if ok {
		wr <- closeSend
	}
}

func (e *eventStream) getClient() rpc.FunctionRpc_EventStreamClient {
	e.lock.Lock()
	client := e.client
	e.lock.Unlock()
	return client
}

func (e *eventStream) send(client rpc.FunctionRpc_EventStreamClient) {
	defer func() {
		close(e.writerCh)
		e.closeConn()
	}()
	var err error
	for err == nil {
		e.writerCh <- e.writer
		msg, ok := <-e.writer
		if !ok {
			return
		}
		if msg == closeSend {
			return
		}
		err = client.Send(msg)
	}
	if err != nil {
		fmt.Println(err)
	}
}

// Send a message through event stream
func (e *eventStream) Send(msg *rpc.StreamingMessage) {
	writer, ok := <-e.writerCh
	if ok {
		writer <- msg
	}
}

// Recv a message from event stream
func (e *eventStream) Recv() (*rpc.StreamingMessage, bool) {
	msg, ok := <-e.reader
	return msg, ok
}

func (e *eventStream) recv(client rpc.FunctionRpc_EventStreamClient) {
	defer func() {
		close(e.reader)
		e.closeConn()
	}()
	var err error
	for err == nil {
		var msg *rpc.StreamingMessage
		msg, err = client.Recv()
		if err == nil {
			e.reader <- msg
		}
	}
	if err != nil {
		fmt.Println(err)
	}
}

func (e *eventStream) Stop() {
	e.closeConn()
}

func (e *eventStream) Start() (err error) {
	if err == nil {
		e.lock.Lock()
		if e.clientCancel != nil {
			e.clientCancel()
		}
		if e.conn != nil {
			e.conn.Close()
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
		if e.grpcMaxMessageLength > 0 {
			opts = append(opts, grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(e.grpcMaxMessageLength),
				grpc.MaxCallSendMsgSize(e.grpcMaxMessageLength),
			))
		}
		e.conn, err = grpc.DialContext(ctx, e.host+":"+e.port, opts...)
		cancel()
		if err == nil {
			ctx, cancel := context.WithCancel(context.Background())
			e.client, err = rpc.NewFunctionRpcClient(e.conn).EventStream(ctx)
			e.clientCancel = cancel
		}
		e.lock.Unlock()
	}
	if err == nil {
		e.reader = make(chan *rpc.StreamingMessage)
		e.writer = make(chan *rpc.StreamingMessage)
		e.writerCh = make(chan chan *rpc.StreamingMessage)
		client := e.getClient()
		go e.recv(client)
		go e.send(client)
	}
	return err
}

// EventStreamOption allows setting optinal values on grpc event stream
type EventStreamOption interface {
	withEventStream(s *eventStream)
}

// MaxGrpcMessageLengthEventStreamOption sets max grpc message length on default grpc event stream
type MaxGrpcMessageLengthEventStreamOption int

func (m MaxGrpcMessageLengthEventStreamOption) withEventStream(s *eventStream) {
	s.grpcMaxMessageLength = int(m)
}

// HostEventStreamOption sets host option on default grpc event stream
type HostEventStreamOption string

func (h HostEventStreamOption) withEventStream(s *eventStream) {
	s.host = string(h)
}

// PortEventStreamOption sets port option on default grpc event stream
type PortEventStreamOption string

func (p PortEventStreamOption) withEventStream(s *eventStream) {
	s.port = string(p)
}

// NewEventStream creates new default grpc event stream with options
func NewEventStream(opts ...EventStreamOption) EventStream {
	e := &eventStream{}
	for _, opt := range opts {
		opt.withEventStream(e)
	}
	return e
}
