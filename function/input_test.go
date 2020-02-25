package function_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/graphql-editor/azure-functions-golang-worker/api"
	functionpkg "github.com/graphql-editor/azure-functions-golang-worker/function"
	"github.com/graphql-editor/azure-functions-golang-worker/mocks"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/stretchr/testify/assert"
)

type keyType string

var testingKey = keyType("testingKey")
var expectedKey = keyType("expectedKey")

type StructType struct {
	Data string
}

type Bytes []byte

type InputTest struct {
	ByteData        []byte
	StringData      string
	IntData         int
	Int8Data        int8
	Int16Data       int16
	Int32Data       int32
	Int64Data       int64
	UintData        uint
	Uint8Data       uint8
	Uint16Data      uint16
	Uint32Data      uint32
	Uint64Data      uint64
	Float32Data     float32
	Float64Data     float64
	BoolData        bool
	PtrData         *string
	StructData      StructType
	MapData         map[string]interface{}
	InterfaceData   interface{}
	IntSliceData    []int
	Int8SliceData   []int8
	Int16SliceData  []int16
	Int32SliceData  []int32
	Int64SliceData  []int64
	UintSliceData   []uint
	Uint16SliceData []uint16
	Uint32SliceData []uint32
	Uint64SliceData []uint64
	StringSliceData []string
	BytesSliceData  [][]byte
	CustomBytes     []Bytes
}

func (i *InputTest) Run(ctx context.Context, logger api.Logger) {
	t := ctx.Value(testingKey).(*testing.T)
	expected := ctx.Value(expectedKey).(InputTest)
	assert.Equal(t, expected, *i)
}

