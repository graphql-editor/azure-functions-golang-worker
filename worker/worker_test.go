package worker_test

import (
	"sync"
	"testing"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/graphql-editor/azure-functions-golang-worker/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChannel struct {
	mock.Mock
	wg sync.WaitGroup
}

func (m *MockChannel) StartStream(requestID string, msg *rpc.StartStream) {
	defer m.wg.Done()
	m.Called(requestID, msg)
}

func (m *MockChannel) InitRequest(requestID string, msg *rpc.WorkerInitRequest) {
	defer m.wg.Done()
	m.Called(requestID, msg)
}

func (m *MockChannel) Heartbeat(requestID string, msg *rpc.WorkerHeartbeat) {
	defer m.wg.Done()
	m.Called(requestID, msg)
}

func (m *MockChannel) Terminate(requestID string, msg *rpc.WorkerTerminate) {
	defer m.wg.Done()
	m.Called(requestID, msg)
}

func (m *MockChannel) StatusRequest(requestID string, msg *rpc.WorkerStatusRequest) {
	defer m.wg.Done()
	m.Called(requestID, msg)
}

func (m *MockChannel) FileChangeEventRequest(requestID string, msg *rpc.FileChangeEventRequest) {
	defer m.wg.Done()
	m.Called(requestID, msg)
}

func (m *MockChannel) FunctionLoadRequest(requestID string, msg *rpc.FunctionLoadRequest) {
	defer m.wg.Done()
	m.Called(requestID, msg)
}

func (m *MockChannel) InvocationRequest(requestID string, msg *rpc.InvocationRequest) {
	defer m.wg.Done()
	m.Called(requestID, msg)
}

func (m *MockChannel) InvocationCancel(requestID string, msg *rpc.InvocationCancel) {
	defer m.wg.Done()
	m.Called(requestID, msg)
}

func (m *MockChannel) FunctionEnvironmentReloadRequest(requestID string, msg *rpc.FunctionEnvironmentReloadRequest) {
	defer m.wg.Done()
	m.Called(requestID, msg)
}

func (m *MockChannel) SetEventStream(stream worker.Sender) {
	m.Called(stream)
}

func (m *MockChannel) SetLoader(loader worker.Loader) {
	m.Called(loader)
}

func TestWorkerCallsChannelStartStream(t *testing.T) {
	data := []struct {
		function string
		expected interface{}
		msg      *rpc.StreamingMessage
	}{
		{
			function: "StartStream",
			expected: &rpc.StartStream{
				WorkerId: "mockWorkerId",
			},
			msg: &rpc.StreamingMessage{
				RequestId: "mockRequestId",
				Content: &rpc.StreamingMessage_StartStream{
					StartStream: &rpc.StartStream{
						WorkerId: "mockWorkerId",
					},
				},
			},
		},
		{
			function: "InitRequest",
			expected: &rpc.WorkerInitRequest{
				HostVersion: "mockHostVersion",
			},
			msg: &rpc.StreamingMessage{
				RequestId: "mockRequestId",
				Content: &rpc.StreamingMessage_WorkerInitRequest{
					WorkerInitRequest: &rpc.WorkerInitRequest{
						HostVersion: "mockHostVersion",
					},
				},
			},
		},
		{
			function: "Heartbeat",
			expected: &rpc.WorkerHeartbeat{},
			msg: &rpc.StreamingMessage{
				RequestId: "mockRequestId",
				Content: &rpc.StreamingMessage_WorkerHeartbeat{
					WorkerHeartbeat: &rpc.WorkerHeartbeat{},
				},
			},
		},
		{
			function: "Terminate",
			expected: &rpc.WorkerTerminate{
				GracePeriod: &duration.Duration{
					Seconds: 1,
				},
			},
			msg: &rpc.StreamingMessage{
				RequestId: "mockRequestId",
				Content: &rpc.StreamingMessage_WorkerTerminate{
					WorkerTerminate: &rpc.WorkerTerminate{
						GracePeriod: &duration.Duration{
							Seconds: 1,
						},
					},
				},
			},
		},
		{
			function: "StatusRequest",
			expected: &rpc.WorkerStatusRequest{},
			msg: &rpc.StreamingMessage{
				RequestId: "mockRequestId",
				Content: &rpc.StreamingMessage_WorkerStatusRequest{
					WorkerStatusRequest: &rpc.WorkerStatusRequest{},
				},
			},
		},
		{
			function: "FileChangeEventRequest",
			expected: &rpc.FileChangeEventRequest{
				Name: "mockName",
			},
			msg: &rpc.StreamingMessage{
				RequestId: "mockRequestId",
				Content: &rpc.StreamingMessage_FileChangeEventRequest{
					FileChangeEventRequest: &rpc.FileChangeEventRequest{
						Name: "mockName",
					},
				},
			},
		},
		{
			function: "FunctionLoadRequest",
			expected: &rpc.FunctionLoadRequest{
				FunctionId: "mockFunctionId",
			},
			msg: &rpc.StreamingMessage{
				RequestId: "mockRequestId",
				Content: &rpc.StreamingMessage_FunctionLoadRequest{
					FunctionLoadRequest: &rpc.FunctionLoadRequest{
						FunctionId: "mockFunctionId",
					},
				},
			},
		},
		{
			function: "InvocationRequest",
			expected: &rpc.InvocationRequest{
				FunctionId: "mockFunctionId",
			},
			msg: &rpc.StreamingMessage{
				RequestId: "mockRequestId",
				Content: &rpc.StreamingMessage_InvocationRequest{
					InvocationRequest: &rpc.InvocationRequest{
						FunctionId: "mockFunctionId",
					},
				},
			},
		},
		{
			function: "InvocationCancel",
			expected: &rpc.InvocationCancel{
				InvocationId: "mockInvocationId",
			},
			msg: &rpc.StreamingMessage{
				RequestId: "mockRequestId",
				Content: &rpc.StreamingMessage_InvocationCancel{
					InvocationCancel: &rpc.InvocationCancel{
						InvocationId: "mockInvocationId",
					},
				},
			},
		},
		{
			function: "FunctionEnvironmentReloadRequest",
			expected: &rpc.FunctionEnvironmentReloadRequest{
				EnvironmentVariables: map[string]string{
					"mock": "mock",
				},
			},
			msg: &rpc.StreamingMessage{
				RequestId: "mockRequestId",
				Content: &rpc.StreamingMessage_FunctionEnvironmentReloadRequest{
					FunctionEnvironmentReloadRequest: &rpc.FunctionEnvironmentReloadRequest{
						EnvironmentVariables: map[string]string{
							"mock": "mock",
						},
					},
				},
			},
		},
	}
	var mockChannel MockChannel
	mockChannel.On("SetEventStream", mock.Anything)
	mockChannel.On("SetLoader", mock.Anything)
	mockChannel.wg.Add(len(data))
	worker := worker.Worker{
		Channel:   &mockChannel,
		WorkerID:  "mockWorkerID",
		RequestID: "mockRequestID",
		Port:      "1234",
	}
	srv, eventStream := setupMockGRPCServer(&worker)
	defer srv.Stop()
	matcher := func(expected interface{}) func(v interface{}) bool {
		return func(v interface{}) bool {
			return assert.Equal(t, expected, v)
		}
	}
	for _, tt := range data {
		mockChannel.On(tt.function, tt.msg.RequestId, mock.MatchedBy(matcher(tt.expected)))
		eventStream.Send(tt.msg)
	}
	mockChannel.wg.Wait()
	for _, tt := range data {
		mockChannel.AssertCalled(t, tt.function, tt.msg.RequestId, mock.MatchedBy(matcher(tt.expected)))
	}
}
