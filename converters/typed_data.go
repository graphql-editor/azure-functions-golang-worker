package converters

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"

	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/pkg/errors"
)

// DecodeHeaders converts rpc headers to http package headers
func DecodeHeaders(h map[string]string) http.Header {
	headers := make(http.Header)
	for k, v := range h {
		headers.Add(k, v)
	}
	return headers
}

// DecodeValues converts rpc values to url package values
func DecodeValues(values map[string]string) url.Values {
	urlValues := make(url.Values)
	for k, v := range values {
		urlValues.Add(k, v)
	}
	return urlValues
}

// DecodeHTTPBody converts rpc.TypedData body to golang native
func DecodeHTTPBody(td *rpc.RpcHttp) (body, rawBody interface{}, err error) {
	if tdBody := td.GetBody(); tdBody != nil {
		body, err = Unmarshal(tdBody)
		if err != nil {
			return
		}
	}
	if tdRawBody := td.GetRawBody(); err == nil && tdRawBody != nil {
		rawBody, err = Unmarshal(tdRawBody)
	} else {
		rawBody = body
	}
	return
}

// TypedDataDecoder converts TypedData into any
type TypedDataDecoder struct{}

var marshalerInterface = reflect.TypeOf((*Marshaler)(nil)).Elem()

// Decode rpc.TypedData into anything
func (t *TypedDataDecoder) Decode(v *rpc.TypedData) (interface{}, error) {
	switch vt := v.Data.(type) {
	case *rpc.TypedData_String_:
		return vt.String_, nil
	case *rpc.TypedData_Json:
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(vt.Json), &m); err != nil {
			return nil, err
		}
		return m, nil
	case *rpc.TypedData_Bytes:
		return vt.Bytes, nil
	case *rpc.TypedData_Http:
		body, rawBody, err := DecodeHTTPBody(vt.Http)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"method":  vt.Http.GetMethod(),
			"url":     vt.Http.GetUrl(),
			"headers": DecodeHeaders(vt.Http.GetHeaders()),
			"query":   DecodeValues(vt.Http.GetQuery()),
			"params":  DecodeValues(vt.Http.GetParams()),
			"body":    body,
			"rawBody": rawBody,
		}, nil
	case *rpc.TypedData_Int:
		return vt.Int, nil
	case *rpc.TypedData_Double:
		return vt.Double, nil
	case *rpc.TypedData_CollectionBytes:
		return vt.CollectionBytes.GetBytes(), nil
	case *rpc.TypedData_CollectionString:
		return vt.CollectionString.GetString_(), nil
	case *rpc.TypedData_CollectionSint64:
		return vt.CollectionSint64.GetSint64(), nil
	case *rpc.TypedData_CollectionDouble:
		return vt.CollectionDouble.GetDouble(), nil
	}
	return nil, errors.Errorf("unsuppported typed data value")
}

// TypedDataEncoder converts any type to typed data
type TypedDataEncoder struct{}

func stringTypedData(data string) *rpc.TypedData_String_ {
	return &rpc.TypedData_String_{
		String_: data,
	}
}

func intTypedData(data interface{}) *rpc.TypedData_Int {
	var i int64
	switch dt := data.(type) {
	case int8:
		i = int64(dt)
	case uint8:
		i = int64(dt)
	case int16:
		i = int64(dt)
	case uint16:
		i = int64(dt)
	case int32:
		i = int64(dt)
	case uint32:
		i = int64(dt)
	case int:
		i = int64(dt)
	case int64:
		i = dt
	}
	return &rpc.TypedData_Int{
		Int: i,
	}
}

func doubleTypedData(data interface{}) *rpc.TypedData_Double {
	var d float64
	switch dt := data.(type) {
	case float32:
		d = float64(dt)
	case float64:
		d = dt
	}
	return &rpc.TypedData_Double{
		Double: d,
	}
}

func bytesDataType(data []byte) *rpc.TypedData_Bytes {
	return &rpc.TypedData_Bytes{
		Bytes: data,
	}
}

func collectionBytesDataType(data [][]byte) *rpc.TypedData_CollectionBytes {
	return &rpc.TypedData_CollectionBytes{
		CollectionBytes: &rpc.CollectionBytes{
			Bytes: data,
		},
	}
}

func collectionStringDataType(data []string) *rpc.TypedData_CollectionString {
	return &rpc.TypedData_CollectionString{
		CollectionString: &rpc.CollectionString{
			String_: data,
		},
	}
}

func collectionIntDataType(data interface{}) *rpc.TypedData_CollectionSint64 {
	var idata []int64
	switch dt := data.(type) {
	case []int8:
		idata = make([]int64, len(dt))
		for i, di := range dt {
			idata[i] = int64(di)
		}
	case []int16:
		idata = make([]int64, len(dt))
		for i, di := range dt {
			idata[i] = int64(di)
		}
	case []uint16:
		idata = make([]int64, len(dt))
		for i, di := range dt {
			idata[i] = int64(di)
		}
	case []int32:
		idata = make([]int64, len(dt))
		for i, di := range dt {
			idata[i] = int64(di)
		}
	case []uint32:
		idata = make([]int64, len(dt))
		for i, di := range dt {
			idata[i] = int64(di)
		}
	case []int:
		idata = make([]int64, len(dt))
		for i, di := range dt {
			idata[i] = int64(di)
		}
	case []int64:
		idata = dt
	}
	return &rpc.TypedData_CollectionSint64{
		CollectionSint64: &rpc.CollectionSInt64{
			Sint64: idata,
		},
	}
}

