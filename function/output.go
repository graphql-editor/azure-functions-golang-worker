package function

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/graphql-editor/azure-functions-golang-worker/converters"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/pkg/errors"
)

var (
	unmarshalerInterface = reflect.TypeOf((*converters.Unmarshaler)(nil)).Elem()
)

func getOutputFieldValue(fieldInfo field, v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	field := v
	for _, i := range fieldInfo.index {
		field = field.Field(i)
		if !field.IsValid() || isNil(field) {
			break
		}
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
	}
	return field
}

type fieldOutputMarshaler struct {
	field field
	get   func(reflect.Value) (*rpc.TypedData, error)
}

func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func (f *fieldOutputMarshaler) marshal(v reflect.Value) (*rpc.TypedData, error) {
	field := getOutputFieldValue(f.field, v)
	if !field.IsValid() || isNil(v) {
		return nil, nil
	}
	return f.get(field)
}

func stringValueGet(v reflect.Value) (*rpc.TypedData, error) {
	return &rpc.TypedData{
		Data: &rpc.TypedData_String_{
			String_: v.String(),
		},
	}, nil
}

func intValueGet(v reflect.Value) (*rpc.TypedData, error) {
	return &rpc.TypedData{
		Data: &rpc.TypedData_Int{
			Int: v.Int(),
		},
	}, nil
}

func uintValueGet(v reflect.Value) (*rpc.TypedData, error) {
	return &rpc.TypedData{
		Data: &rpc.TypedData_Int{
			Int: int64(v.Uint()),
		},
	}, nil
}

func floatValueGet(v reflect.Value) (*rpc.TypedData, error) {
	return &rpc.TypedData{
		Data: &rpc.TypedData_Double{
			Double: v.Float(),
		},
	}, nil
}

func boolValueGet(v reflect.Value) (*rpc.TypedData, error) {
	var i int64
	if v.Bool() {
		i = 1
	}
	return &rpc.TypedData{
		Data: &rpc.TypedData_Int{
			Int: i,
		},
	}, nil
}

func bytesValueGet(v reflect.Value) (*rpc.TypedData, error) {
	return &rpc.TypedData{
		Data: &rpc.TypedData_Bytes{
			Bytes: v.Bytes(),
		},
	}, nil
}

func mapStructGet(v reflect.Value) (*rpc.TypedData, error) {
	b, err := json.Marshal(v.Interface())
	if err != nil {
		return nil, err
	}
	return &rpc.TypedData{
		Data: &rpc.TypedData_Json{
			Json: string(b),
		},
	}, nil
}

var sliceOfInt64Type = reflect.TypeOf([]int64{})

func sliceOfInt64Get(v reflect.Value) (*rpc.TypedData, error) {
	if v.Type().ConvertibleTo(sliceOfInt64Type) {
		return &rpc.TypedData{
			Data: &rpc.TypedData_CollectionSint64{
				CollectionSint64: &rpc.CollectionSInt64{
					Sint64: v.Convert(sliceOfInt64Type).Interface().([]int64),
				},
			},
		}, nil
	}
	// TODO: Long path attempting to convert types that do not directly convert to []int64
	return nil, errors.Errorf("cannot set %s a []int64", v.Type().String())
}

var sliceOfFloat64Type = reflect.TypeOf([]float64{})

func sliceOfFloat64Get(v reflect.Value) (*rpc.TypedData, error) {
	if v.Type().ConvertibleTo(sliceOfFloat64Type) {
		return &rpc.TypedData{
			Data: &rpc.TypedData_CollectionDouble{
				CollectionDouble: &rpc.CollectionDouble{
					Double: v.Convert(sliceOfFloat64Type).Interface().([]float64),
				},
			},
		}, nil
	}
	// TODO: Long path attempting to convert types that do not directly convert to []float64
	return nil, errors.Errorf("cannot set %s a []float64", v.Type().String())
}

var sliceOfStringType = reflect.TypeOf([]string{})

func sliceOfStringGet(v reflect.Value) (*rpc.TypedData, error) {
	if v.Type().ConvertibleTo(sliceOfStringType) {
		return &rpc.TypedData{
			Data: &rpc.TypedData_CollectionString{
				CollectionString: &rpc.CollectionString{
					String_: v.Convert(sliceOfStringType).Interface().([]string),
				},
			},
		}, nil
	}
	return nil, errors.Errorf("cannot set %s a []string", v.Type().String())
}

var sliceOfBytesType = reflect.TypeOf([][]byte{})

func sliceOfBytesGet(v reflect.Value) (*rpc.TypedData, error) {
	if v.Type().ConvertibleTo(sliceOfBytesType) {
		return &rpc.TypedData{
			Data: &rpc.TypedData_CollectionBytes{
				CollectionBytes: &rpc.CollectionBytes{
					Bytes: v.Convert(sliceOfBytesType).Interface().([][]byte),
				},
			},
		}, nil
	}
	// TODO: Long path attempting to convert types that do not directly convert to [][]byte
	return nil, errors.Errorf("cannot set %s a [][]byte", v.Type().String())
}

var marshalerInterface = reflect.TypeOf((*converters.Marshaler)(nil)).Elem()

func marshalerForKind(t reflect.Type) marshaler {
	if t.Implements(marshalerInterface) {
		return func(v reflect.Value) (*rpc.TypedData, error) {
			return v.Interface().(converters.Marshaler).Marshal()
		}
	}
	switch t.Kind() {
	case reflect.Interface:
		return interfaceValueGet
	case reflect.Bool:
		return boolValueGet
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intValueGet
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintValueGet
	case reflect.Float32, reflect.Float64:
		return floatValueGet
	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Uint8:
			return bytesValueGet
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return sliceOfInt64Get
		case reflect.Float32, reflect.Float64:
			return sliceOfFloat64Get
		case reflect.String:
			return sliceOfStringGet
		case reflect.Slice:
			if t.Elem().Elem().Kind() == reflect.Uint8 {
				return sliceOfBytesGet
			}
		}
	case reflect.Map, reflect.Struct:
		return mapStructGet
	case reflect.String:
		return stringValueGet
	}
	return nil
}

func interfaceValueGet(v reflect.Value) (*rpc.TypedData, error) {
	if !v.IsValid() {
		return nil, nil
	}
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if marshaler, ok := v.Interface().(converters.Marshaler); ok {
		return marshaler.Marshal()
	}
	marshaler := marshalerForKind(v.Type())
	if marshaler == nil {
		return nil, errors.Errorf("cannot find marshaler for type %s", v.Type().String())
	}
	return marshaler(v)
}

func newFieldOutputMarshaler(binding Binding, t reflect.Type) marshaler {
	fields := cachedTypeFields(t)
	var field *field
	for _, f := range fields {
		if strings.EqualFold(binding.Name, f.name) {
			field = &f
			break
		}
	}
	if field == nil {
		return nil
	}
	marshaler := fieldOutputMarshaler{
		field: *field,
		get:   marshalerForKind(field.typ),
	}
	if marshaler.get == nil {
		marshaler.get = func(reflect.Value) (*rpc.TypedData, error) {
			return nil, errors.Errorf("type %s could not be marshaled", field.typ.String())
		}
	}
	return marshaler.marshal
}

func newOutputMarshaler(binding Binding, t reflect.Type, kind kind) marshaler {
	switch kind {
	case mapFunction:
		return mapMarshaler(binding)
	case structFunction, returnStructFunction:
		return newFieldOutputMarshaler(binding, t)
	}
	return nil
}
