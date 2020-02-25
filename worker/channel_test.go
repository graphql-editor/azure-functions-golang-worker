package worker_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/graphql-editor/azure-functions-golang-worker/api"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/graphql-editor/azure-functions-golang-worker/worker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSender struct {
	mock.Mock
}

func (m *MockSender) Send(msg *rpc.StreamingMessage) {
	m.Called(msg)
}

type MockFuncForChannel map[string]interface{}

func (m MockFuncForChannel) Run(ctx context.Context, logger api.Logger) {
	m["copy"] = append(m["original"].([]byte), []byte("-copied")...)
	m["res"] = api.Response{
		StatusCode: http.StatusOK,
		Body:       []byte("body"),
	}
}

func TestDefaultChannel(t *testing.T) {
	mockFunctionType := reflect.TypeOf((*MockFuncForChannel)(nil)).Elem()
	var mockLoader MockTypeLoader
	mockLoader.On("GetFunctionType", mock.Anything, mock.Anything).Return(mockFunctionType, nil)
	var mockSender MockSender
	mockSender.On("Send", mock.Anything)
	ch := worker.NewChannel()
	ch.SetEventStream(&mockSender)
	ch.SetLoader(worker.Loader{
		TypeLoader:      &mockLoader,
		LoadedFunctions: map[string]worker.Function{},
	})
	ch.InitRequest("mockRequestID", &rpc.WorkerInitRequest{})
	mockSender.AssertCalled(t, "Send", &rpc.StreamingMessage{
		RequestId: "mockRequestID",
		Content: &rpc.StreamingMessage_WorkerInitResponse{
			WorkerInitResponse: &rpc.WorkerInitResponse{
				Result: &rpc.StatusResult{
					Status: rpc.StatusResult_Success,
				},
				Capabilities: map[string]string{
					"RpcHttpTriggerMetadataRemoved": "true",
					"RpcHttpBodyOnly":               "true",
					"RawHttpBodyBytes":              "true",
				},
			},
		},
	})
	ch.FunctionLoadRequest("mockRequestID", &rpc.FunctionLoadRequest{
		FunctionId: "mockFunctionID",
		Metadata: &rpc.RpcFunctionMetadata{
			Name: "func",
			Bindings: map[string]*rpc.BindingInfo{
				"trigger": &rpc.BindingInfo{
					Type:      "httpTrigger",
					Direction: rpc.BindingInfo_in,
				},
				"original": &rpc.BindingInfo{
					Type:      "blob",
					Direction: rpc.BindingInfo_in,
				},
				"copy": &rpc.BindingInfo{
					Type:      "blob",
					Direction: rpc.BindingInfo_out,
				},
				"res": &rpc.BindingInfo{
					Type:      "http",
					Direction: rpc.BindingInfo_out,
				},
			},
		},
	})
	mockSender.AssertCalled(t, "Send", &rpc.StreamingMessage{
		RequestId: "mockRequestID",
		Content: &rpc.StreamingMessage_FunctionLoadResponse{
			FunctionLoadResponse: &rpc.FunctionLoadResponse{
				FunctionId: "mockFunctionID",
				Result: &rpc.StatusResult{
					Status: rpc.StatusResult_Success,
				},
			},
		},
	})
	ch.InvocationRequest("mockRequestID", &rpc.InvocationRequest{
		FunctionId:   "mockFunctionID",
		InvocationId: "mockInvocationID",
		InputData: []*rpc.ParameterBinding{
			&rpc.ParameterBinding{
				Name: "trigger",
				Data: &rpc.TypedData{
					Data: &rpc.TypedData_Http{
						Http: &rpc.RpcHttp{},
					},
				},
			},
			&rpc.ParameterBinding{
				Name: "original",
				Data: &rpc.TypedData{
					Data: &rpc.TypedData_Bytes{
						Bytes: []byte("original"),
					},
				},
			},
		},
	})
	mockSender.AssertCalled(t, "Send", mock.MatchedBy(func(v interface{}) bool {
		msg, result := v.(*rpc.StreamingMessage)
		if result {
			var resp *rpc.StreamingMessage_InvocationResponse
			resp, result = msg.Content.(*rpc.StreamingMessage_InvocationResponse)
			if result {
				outputData := resp.InvocationResponse.OutputData
				resp.InvocationResponse.OutputData = nil
				result = result && assert.Equal(t, &rpc.StreamingMessage{
					RequestId: "mockRequestID",
					Content: &rpc.StreamingMessage_InvocationResponse{
						InvocationResponse: &rpc.InvocationResponse{
							InvocationId: "mockInvocationID",
							Result: &rpc.StatusResult{
								Status: rpc.StatusResult_Success,
							},
						},
					},
				}, v)
				// parameter bindings do not need to be sorted
				result = result && assert.Len(t, outputData, 2)
				result = result && assert.Contains(t, outputData, &rpc.ParameterBinding{
					Name: "copy",
					Data: &rpc.TypedData{
						Data: &rpc.TypedData_Bytes{
							Bytes: []byte("original-copied"),
						},
					},
				})
				result = result && assert.Contains(t, outputData, &rpc.ParameterBinding{
					Name: "res",
					Data: &rpc.TypedData{
						Data: &rpc.TypedData_Http{
							Http: &rpc.RpcHttp{
								StatusCode: "200",
								Body: &rpc.TypedData{
									Data: &rpc.TypedData_Bytes{
										Bytes: []byte("body"),
									},
								},
							},
						},
					},
				})
				resp.InvocationResponse.OutputData = outputData
			}
		}
		return result
	}))
}

