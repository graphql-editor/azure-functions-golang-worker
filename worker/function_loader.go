package worker

import (
	"reflect"

	"github.com/graphql-editor/azure-functions-golang-worker/api"
	"github.com/graphql-editor/azure-functions-golang-worker/function"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/pkg/errors"
)

// TypeLoader loads type that represents function
// Type representing function must implement api.Function or api.ReturnFunction.
type TypeLoader interface {
	GetFunctionType(FunctionInfo, api.Logger) (reflect.Type, error)
}

// Function is loaded function in worker
type Function struct {
	Info       FunctionInfo
	ObjectType function.ObjectType
}

// Loader loads function object for given function id.
type Loader struct {
	TypeLoader
	LoadedFunctions map[string]Function
}

// Info returns function info for given function id
func (l *Loader) Info(functionID string) (FunctionInfo, error) {
	f, ok := l.LoadedFunctions[functionID]
	if !ok {
		return FunctionInfo{}, errors.Errorf("function info for %s is not loaded and cannot be invoked", functionID)
	}
	return f.Info, nil
}

// Func returns function object type to execute
func (l *Loader) Func(functionID string) (function.ObjectType, error) {
	f, ok := l.LoadedFunctions[functionID]
	if !ok {
		return function.ObjectType{}, errors.Errorf("function code for %s is not loaded and cannot be invoked", functionID)
	}
	return f.ObjectType, nil
}

// Load returns object type for given function id
func (l *Loader) Load(functionID string, metadata *rpc.RpcFunctionMetadata, logger api.Logger) error {
	info, err := NewFunctionInfo(metadata)
	if err != nil {
		return err
	}
	t, err := l.GetFunctionType(info, logger)
	if err != nil {
		return err
	}
	inputBindings := make(function.Bindings, 0, len(info.InputBindings))
	for k, v := range info.InputBindings {
		inputBindings = append(inputBindings, function.Binding{
			Name: k,
			Type: v.Type,
		})
	}
	outputBindings := make(function.Bindings, 0, len(info.OutputBindings))
	for k, v := range info.OutputBindings {
		outputBindings = append(outputBindings, function.Binding{
			Name: k,
			Type: v.Type,
		})
	}
	ot, err := function.NewObjectType(
		t,
		function.TriggerType(info.Trigger.Type),
		inputBindings,
		outputBindings,
	)
	if err == nil {
		l.LoadedFunctions[functionID] = Function{
			Info:       info,
			ObjectType: ot,
		}
	}
	return err
}
