package worker_test

import (
	"testing"

	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/graphql-editor/azure-functions-golang-worker/worker"
	"github.com/stretchr/testify/assert"
)

func TestNewFunctionInfo(t *testing.T) {
	fi, err := worker.NewFunctionInfo(
		&rpc.RpcFunctionMetadata{
			Name:       "function",
			Directory:  "mock/path",
			ScriptFile: "main.go",
			Bindings: map[string]*rpc.BindingInfo{
				"trigger": &rpc.BindingInfo{
					Type:      "httpTrigger",
					Direction: rpc.BindingInfo_in,
					DataType:  rpc.BindingInfo_binary,
				},
				"input": &rpc.BindingInfo{
					Type:      "blob",
					Direction: rpc.BindingInfo_in,
					DataType:  rpc.BindingInfo_string,
				},
				"output": &rpc.BindingInfo{
					Type:      "blob",
					Direction: rpc.BindingInfo_out,
					DataType:  rpc.BindingInfo_stream,
				},
				"inputoutput": &rpc.BindingInfo{
					Type:      "blob",
					Direction: rpc.BindingInfo_inout,
					DataType:  rpc.BindingInfo_undefined,
				},
				"$return": &rpc.BindingInfo{
					Type:      "http",
					Direction: rpc.BindingInfo_out,
					DataType:  rpc.BindingInfo_undefined,
				},
			},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, worker.FunctionInfo{
		Name:               "function",
		Directory:          "mock/path",
		ScriptFile:         "main.go",
		EntryPoint:         "Function",
		TriggerBindingName: "trigger",
		Trigger: worker.BindingInfo{
			Type:      "httpTrigger",
			Direction: worker.In,
			DataType:  worker.Binary,
		},
		InputBindings: worker.Bindings{
			"input": {
				Type:      "blob",
				Direction: worker.In,
				DataType:  worker.String,
			},
			"inputoutput": {
				Type:      "blob",
				Direction: worker.InOut,
				DataType:  worker.Undefined,
			},
		},
		OutputBindings: worker.Bindings{
			"output": {
				Type:      "blob",
				Direction: worker.Out,
				DataType:  worker.Stream,
			},
			"inputoutput": {
				Type:      "blob",
				Direction: worker.InOut,
				DataType:  worker.Undefined,
			},
			"$return": {
				Type:      "http",
				Direction: worker.Out,
				DataType:  worker.Undefined,
			},
		},
	}, fi)
}

func TestEntrypointValidation(t *testing.T) {
	fi, err := worker.NewFunctionInfo(&rpc.RpcFunctionMetadata{
		EntryPoint: "entryPoint",
	})
	assert.NoError(t, err)
	assert.Equal(t, "EntryPoint", fi.EntryPoint)
	_, err = worker.NewFunctionInfo(&rpc.RpcFunctionMetadata{
		EntryPoint: "not a valid entry point",
	})
	assert.Error(t, err)
}