func collectionDoubleDataType(data interface{}) *rpc.TypedData_CollectionDouble {
	var ddata []float64
	switch dt := data.(type) {
	case []float32:
		ddata = make([]float64, len(dt))
		for i, dd := range dt {
			ddata[i] = float64(dd)
		}
	case []float64:
		ddata = dt
	}
	return &rpc.TypedData_CollectionDouble{
		CollectionDouble: &rpc.CollectionDouble{
			Double: ddata,
		},
	}
}

var (
	stringType = reflect.ValueOf("").Type()
	byteType   = reflect.ValueOf(uint8(0)).Type()
	bytesType  = reflect.ValueOf([]byte{}).Type()
	intType    = reflect.ValueOf(int64(0)).Type()
	doubleType = reflect.ValueOf(float64(0)).Type()
)

func valueAsBytes(v reflect.Value) []byte {
	slice := v
	if v.Kind() == reflect.Array {
		slice = reflect.MakeSlice(
			reflect.SliceOf(v.Type().Elem()),
			v.Len(),
			v.Len(),
		)
		reflect.Copy(slice, v)
	}
	return slice.Bytes()
}

func sliceValueSet(dst, src reflect.Value) {
	switch dst.Kind() {
	case reflect.String:
		dst.SetString(src.String())
	case reflect.Int64:
		switch src.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dst.SetInt(src.Int())
		case reflect.Uint16, reflect.Uint32:
			dst.SetInt(int64(src.Uint()))
		}
	case reflect.Float64:
		dst.SetFloat(src.Float())
	case reflect.Slice:
		dst.SetBytes(valueAsBytes(src))
	}
}

func copySlice(slice reflect.Value, ofType reflect.Type) reflect.Value {
	newSlice := reflect.MakeSlice(reflect.SliceOf(ofType), slice.Len(), slice.Len())
	if slice.Type().Elem() == ofType {
		// short path
		reflect.Copy(newSlice, slice)
	} else {
		// long path
		for i := 0; i < slice.Len(); i++ {
			sliceValueSet(newSlice.Index(i), slice.Index(i))
		}
	}
	return newSlice
}

// Encode any data into rpc.TypedData
func (t *TypedDataEncoder) Encode(v interface{}) (td *rpc.TypedData, err error) {
	var ltd rpc.TypedData
	// short path for known types
	switch vt := v.(type) {
	case string:
		ltd.Data = stringTypedData(vt)
	case []byte:
		ltd.Data = bytesDataType(vt)
	case *rpc.RpcHttp:
		ltd.Data = &rpc.TypedData_Http{
			Http: vt,
		}
	case int, int8, int16, int32, int64, uint8, uint16, uint32:
		ltd.Data = intTypedData(vt)
	case float32, float64:
		ltd.Data = doubleTypedData(vt)
	case [][]byte:
		ltd.Data = collectionBytesDataType(vt)
	case []string:
		ltd.Data = collectionStringDataType(vt)
	case []int, []int8, []int16, []int32, []int64, []uint16, []uint32:
		ltd.Data = collectionIntDataType(vt)
	case []float32, []float64:
		ltd.Data = collectionDoubleDataType(vt)
	}
	if ltd.Data == nil && err == nil {
		// Long path
		// Use reflection to handle similar types, like type definitions using strings,
		// different sizes of integers and floating points etc.

		vr := reflect.ValueOf(v)
		if vr.IsValid() {
			if vr.Type().Implements(marshalerInterface) {
				return v.(Marshaler).Marshal()
			}
			if vr.Kind() == reflect.Ptr || vr.Kind() == reflect.Interface {
				if !vr.IsNil() {
					vr = vr.Elem()
				}
			}

			switch vr.Kind() {
			case reflect.String:
				ltd.Data = stringTypedData(vr.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ltd.Data = intTypedData(vr.Int())
			case reflect.Uint8, reflect.Uint16, reflect.Uint32:
				ltd.Data = intTypedData(int64(vr.Uint()))
			case reflect.Float32, reflect.Float64:
				ltd.Data = doubleTypedData(vr.Float())
			case reflect.Slice, reflect.Array:
				elem := vr.Type().Elem()
				switch elem.Kind() {
				case reflect.Uint8:
					// Bytes
					ltd.Data = bytesDataType(valueAsBytes(vr))
				case reflect.String:
					// CollectionString
					// convert array/slice of strings and similar to a slice of strings
					ltd.Data = collectionStringDataType(copySlice(
						vr,
						stringType,
					).Interface().([]string))
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					// CollectionSInt64
					// convert array/slice of ints and similar to a slice of int64
					ltd.Data = collectionIntDataType(copySlice(
						vr,
						intType,
					).Interface().([]int64))
				case reflect.Uint16, reflect.Uint32:
					// CollectionSInt64
					// convert array/slice of unsigned and similar to a slice of int64
					ltd.Data = collectionIntDataType(copySlice(
						vr,
						intType,
					).Interface().([]int64))
				case reflect.Float32, reflect.Float64:
					// CollectionDouble
					// convert array/slice of floating points and similar to a slice of float64
					ltd.Data = collectionDoubleDataType(copySlice(
						vr,
						doubleType,
					).Interface().([]float64))
				case reflect.Slice, reflect.Array:
					if elem.Elem().Kind() == reflect.Uint8 {
						// CollectionBytes
						// convert array/slice of bytes and similar to a slice of []byte
						ltd.Data = collectionBytesDataType(copySlice(
							vr,
							bytesType,
						).Interface().([][]byte))
					}
				}
			}
		}
		// Fallback to JSON for anything else
		if ltd.Data == nil && err == nil {
			var b []byte
			b, err = json.Marshal(v)
			if err == nil {
				ltd.Data = &rpc.TypedData_Json{
					Json: string(b),
				}
			}
		}
	}
	if err == nil {
		td = &ltd
	}
	return
}
