package worker_test

import (
	"strings"
	"testing"

	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/graphql-editor/azure-functions-golang-worker/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogger(t *testing.T) {
	var mockSender MockSender
	mockSender.On("Send", mock.Anything)

	logger := worker.Logger{
		InvocationID: "mockInvocationID",
		EventID:      "mockEventID",
		Cat:          rpc.RpcLog_User,
		Stream:       &mockSender,
	}
	logger.Trace("trace msg")
	mockSender.AssertCalled(t, "Send", mock.MatchedBy(func(v interface{}) bool {
		msg, result := v.(*rpc.StreamingMessage)
		if result {
			var content *rpc.StreamingMessage_RpcLog
			content, result = msg.Content.(*rpc.StreamingMessage_RpcLog)
			if result {
				msg := content.RpcLog.Message
				content.RpcLog.Message = ""
				result = result && assert.Equal(t, &rpc.StreamingMessage{
					Content: &rpc.StreamingMessage_RpcLog{
						RpcLog: &rpc.RpcLog{
							InvocationId: "mockInvocationID",
							EventId:      "mockEventID",
							Level:        rpc.RpcLog_Trace,
							Category:     "User",
						},
					},
				}, v)
				traceLines := strings.Split(msg, "\n")
				assert.Equal(t, "trace msg", traceLines[0])
				assert.Equal(t, "github.com/graphql-editor/azure-functions-golang-worker/worker_test.TestLogger", traceLines[2])
				content.RpcLog.Message = msg
			}
		}
		return result
	}))
	logger.Debug("debug msg")
	mockSender.AssertCalled(t, "Send", &rpc.StreamingMessage{
		Content: &rpc.StreamingMessage_RpcLog{
			RpcLog: &rpc.RpcLog{
				InvocationId: "mockInvocationID",
				EventId:      "mockEventID",
				Message:      "debug msg",
				Level:        rpc.RpcLog_Debug,
				Category:     "User",
			},
		},
	})
	logger.Info("info msg")
	mockSender.AssertCalled(t, "Send", &rpc.StreamingMessage{
		Content: &rpc.StreamingMessage_RpcLog{
			RpcLog: &rpc.RpcLog{
				InvocationId: "mockInvocationID",
				EventId:      "mockEventID",
				Message:      "info msg",
				Level:        rpc.RpcLog_Information,
				Category:     "User",
			},
		},
	})
	logger.Warn("warn msg")
	mockSender.AssertCalled(t, "Send", &rpc.StreamingMessage{
		Content: &rpc.StreamingMessage_RpcLog{
			RpcLog: &rpc.RpcLog{
				InvocationId: "mockInvocationID",
				EventId:      "mockEventID",
				Message:      "warn msg",
				Level:        rpc.RpcLog_Warning,
				Category:     "User",
			},
		},
	})
	logger.Error("error msg")
	mockSender.AssertCalled(t, "Send", &rpc.StreamingMessage{
		Content: &rpc.StreamingMessage_RpcLog{
			RpcLog: &rpc.RpcLog{
				InvocationId: "mockInvocationID",
				EventId:      "mockEventID",
				Message:      "error msg",
				Level:        rpc.RpcLog_Error,
				Category:     "User",
			},
		},
	})
	logger.Fatal("fatal msg")
	mockSender.AssertCalled(t, "Send", &rpc.StreamingMessage{
		Content: &rpc.StreamingMessage_RpcLog{
			RpcLog: &rpc.RpcLog{
				InvocationId: "mockInvocationID",
				EventId:      "mockEventID",
				Message:      "fatal msg",
				Level:        rpc.RpcLog_Critical,
				Category:     "User",
			},
		},
	})
}
