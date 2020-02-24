package worker

import (
	"go/token"
	"strings"

	"github.com/graphql-editor/azure-functions-golang-worker/function"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/pkg/errors"
)

const returnBindingKey = "$return"

// Direction can be in, out or inout
type Direction uint8

const (
	// In represents input binding
	In Direction = iota
	// Out represents output binding
	Out
	// InOut represents binding used both as input and output
	InOut
)

// DataType is a hint for binding data type.
type DataType uint8

const (
	// Undefined data type
	Undefined DataType = iota
	// String data type
	String
	// Binary data type
	Binary
	// Stream data type
	Stream
)

// BindingInfo represents functions input and output bindings
type BindingInfo struct {
	// Type of binding (e.g. HttpTrigger)
	Type string
	// Direction of the given binding
	Direction Direction
	// DataType of binding
	DataType DataType
}

// Bindings configured in function
type Bindings map[string]BindingInfo

// FunctionInfo data
type FunctionInfo struct {
	Name               string
	Directory          string
	ScriptFile         string
	EntryPoint         string
	TriggerBindingName string
	Trigger            BindingInfo
	InputBindings      Bindings
	OutputBindings     Bindings
}

func isTrigger(typ string) bool {
	switch typ {
	case string(function.HTTPTrigger):
		return true
	default:
		return false
	}
}

func validateEntrypoint(e string) (string, error) {
	if e == "" {
		return "Function", nil
	}
	var entryPoint string
	if !token.IsExported(e) {
		entryPoint = strings.Title(e)
	}
	if !token.IsIdentifier(entryPoint) {
		return "", errors.Errorf("%s is not a valid go identifier", e)
	}
	// no need to check for keywords as all keywords start with lower case
	return entryPoint, nil
}

// NewFunctionInfo from rpc function metadata
func NewFunctionInfo(metadata *rpc.RpcFunctionMetadata) (FunctionInfo, error) {
	ep, err := validateEntrypoint(metadata.GetEntryPoint())
	if err != nil {
		return FunctionInfo{}, err
	}
	fi := FunctionInfo{
		Name:           metadata.GetName(),
		Directory:      metadata.GetDirectory(),
		ScriptFile:     metadata.GetScriptFile(),
		EntryPoint:     ep,
		InputBindings:  make(Bindings),
		OutputBindings: make(Bindings),
	}
	bindings := metadata.GetBindings()
	for k, v := range bindings {
		b := BindingInfo{
			Type: v.GetType(),
		}
		switch v.GetDataType() {
		case rpc.BindingInfo_undefined:
			b.DataType = Undefined
		case rpc.BindingInfo_string:
			b.DataType = String
		case rpc.BindingInfo_binary:
			b.DataType = Binary
		case rpc.BindingInfo_stream:
			b.DataType = Stream
		}
		if isTrigger(b.Type) {
			fi.TriggerBindingName = k
			b.Direction = In
			fi.Trigger = b
		} else {
			switch v.GetDirection() {
			case rpc.BindingInfo_inout:
				b.Direction = InOut
				fi.InputBindings[k] = b
				fi.OutputBindings[k] = b
			case rpc.BindingInfo_in:
				b.Direction = In
				fi.InputBindings[k] = b
			case rpc.BindingInfo_out:
				b.Direction = Out
				fi.OutputBindings[k] = b
			}
		}
	}
	return fi, nil
}
