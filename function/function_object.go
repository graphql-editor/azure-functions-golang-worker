package function

import (
	"context"
	"reflect"

	"github.com/graphql-editor/azure-functions-golang-worker/api"
	"github.com/graphql-editor/azure-functions-golang-worker/converters"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/pkg/errors"
)

// Binding represents a named binding with type and direction
type Binding struct {
	Name, Type string
}

// Bindings list of bindings defined in function.json except for trigger and $return
type Bindings []Binding

type unmarshaler func(data *rpc.TypedData, v reflect.Value) error
type marshaler func(v reflect.Value) (*rpc.TypedData, error)

var (
	functionInterfaceType       = reflect.TypeOf((*api.Function)(nil)).Elem()
	returnFunctionInterfaceType = reflect.TypeOf((*api.ReturnFunction)(nil)).Elem()
	stringMapOfAny              = reflect.TypeOf((*map[string]interface{})(nil)).Elem()
)

func isStructFunctionType(t reflect.Type) (bool, error) {
	if t.Kind() != reflect.Struct {
		return false, nil
	}
	ok := t.Implements(functionInterfaceType)
	if ok {
		return false, errors.Errorf("Run method of %s must have a pointer reciever", t.Name())
	}
	return reflect.PtrTo(t).Implements(functionInterfaceType), nil
}

func isReturnStructFunctionType(t reflect.Type) (bool, error) {
	if t.Kind() != reflect.Struct {
		return false, nil
	}
	ok := t.Implements(returnFunctionInterfaceType)
	if ok {
		return false, errors.Errorf("Run method of %s must have a pointer reciever", t.Name())
	}
	return reflect.PtrTo(t).Implements(returnFunctionInterfaceType), nil
}

// TriggerType of supported triggers
type TriggerType string

const (
	// HTTPTrigger represents httpTrigger defined in function.json
	HTTPTrigger TriggerType = "httpTrigger"
)

type kind uint8

const (
	structFunction kind = iota
	returnStructFunction
	mapFunction
	invalid
)

func getFunctionType(t reflect.Type, trigger TriggerType) (reflect.Type, kind, error) {
	rt := t
	if t.Kind() == reflect.Ptr {
		rt = t.Elem()
	}
	ok, err := isStructFunctionType(rt)
	if ok || err != nil {
		return rt, structFunction, err
	}
	ok, err = isReturnStructFunctionType(rt)
	if ok {
		return rt, returnStructFunction, err
	}
	if !t.Implements(functionInterfaceType) && !t.Implements(returnFunctionInterfaceType) {
		return nil, invalid, errors.Errorf("type must implement either api.Function or api.ReturnFunction")
	}
	if rt.Kind() == reflect.Map && rt.ConvertibleTo(stringMapOfAny) {
		return rt, mapFunction, nil
	}
	return rt, invalid, errors.Errorf("unsupported function type: %s", rt.String())
}

// ObjectType represents unified interface for user function object type
type ObjectType struct {
	objectType         reflect.Type
	kind               kind
	triggerType        TriggerType
	triggerUnmarshaler unmarshaler
	returnMarshaler    marshaler
	inputUnmarshalers  map[string]unmarshaler
	outputMarshalers   map[string]marshaler
	httpOutBindings    []string
}

// NewObjectType creates new user function object type
func NewObjectType(
	t reflect.Type,
	trigger TriggerType,
	inputBindings Bindings,
	outputBindings Bindings,
) (ObjectType, error) {
	tt, kind, err := getFunctionType(t, trigger)
	if err != nil {
		return ObjectType{}, err
	}
	objectType := ObjectType{
		objectType:  tt,
		kind:        kind,
		triggerType: HTTPTrigger,
		triggerUnmarshaler: newInputUnmarshaler(Binding{
			Name: string(trigger),
		}, tt, kind),
		inputUnmarshalers: map[string]unmarshaler{},
		outputMarshalers:  map[string]marshaler{},
		httpOutBindings:   []string{},
	}
	if t.Implements(returnFunctionInterfaceType) {
		objectType.returnMarshaler = interfaceValueGet
	}
	for _, binding := range inputBindings {
		unmarshaler := newInputUnmarshaler(
			binding,
			tt,
			kind,
		)
		if unmarshaler != nil {
			objectType.inputUnmarshalers[binding.Name] = unmarshaler
		}
	}
	for _, binding := range outputBindings {
		marshaler := newOutputMarshaler(
			binding,
			tt,
			kind,
		)
		if marshaler != nil {
			objectType.outputMarshalers[binding.Name] = marshaler
		}
		if binding.Type == "http" {
			objectType.httpOutBindings = append(objectType.httpOutBindings, binding.Name)
		}
	}
	if trigger == HTTPTrigger {
		objectType.httpOutBindings = append(objectType.httpOutBindings, "$return")
	}
	return objectType, nil
}

