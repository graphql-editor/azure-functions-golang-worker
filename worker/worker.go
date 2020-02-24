package worker

import (
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/pkg/errors"
)

// Sender interface supports sending StreamingMessages from client to server
type Sender interface {
	Send(*rpc.StreamingMessage)
}

// EventStream interface capable of recieving and sending StreamingMessage
type EventStream interface {
	Recv() (*rpc.StreamingMessage, bool)
	Send(*rpc.StreamingMessage)
	Start() error
	Stop()
}

// Channel interface used to route streaming messages from host
type Channel interface {
	SetEventStream(Sender)
	SetLoader(Loader)
	StartStream(requestID string, msg *rpc.StartStream)
	InitRequest(requestID string, msg *rpc.WorkerInitRequest)
	Heartbeat(requestID string, msg *rpc.WorkerHeartbeat)
	Terminate(requestID string, msg *rpc.WorkerTerminate)
	StatusRequest(requestID string, msg *rpc.WorkerStatusRequest)
	FileChangeEventRequest(requestID string, msg *rpc.FileChangeEventRequest)
	FunctionLoadRequest(requestID string, msg *rpc.FunctionLoadRequest)
	InvocationRequest(requestID string, msg *rpc.InvocationRequest)
	InvocationCancel(requestID string, msg *rpc.InvocationCancel)
	FunctionEnvironmentReloadRequest(requestID string, msg *rpc.FunctionEnvironmentReloadRequest)
}

// Worker handles incoming messages from event stream
type Worker struct {
	Host      string
	Port      string
	WorkerID  string
	RequestID string
	Stream    EventStream
	Channel   Channel
	Loader    Loader
}

func (w *Worker) checkConfig() bool {
	return w.Port != ""
}

func (w *Worker) eventStream() (stream EventStream, err error) {
	stream = w.Stream
	if stream == nil {
		if !w.checkConfig() {
			err = errors.Errorf("invalid worker config")
		}
		if err == nil {
			stream = &eventStream{
				host: w.Host,
				port: w.Port,
			}
		}
	}
	return
}

func (w *Worker) getChannel(stream EventStream) Channel {
	ch := w.Channel
	if ch == nil {
		ch = NewChannel()
	}
	ch.SetEventStream(stream)
	ch.SetLoader(w.Loader)
	return ch
}

// Listen handles incoming and outgoing streaming messages from event stream
//
// It is not safe to call listen concurrently
func (w *Worker) Listen() (err error) {
	if w.WorkerID == "" || w.RequestID == "" {
		return errors.Errorf("workerID and requestID required")
	}
	stream, err := w.eventStream()
	if err == nil {
		if err = stream.Start(); err != nil {
			return
		}
		defer stream.Stop()
		stream.Send(&rpc.StreamingMessage{
			Content: &rpc.StreamingMessage_StartStream{
				StartStream: &rpc.StartStream{
					WorkerId: w.WorkerID,
				},
			},
		})
		ch := w.getChannel(stream)
		msg, ok := stream.Recv()
		for ok {
			switch msgT := msg.Content.(type) {
			case *rpc.StreamingMessage_StartStream:
				ch.StartStream(msg.RequestId, msgT.StartStream)
			case *rpc.StreamingMessage_WorkerInitRequest:
				ch.InitRequest(msg.RequestId, msgT.WorkerInitRequest)
			case *rpc.StreamingMessage_WorkerHeartbeat:
				ch.Heartbeat(msg.RequestId, msgT.WorkerHeartbeat)
			case *rpc.StreamingMessage_WorkerTerminate:
				ch.Terminate(msg.RequestId, msgT.WorkerTerminate)
			case *rpc.StreamingMessage_WorkerStatusRequest:
				ch.StatusRequest(msg.RequestId, msgT.WorkerStatusRequest)
			case *rpc.StreamingMessage_FileChangeEventRequest:
				ch.FileChangeEventRequest(msg.RequestId, msgT.FileChangeEventRequest)
			case *rpc.StreamingMessage_FunctionLoadRequest:
				ch.FunctionLoadRequest(msg.RequestId, msgT.FunctionLoadRequest)
			case *rpc.StreamingMessage_InvocationRequest:
				ch.InvocationRequest(msg.RequestId, msgT.InvocationRequest)
			case *rpc.StreamingMessage_InvocationCancel:
				ch.InvocationCancel(msg.RequestId, msgT.InvocationCancel)
			case *rpc.StreamingMessage_FunctionEnvironmentReloadRequest:
				ch.FunctionEnvironmentReloadRequest(msg.RequestId, msgT.FunctionEnvironmentReloadRequest)
			}
			msg, ok = stream.Recv()
		}
	}
	return
}
