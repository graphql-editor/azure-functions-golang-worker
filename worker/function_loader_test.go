package worker_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/graphql-editor/azure-functions-golang-worker/api"
	"github.com/graphql-editor/azure-functions-golang-worker/mocks"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/graphql-editor/azure-functions-golang-worker/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFunction map[string]interface{}

func (m MockFunction) Run(ctx context.Context, logger api.Logger) {}

func TestFunctionLoader(t *testing.T) {
	mockFunctionType := reflect.TypeOf((*MockFunction)(nil)).Elem()
	fi := worker.FunctionInfo{
		Name:               "someFunction",
		EntryPoint:         "Function",
		TriggerBindingName: "trigger",
		Trigger: worker.BindingInfo{
			Type: "httpTrigger",
		},
		InputBindings: worker.Bindings{
			"input": worker.BindingInfo{
				Type: "blob",
			},
		},
		OutputBindings: worker.Bindings{
			"output": worker.BindingInfo{
				Type:      "blob",
				Direction: worker.Out,
			},
		},
	}
	var mockLoader mocks.TypeLoader
	mockLoader.On("GetFunctionType", fi, mock.Anything).Return(mockFunctionType, nil)
	loader := worker.Loader{
		TypeLoader:      &mockLoader,
		LoadedFunctions: make(map[string]worker.Function),
	}
	assert.NoError(t, loader.Load("mockID", &rpc.RpcFunctionMetadata{
		Name: "someFunction",
		Bindings: map[string]*rpc.BindingInfo{
			"trigger": &rpc.BindingInfo{
				Type: "httpTrigger",
			},
			"input": &rpc.BindingInfo{
				Type: "blob",
			},
			"output": &rpc.BindingInfo{
				Type:      "blob",
				Direction: rpc.BindingInfo_out,
			},
		},
	}, nil))
	storedFi, err := loader.Info("mockID")
	assert.NoError(t, err)
	assert.Equal(t, fi, storedFi)
	_, err = loader.Func("mockID")
	assert.NoError(t, err)
	mockLoader.AssertCalled(t, "GetFunctionType", fi, mock.Anything)
}

func TestErrorOnBadEntrypoint(t *testing.T) {
	var mockLoader mocks.TypeLoader
	loader := worker.Loader{
		TypeLoader:      &mockLoader,
		LoadedFunctions: make(map[string]worker.Function),
	}
	assert.Error(t, loader.Load("mockID", &rpc.RpcFunctionMetadata{
		EntryPoint: "bad entrypoint",
	}, nil))
}

func TestErrorOnLoaderError(t *testing.T) {
	var mockLoader mocks.TypeLoader
	mockLoader.On("GetFunctionType", worker.FunctionInfo{
		EntryPoint:     "Function",
		InputBindings:  worker.Bindings{},
		OutputBindings: worker.Bindings{},
	}, mock.Anything).Return(nil, errors.New(""))
	loader := worker.Loader{
		TypeLoader:      &mockLoader,
		LoadedFunctions: make(map[string]worker.Function),
	}
	assert.Error(t, loader.Load("mockID", &rpc.RpcFunctionMetadata{}, nil))
}

func TestMissingFunctionErrors(t *testing.T) {
	loader := worker.Loader{
		LoadedFunctions: make(map[string]worker.Function),
	}
	_, err := loader.Info("mockID")
	assert.Error(t, err)
	_, err = loader.Func("mockID")
	assert.Error(t, err)
}