// New creates new instance of user function object
func (f *ObjectType) New() Object {
	t := f.objectType
	instance := reflect.New(t)
	switch t.Kind() {
	case reflect.Map:
		instance.Elem().Set(reflect.MakeMap(t))
	case reflect.Slice:
		instance.Elem().Set(reflect.MakeSlice(t, 0, 0))
	}
	return Object{tp: f, instance: instance}
}

// Object represents unified interface for user function object instance
type Object struct {
	tp          *ObjectType
	instance    reflect.Value
	returnValue interface{}
}

// BindingData for user defined bindings
type BindingData struct {
	Name string
	Data *rpc.TypedData
}

// Call user function object with bound data
func (f *Object) Call(
	ctx context.Context,
	logger api.Logger,
	TriggerData *rpc.TypedData,
	TriggerMetaData map[string]*rpc.TypedData,
	inputBindings ...BindingData,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("%#v", r)
		}
	}()
	if TriggerData == nil {
		err = errors.Errorf("missing trigger data")
	}
	if err == nil && f.tp.triggerUnmarshaler != nil {
		err = f.tp.triggerUnmarshaler(TriggerData, f.instance)
	}
	for _, bd := range inputBindings {
		if err != nil {
			break
		}
		if unmarshaler, ok := f.tp.inputUnmarshalers[bd.Name]; ok {
			err = unmarshaler(bd.Data, f.instance)
		}
	}
	if len(TriggerMetaData) > 0 {
		triggerMetadata := make(map[string]interface{})
		for k, v := range TriggerMetaData {
			raw, err := converters.Unmarshal(v)
			if err == nil {
				triggerMetadata[k] = raw
			}
		}
		ctx = context.WithValue(ctx, api.TriggerMetadataKey, triggerMetadata)
	}
	if err == nil {
		switch fn := f.instance.Interface().(type) {
		case api.Function:
			fn.Run(ctx, logger)
		case api.ReturnFunction:
			f.returnValue = fn.Run(ctx, logger)
		}
	}
	return
}

// special case to allow arbitrary data to be returned through http response
func (f *Object) wrapHTTPOut(data *rpc.TypedData, name string) *rpc.TypedData {
	isHTTPOut := false
	for _, httpOut := range f.tp.httpOutBindings {
		if httpOut == name {
			isHTTPOut = true
		}
	}
	if isHTTPOut {
		_, ok := data.Data.(*rpc.TypedData_Http)
		if !ok {
			data = &rpc.TypedData{
				Data: &rpc.TypedData_Http{
					Http: &rpc.RpcHttp{
						StatusCode: "200",
						Body:       data,
					},
				},
			}
		}
	}
	return data
}

// ReturnValue returns marshaled function call return value
func (f *Object) ReturnValue() (*rpc.TypedData, bool, error) {
	_, ok := f.instance.Interface().(api.ReturnFunction)
	if !ok {
		return nil, ok, nil
	}
	td, err := f.tp.returnMarshaler(reflect.ValueOf(f.returnValue))
	if err == nil {
		td = f.wrapHTTPOut(td, "$return")
	}
	return td, ok, err
}

// GetOutput returns output binding value from user function
func (f *Object) GetOutput(name string) (*rpc.TypedData, bool, error) {
	fn, ok := f.tp.outputMarshalers[name]
	if !ok {
		return nil, ok, nil
	}
	td, err := fn(f.instance)
	if err == nil {
		td = f.wrapHTTPOut(td, name)
	}
	return td, ok, err
}
