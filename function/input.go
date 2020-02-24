package function

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/graphql-editor/azure-functions-golang-worker/converters"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/pkg/errors"
)

func getFieldValue(fieldInfo field, v reflect.Value) (reflect.Value, bool) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	field := v
	var isPtr bool
	for _, i := range fieldInfo.index {
		field = field.Field(i)
		isPtr = field.Kind() == reflect.Ptr
		if isPtr {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			field = field.Elem()
		}
	}
	return field, isPtr
}

type fieldInputUnmarshaler struct {
	field field
	set   func(*rpc.TypedData, reflect.Value) error
}

func (f *fieldInputUnmarshaler) unmarshal(data *rpc.TypedData, v reflect.Value) error {
	field, isPtr := getFieldValue(f.field, v)
	if f.field.typ.Implements(unmarshalerInterface) || (isPtr && reflect.PtrTo(f.field.typ).Implements(unmarshalerInterface)) {
		if isPtr {
			field = field.Addr()
		}
		return field.Interface().(converters.Unmarshaler).Unmarshal(data)
	}
	return f.set(data, field)
}

func stringValueSet(data *rpc.TypedData, v reflect.Value) error {
	switch td := data.Data.(type) {
	case *rpc.TypedData_String_:
		v.SetString(td.String_)
	case *rpc.TypedData_Int:
		v.SetString(strconv.FormatInt(td.Int, 10))
	case *rpc.TypedData_Double:
		v.SetString(strconv.FormatFloat(td.Double, 'f', -1, 64))
	case *rpc.TypedData_Bytes:
		v.SetString(string(td.Bytes))
	case *rpc.TypedData_Stream:
		v.SetString(string(td.Stream))
	case *rpc.TypedData_Json:
		v.SetString(td.Json)
	default:
		return errors.Errorf("unsupported typedData for string value")
	}
	return nil
}