func TestLoadingError(t *testing.T) {
	var mockLoader MockTypeLoader
	mockLoader.On("GetFunctionType", mock.Anything, mock.Anything).Return(nil, errors.WithStack(errors.New("some error")))
	var mockSender MockSender
	mockSender.On("Send", mock.Anything)
	ch := worker.NewChannel()
	ch.SetEventStream(&mockSender)
	ch.SetLoader(worker.Loader{
		TypeLoader:      &mockLoader,
		LoadedFunctions: map[string]worker.Function{},
	})
	ch.FunctionLoadRequest("mockRequestID", &rpc.FunctionLoadRequest{
		FunctionId: "mockFunctionID",
		Metadata: &rpc.RpcFunctionMetadata{
			Name: "func",
			Bindings: map[string]*rpc.BindingInfo{
				"trigger": &rpc.BindingInfo{
					Type:      "httpTrigger",
					Direction: rpc.BindingInfo_in,
				},
				"res": &rpc.BindingInfo{
					Type:      "http",
					Direction: rpc.BindingInfo_out,
				},
			},
		},
	})
	mockSender.AssertCalled(t, "Send", mock.MatchedBy(func(v interface{}) bool {
		result := true
		msg, ok := v.(*rpc.StreamingMessage)
		result = result && ok
		if !result {
			return result
		}
		result = result && msg.RequestId == "mockRequestID"
		resp, ok := msg.Content.(*rpc.StreamingMessage_FunctionLoadResponse)
		result = result && ok
		if !result {
			return result
		}
		result = result && "mockFunctionID" == resp.FunctionLoadResponse.FunctionId
		status := resp.FunctionLoadResponse.Result
		result = result && rpc.StatusResult_Failure == status.Status
		result = result && status.Exception != nil
		result = result && status.Exception.Message != ""
		result = result && status.Exception.StackTrace != ""
		return result
	}))
}

func TestFunctionEnvironmentReloadRequest(t *testing.T) {
	env := os.Environ()
	pwd, _ := os.Getwd()
	defer func() {
		for _, e := range env {
			var key, value string
			splitIdx := strings.Index(e, "=")
			if splitIdx != -1 {
				key = e[:splitIdx]
				value = e[splitIdx+1:]
			} else {
				key = e
			}
			os.Setenv(key, value)
		}
		os.Chdir(pwd)
	}()
	var mockSender MockSender
	mockSender.On("Send", mock.Anything)

	ch := worker.NewChannel()
	ch.SetEventStream(&mockSender)

	nwd := filepath.Join(pwd, "testdir")
	os.Mkdir(nwd, 0777)

	ch.FunctionEnvironmentReloadRequest("mockRequestID", &rpc.FunctionEnvironmentReloadRequest{
		EnvironmentVariables: map[string]string{
			"VAR": "VALUE",
		},
		FunctionAppDirectory: nwd,
	})
	mockSender.AssertCalled(t, "Send", &rpc.StreamingMessage{
		RequestId: "mockRequestID",
		Content: &rpc.StreamingMessage_FunctionEnvironmentReloadResponse{
			FunctionEnvironmentReloadResponse: &rpc.FunctionEnvironmentReloadResponse{
				Result: &rpc.StatusResult{
					Status: rpc.StatusResult_Success,
				},
			},
		},
	})
	assert.Equal(t, os.Getenv("VAR"), "VALUE")
	actualwd, _ := os.Getwd()
	assert.Equal(t, nwd, actualwd)
	os.Chdir(pwd)
	os.Remove(nwd)
}
