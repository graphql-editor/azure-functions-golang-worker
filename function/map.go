package function

import (
	"reflect"

	"github.com/graphql-editor/azure-functions-golang-worker/converters"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
)

func mapUnmarshaler(binding Binding) func(data *rpc.TypedData, v reflect.Value) error {
	return func(data *rpc.TypedData, v reflect.Value) error {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		m := v.Convert(stringMapOfAny).Interface().(map[string]interface{})
		bindingValue, err := converters.Unmarshal(data)
		if err != nil {
			return err
		}
		m[binding.Name] = bindingValue
		return nil
	}
}

func mapMarshaler(binding Binding) func(reflect.Value) (*rpc.TypedData, error) {
	mapKey := reflect.ValueOf(binding.Name)
	return func(v reflect.Value) (*rpc.TypedData, error) {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		output := v.MapIndex(mapKey)
		for output.Kind() == reflect.Interface || output.Kind() == reflect.Ptr {
			if !output.IsValid() || isNil(output) {
				return nil, nil
			}
			output = output.Elem()
		}
		if isNil(output) {
			return nil, nil
		}
		return marshalerForKind(output.Type())(output)
	}
}