func intValueSet(data *rpc.TypedData, v reflect.Value) error {
	switch td := data.Data.(type) {
	case *rpc.TypedData_String_:
		i, err := strconv.ParseInt(td.String_, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case *rpc.TypedData_Int:
		v.SetInt(td.Int)
	case *rpc.TypedData_Json:
		var i int64
		if err := json.Unmarshal([]byte(td.Json), &i); err != nil {
			return err
		}
		v.SetInt(i)
	default:
		return errors.Errorf("unsupported typedData for int value")
	}
	return nil
}

func uintValueSet(data *rpc.TypedData, v reflect.Value) error {
	switch td := data.Data.(type) {
	case *rpc.TypedData_String_:
		u, err := strconv.ParseUint(td.String_, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(u)
	case *rpc.TypedData_Int:
		v.SetUint(uint64(td.Int))
	case *rpc.TypedData_Json:
		var u uint64
		if err := json.Unmarshal([]byte(td.Json), &u); err != nil {
			return err
		}
		v.SetUint(u)
	default:
		return errors.Errorf("unsupported typedData for uint value")
	}
	return nil
}

func floatValueSet(data *rpc.TypedData, v reflect.Value) error {
	switch td := data.Data.(type) {
	case *rpc.TypedData_String_:
		f, err := strconv.ParseFloat(td.String_, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)
	case *rpc.TypedData_Int:
		v.SetFloat(float64(td.Int))
	case *rpc.TypedData_Json:
		var f float64
		if err := json.Unmarshal([]byte(td.Json), &f); err != nil {
			return err
		}
		v.SetFloat(f)
	case *rpc.TypedData_Double:
		v.SetFloat(td.Double)
	default:
		return errors.Errorf("unsupported typedData for float value")
	}
	return nil
}

func boolValueSet(data *rpc.TypedData, v reflect.Value) error {
	switch td := data.Data.(type) {
	case *rpc.TypedData_String_:
		b, err := strconv.ParseBool(td.String_)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case *rpc.TypedData_Int:
		v.SetBool(td.Int > 0)
	case *rpc.TypedData_Json:
		var b bool
		if err := json.Unmarshal([]byte(td.Json), &b); err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return errors.Errorf("unsupported typedData for bool value")
	}
	return nil
}

func bytesValueSet(data *rpc.TypedData, v reflect.Value) error {
	switch td := data.Data.(type) {
	case *rpc.TypedData_String_:
		v.SetBytes([]byte(td.String_))
	case *rpc.TypedData_Json:
		v.SetBytes([]byte(td.Json))
	case *rpc.TypedData_Bytes:
		v.SetBytes(td.Bytes)
	case *rpc.TypedData_Stream:
		v.SetBytes(td.Stream)
	default:
		return errors.Errorf("unsupported typedData for []byte value")
	}
	return nil
}

func mapStructSet(data *rpc.TypedData, v reflect.Value) error {
	if !v.CanAddr() {
		return errors.Errorf("cannot unmarshal into non addressable map or struct")
	}
	var jsonString []byte
	switch td := data.Data.(type) {
	case *rpc.TypedData_Bytes:
		jsonString = td.Bytes
	case *rpc.TypedData_Stream:
		jsonString = td.Stream
	case *rpc.TypedData_String_:
		jsonString = []byte(td.String_)
	case *rpc.TypedData_Json:
		jsonString = []byte(td.Json)
	default:
		return errors.Errorf("unsupported typedData %s for %s", reflect.TypeOf(data.Data).String(), v.Type().String())
	}
	return json.Unmarshal(jsonString, v.Addr().Interface())
}

func interfaceValueSet(data *rpc.TypedData, v reflect.Value) (err error) {
	var i interface{}
	switch td := data.Data.(type) {
	case *rpc.TypedData_Bytes:
		i = td.Bytes
	case *rpc.TypedData_String_:
		i = td.String_
	case *rpc.TypedData_Json:
		var m map[string]interface{}
		err = json.Unmarshal([]byte(td.Json), &m)
		if err == nil {
			i = m
		}
	case *rpc.TypedData_Int:
		i = td.Int
	case *rpc.TypedData_Double:
		i = td.Double
	case *rpc.TypedData_Stream:
		i = td.Stream
	case *rpc.TypedData_CollectionBytes:
		i = td.CollectionBytes.GetBytes()
	case *rpc.TypedData_CollectionDouble:
		i = td.CollectionDouble.GetDouble()
	case *rpc.TypedData_CollectionSint64:
		i = td.CollectionSint64.GetSint64()
	case *rpc.TypedData_CollectionString:
		i = td.CollectionString.GetString_()
	default:
		err = errors.Errorf("unsupported typedData for interface value")
	}
	if err == nil {
		v.Set(reflect.ValueOf(i))
	}
	return
}

func intSliceSet(data *rpc.TypedData, v reflect.Value) error {
	switch td := data.Data.(type) {
	case *rpc.TypedData_CollectionSint64:
		intValue := reflect.ValueOf(td.CollectionSint64.Sint64)
		if intValue.Type().ConvertibleTo(v.Type()) {
			v.Set(intValue.Convert(v.Type()))
		} else {
			n := len(td.CollectionSint64.Sint64)
			v.Set(reflect.MakeSlice(v.Type(), n, n))
			for i := 0; i < n; i++ {
				v.Index(i).SetInt(td.CollectionSint64.Sint64[i])
			}
		}
	default:
		return errors.Errorf("unsupported typedData for int slice")
	}
	return nil
}

func uintSliceSet(data *rpc.TypedData, v reflect.Value) error {
	switch td := data.Data.(type) {
	case *rpc.TypedData_CollectionSint64:
		n := len(td.CollectionSint64.Sint64)
		v.Set(reflect.MakeSlice(v.Type(), n, n))
		for i := 0; i < n; i++ {
			v.Index(i).SetUint(uint64(td.CollectionSint64.Sint64[i]))
		}
	default:
		return errors.Errorf("unsupported typedData for int slice")
	}
	return nil
}

func bytesSliceSet(data *rpc.TypedData, v reflect.Value) error {
	switch td := data.Data.(type) {
	case *rpc.TypedData_CollectionBytes:
		bytesValue := reflect.ValueOf(td.CollectionBytes.Bytes)
		if bytesValue.Type().ConvertibleTo(v.Type()) {
			v.Set(bytesValue.Convert(v.Type()))
		} else {
			n := len(td.CollectionBytes.Bytes)
			v.Set(reflect.MakeSlice(v.Type(), n, n))
			for i := 0; i < n; i++ {
				v.Index(i).Set(bytesValue.Index(i).Convert(v.Type().Elem()))
			}
		}
	default:
		return errors.Errorf("unsupported typedData for int slice")
	}
	return nil
}

func stringSliceSet(data *rpc.TypedData, v reflect.Value) error {
	switch td := data.Data.(type) {
	case *rpc.TypedData_CollectionString:
		v.Set(reflect.ValueOf(td.CollectionString.String_).Convert(v.Type()))
	default:
		return errors.Errorf("unsupported typedData for int slice")
	}
	return nil
}

func newFieldInputUnmarshaler(binding Binding, t reflect.Type) unmarshaler {
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
	unmarshaler := fieldInputUnmarshaler{
		field: *field,
	}
	switch field.typ.Kind() {
	case reflect.String:
		unmarshaler.set = stringValueSet
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		unmarshaler.set = intValueSet
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		unmarshaler.set = uintValueSet
	case reflect.Float32, reflect.Float64:
		unmarshaler.set = floatValueSet
	case reflect.Bool:
		unmarshaler.set = boolValueSet
	case reflect.Slice:
		switch field.typ.Elem().Kind() {
		case reflect.Uint8:
			unmarshaler.set = bytesValueSet
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			unmarshaler.set = intSliceSet
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			unmarshaler.set = uintSliceSet
		case reflect.String:
			unmarshaler.set = stringSliceSet
		case reflect.Slice:
			if field.typ.Elem().Elem().Kind() == reflect.Uint8 {
				unmarshaler.set = bytesSliceSet
			}
		}
	case reflect.Interface:
		if field.typ.NumMethod() == 0 {
			unmarshaler.set = interfaceValueSet
		}
	case reflect.Map, reflect.Struct:
		unmarshaler.set = mapStructSet
	}
	if unmarshaler.set == nil {
		unmarshaler.set = func(*rpc.TypedData, reflect.Value) error {
			return errors.Errorf("type %s could not be unmarshaled", field.typ.String())
		}
	}
	return unmarshaler.unmarshal
}

func newInputUnmarshaler(binding Binding, t reflect.Type, kind kind) unmarshaler {
	switch kind {
	case mapFunction:
		return mapUnmarshaler(binding)
	case structFunction, returnStructFunction:
		return newFieldInputUnmarshaler(binding, t)
	}
	return nil
}