func TestInputUnmarshaling(t *testing.T) {
	mockDataString := "mock-data"
	data := []struct {
		binding  string
		data     *rpc.TypedData
		expected InputTest
	}{
		{
			binding: "byteData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Bytes{
					Bytes: []byte(mockDataString),
				},
			},
			expected: InputTest{
				ByteData: []byte(mockDataString),
			},
		},
		{
			binding: "byteData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Stream{
					Stream: []byte(mockDataString),
				},
			},
			expected: InputTest{
				ByteData: []byte(mockDataString),
			},
		},
		{
			binding: "byteData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: mockDataString,
				},
			},
			expected: InputTest{
				ByteData: []byte(mockDataString),
			},
		},
		{
			binding: "byteData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: mockDataString,
				},
			},
			expected: InputTest{
				ByteData: []byte(mockDataString),
			},
		},
		{
			binding: "stringData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: mockDataString,
				},
			},
			expected: InputTest{
				StringData: mockDataString,
			},
		},
		{
			binding: "stringData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				StringData: "10",
			},
		},
		{
			binding: "stringData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Double{
					Double: 1.1,
				},
			},
			expected: InputTest{
				StringData: "1.1",
			},
		},
		{
			binding: "stringData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Bytes{
					Bytes: []byte(mockDataString),
				},
			},
			expected: InputTest{
				StringData: mockDataString,
			},
		},
		{
			binding: "stringData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: `"` + mockDataString + `"`,
				},
			},
			expected: InputTest{
				StringData: `"` + mockDataString + `"`,
			},
		},
		{
			binding: "intData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				IntData: 10,
			},
		},
		{
			binding: "int8Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				Int8Data: 10,
			},
		},
		{
			binding: "int16Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				Int16Data: 10,
			},
		},
		{
			binding: "int32Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				Int32Data: 10,
			},
		},
		{
			binding: "int64Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				Int64Data: 10,
			},
		},
		{
			binding: "int64Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: "10",
				},
			},
			expected: InputTest{
				Int64Data: 10,
			},
		},
		{
			binding: "int64Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: "10",
				},
			},
			expected: InputTest{
				Int64Data: 10,
			},
		},
		{
			binding: "uintData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				UintData: 10,
			},
		},
		{
			binding: "uint8Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				Uint8Data: 10,
			},
		},
		{
			binding: "uint16Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				Uint16Data: 10,
			},
		},
		{
			binding: "uint32Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				Uint32Data: 10,
			},
		},
		{
			binding: "uint64Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				Uint64Data: 10,
			},
		},
		{
			binding: "uint64Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: "10",
				},
			},
			expected: InputTest{
				Uint64Data: 10,
			},
		},
		{
			binding: "uint64Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: "10",
				},
			},
			expected: InputTest{
				Uint64Data: 10,
			},
		},
		{
			binding: "float32Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Double{
					Double: 1.0,
				},
			},
			expected: InputTest{
				Float32Data: 1.0,
			},
		},
		{
			binding: "float64Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Double{
					Double: 1.0,
				},
			},
			expected: InputTest{
				Float64Data: 1.0,
			},
		},
		{
			binding: "float64Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 1,
				},
			},
			expected: InputTest{
				Float64Data: 1.0,
			},
		},
		{
			binding: "float64Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: "1.1",
				},
			},
			expected: InputTest{
				Float64Data: 1.1,
			},
		},
		{
			binding: "float64Data",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: "1.1",
				},
			},
			expected: InputTest{
				Float64Data: 1.1,
			},
		},
		{
			binding: "boolData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 1,
				},
			},
			expected: InputTest{
				BoolData: true,
			},
		},
		{
			binding: "boolData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: "true",
				},
			},
			expected: InputTest{
				BoolData: true,
			},
		},
		{
			binding: "boolData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: "true",
				},
			},
			expected: InputTest{
				BoolData: true,
			},
		},
		{
			binding: "ptrData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: mockDataString,
				},
			},
			expected: InputTest{
				PtrData: &mockDataString,
			},
		},
		{
			binding: "structData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: `{"data": "data"}`,
				},
			},
			expected: InputTest{
				StructData: StructType{
					Data: "data",
				},
			},
		},
		{
			binding: "structData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: `{"data": "data"}`,
				},
			},
			expected: InputTest{
				StructData: StructType{
					Data: "data",
				},
			},
		},
		{
			binding: "structData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Bytes{
					Bytes: []byte(`{"data": "data"}`),
				},
			},
			expected: InputTest{
				StructData: StructType{
					Data: "data",
				},
			},
		},
		{
			binding: "structData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Stream{
					Stream: []byte(`{"data": "data"}`),
				},
			},
			expected: InputTest{
				StructData: StructType{
					Data: "data",
				},
			},
		},
		{
			binding: "mapData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: `{"data": "data"}`,
				},
			},
			expected: InputTest{
				MapData: map[string]interface{}{
					"data": "data",
				},
			},
		},
		{
			binding: "mapData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: `{"data": "data"}`,
				},
			},
			expected: InputTest{
				MapData: map[string]interface{}{
					"data": "data",
				},
			},
		},
		{
			binding: "mapData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Bytes{
					Bytes: []byte(`{"data": "data"}`),
				},
			},
			expected: InputTest{
				MapData: map[string]interface{}{
					"data": "data",
				},
			},
		},
		{
			binding: "mapData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Stream{
					Stream: []byte(`{"data": "data"}`),
				},
			},
			expected: InputTest{
				MapData: map[string]interface{}{
					"data": "data",
				},
			},
		},
		{
			binding: "interfaceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Bytes{
					Bytes: []byte(mockDataString),
				},
			},
			expected: InputTest{
				InterfaceData: []byte(mockDataString),
			},
		},
		{
			binding: "interfaceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: mockDataString,
				},
			},
			expected: InputTest{
				InterfaceData: mockDataString,
			},
		},
		{
			binding: "interfaceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: `{"data": "data"}`,
				},
			},
			expected: InputTest{
				InterfaceData: map[string]interface{}{
					"data": "data",
				},
			},
		},
		{
			binding: "interfaceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
			expected: InputTest{
				InterfaceData: int64(10),
			},
		},
		{
			binding: "interfaceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Double{
					Double: 1.0,
				},
			},
			expected: InputTest{
				InterfaceData: 1.0,
			},
		},
		{
			binding: "interfaceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_Stream{
					Stream: []byte(mockDataString),
				},
			},
			expected: InputTest{
				InterfaceData: []byte(mockDataString),
			},
		},
		{
			binding: "interfaceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionBytes{
					CollectionBytes: &rpc.CollectionBytes{
						Bytes: [][]byte{[]byte(mockDataString)},
					},
				},
			},
			expected: InputTest{
				InterfaceData: [][]byte{[]byte(mockDataString)},
			},
		},
		{
			binding: "interfaceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionDouble{
					CollectionDouble: &rpc.CollectionDouble{
						Double: []float64{1.0},
					},
				},
			},
			expected: InputTest{
				InterfaceData: []float64{1.0},
			},
		},
		{
			binding: "interfaceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(10)},
					},
				},
			},
			expected: InputTest{
				InterfaceData: []int64{int64(10)},
			},
		},
		{
			binding: "interfaceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionString{
					CollectionString: &rpc.CollectionString{
						String_: []string{mockDataString},
					},
				},
			},
			expected: InputTest{
				InterfaceData: []string{mockDataString},
			},
		},
		{
			binding: "intSliceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(10)},
					},
				},
			},
			expected: InputTest{
				IntSliceData: []int{10},
			},
		},
		{
			binding: "int8SliceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(10)},
					},
				},
			},
			expected: InputTest{
				Int8SliceData: []int8{int8(10)},
			},
		},
		{
			binding: "int16SliceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(10)},
					},
				},
			},
			expected: InputTest{
				Int16SliceData: []int16{int16(10)},
			},
		},
		{
			binding: "int32SliceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(10)},
					},
				},
			},
			expected: InputTest{
				Int32SliceData: []int32{int32(10)},
			},
		},
		{
			binding: "int64SliceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(10)},
					},
				},
			},
			expected: InputTest{
				Int64SliceData: []int64{int64(10)},
			},
		},
		{
			binding: "uintSliceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(10)},
					},
				},
			},
			expected: InputTest{
				UintSliceData: []uint{uint(10)},
			},
		},
		{
			binding: "uint16SliceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(10)},
					},
				},
			},
			expected: InputTest{
				Uint16SliceData: []uint16{uint16(10)},
			},
		},
		{
			binding: "uint32SliceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(10)},
					},
				},
			},
			expected: InputTest{
				Uint32SliceData: []uint32{uint32(10)},
			},
		},
		{
			binding: "uint64SliceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(10)},
					},
				},
			},
			expected: InputTest{
				Uint64SliceData: []uint64{uint64(10)},
			},
		},
		{
			binding: "stringSliceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionString{
					CollectionString: &rpc.CollectionString{
						String_: []string{mockDataString},
					},
				},
			},
			expected: InputTest{
				StringSliceData: []string{mockDataString},
			},
		},
		{
			binding: "bytesSliceData",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionBytes{
					CollectionBytes: &rpc.CollectionBytes{
						Bytes: [][]byte{[]byte(mockDataString)},
					},
				},
			},
			expected: InputTest{
				BytesSliceData: [][]byte{[]byte(mockDataString)},
			},
		},
		{
			binding: "customBytes",
			data: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionBytes{
					CollectionBytes: &rpc.CollectionBytes{
						Bytes: [][]byte{[]byte(mockDataString)},
					},
				},
			},
			expected: InputTest{
				CustomBytes: []Bytes{Bytes(mockDataString)},
			},
		},
	}
	ctx := context.WithValue(context.Background(), testingKey, t)
	for _, tt := range data {
		var function *InputTest
		objectType, err := functionpkg.NewObjectType(
			reflect.TypeOf(function),
			functionpkg.HTTPTrigger,
			functionpkg.Bindings{
				functionpkg.Binding{
					Name: tt.binding,
				},
			},
			functionpkg.Bindings{},
		)
		assert.NoError(t, err)
		object := objectType.New()
		assert.NoError(t, object.Call(
			context.WithValue(ctx, expectedKey, tt.expected),
			&mocks.Logger{},
			inputRPCHttpData,
			nil,
			functionpkg.BindingData{
				Name: tt.binding,
				Data: tt.data,
			},
		))
	}
}
