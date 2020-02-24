package worker_test

import (
	"net"
	"time"

	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/graphql-editor/azure-functions-golang-worker/worker"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type mockFunctionRPCServer struct {
	*grpc.Server
	mock.Mock
	ev chan rpc.FunctionRpc_EventStreamServer
}

func (m *mockFunctionRPCServer) EventStream(ev rpc.FunctionRpc_EventStreamServer) error {
	m.ev <- ev
	var err error
	for err == nil {
		_, err = ev.Recv()
	}
	return err
}

func probeServer() error {
	conn, err := net.DialTimeout("tcp", ":1234", time.Second*5)
	if conn != nil {
		conn.Close()
	}
	return err
}

func setupMockGRPCServer(w *worker.Worker) (*mockFunctionRPCServer, rpc.FunctionRpc_EventStreamServer) {
	lis, err := net.Listen("tcp", w.Host+":"+w.Port)
	if err != nil {
		panic(err)
	}
	srv := mockFunctionRPCServer{ev: make(chan rpc.FunctionRpc_EventStreamServer)}
	srv.Server = grpc.NewServer()
	rpc.RegisterFunctionRpcServer(srv.Server, &srv)
	go srv.Serve(lis)
	retry := 0
	err = probeServer()
	for err != nil {
		if retry > 3 {
			panic("server did not start")
		}
		err = probeServer()
		retry++
	}
	go func() {
		if err := w.Listen(); err != nil {
			panic(err)
		}
	}()
	return &srv, <-srv.ev
}
