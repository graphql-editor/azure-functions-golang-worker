package function_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/graphql-editor/azure-functions-golang-worker/api"
	functionpkg "github.com/graphql-editor/azure-functions-golang-worker/function"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/stretchr/testify/assert"
)

type OutputTest struct {
	BytesValue        []byte
	StringValue       string
	IntValue          int
	Int8Value         int8
	Int16Value        int16
	Int32Value        int32
	Int64Value        int64
	UintValue         uint
	Uint8Value        uint8
	Uint16Value       uint16
	Uint32Value       uint32
	Uint64Value       uint64
	Float32Value      float32
	Float64Value      float64
	BoolValue         bool
	PtrValue          *string
	StructValue       StructType
	MapValue          map[string]interface{}
	Int64SliceValue   []int64
	Float64SliceValue []float64
	StringSliceValue  []string
	BytesSliceValue   [][]byte
}

func (o *OutputTest) Run(ctx context.Context, logger api.Logger) {
	expected := ctx.Value(expectedKey).(OutputTest)
	*o = expected
}

func TestOutputMarshal(t *testing.T) {
	mockDataString := "mock-data"
	data := []struct {
		binding  string
		output   OutputTest
		expected *rpc.TypedData
	}{
		{
			binding: "bytesValue",
			output:  OutputTest{BytesValue: []byte(mockDataString)},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Bytes{
					Bytes: []byte(mockDataString),
				},
			},
		},
		{
			binding: "stringValue",
			output:  OutputTest{StringValue: mockDataString},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: mockDataString,
				},
			},
		},
		{
			binding: "intValue",
			output:  OutputTest{IntValue: 10},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
		},
		{
			binding: "int8Value",
			output:  OutputTest{Int8Value: 10},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
		},
		{
			binding: "int16Value",
			output:  OutputTest{Int16Value: 10},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
		},
		{
			binding: "int32Value",
			output:  OutputTest{Int32Value: 10},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
		},
		{
			binding: "int64Value",
			output:  OutputTest{Int64Value: 10},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
		},
		{
			binding: "uintValue",
			output:  OutputTest{UintValue: 10},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
		},
		{
			binding: "uint8Value",
			output:  OutputTest{Uint8Value: 10},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
		},
		{
			binding: "uint16Value",
			output:  OutputTest{Uint16Value: 10},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
		},
		{
			binding: "uint32Value",
			output:  OutputTest{Uint32Value: 10},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
		},
		{
			binding: "uint64Value",
			output:  OutputTest{Uint64Value: 10},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 10,
				},
			},
		},
		{
			binding: "float32Value",
			output:  OutputTest{Float32Value: 1.0},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Double{
					Double: 1.0,
				},
			},
		},
		{
			binding: "float64Value",
			output:  OutputTest{Float64Value: 1.0},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Double{
					Double: 1.0,
				},
			},
		},
		{
			binding: "boolValue",
			output:  OutputTest{BoolValue: true},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: 1,
				},
			},
		},
		{
			binding: "ptrValue",
			output:  OutputTest{PtrValue: &mockDataString},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: mockDataString,
				},
			},
		},
		{
			binding: "structValue",
			output:  OutputTest{StructValue: StructType{Data: "data"}},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: `{"Data":"data"}`,
				},
			},
		},
		{
			binding: "mapValue",
			output:  OutputTest{MapValue: map[string]interface{}{"data": "data"}},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: `{"data":"data"}`,
				},
			},
		},
		{
			binding: "int64SliceValue",
			output:  OutputTest{Int64SliceValue: []int64{10}},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{10},
					},
				},
			},
		},
		{
			binding: "float64SliceValue",
			output:  OutputTest{Float64SliceValue: []float64{1.0}},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionDouble{
					CollectionDouble: &rpc.CollectionDouble{
						Double: []float64{1.0},
					},
				},
			},
		},
		{
			binding: "stringSliceValue",
			output:  OutputTest{StringSliceValue: []string{mockDataString}},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionString{
					CollectionString: &rpc.CollectionString{
						String_: []string{mockDataString},
					},
				},
			},
		},
		{
			binding: "bytesSliceValue",
			output:  OutputTest{BytesSliceValue: [][]byte{[]byte(mockDataString)}},
			expected: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionBytes{
					CollectionBytes: &rpc.CollectionBytes{
						Bytes: [][]byte{[]byte(mockDataString)},
					},
				},
			},
		},
	}
	ctx := context.Background()
	for _, tt := range data {
		var function *OutputTest
		objectType, err := functionpkg.NewObjectType(
			reflect.TypeOf(function),
			functionpkg.HTTPTrigger,
			functionpkg.Bindings{},
			functionpkg.Bindings{
				functionpkg.Binding{
					Name: tt.binding,
				},
			},
		)
		assert.NoError(t, err)
		object := objectType.New()
		assert.NoError(t, object.Call(
			context.WithValue(ctx, expectedKey, tt.output),
			&loggerMock{},
			inputRPCHttpData,
			nil,
			functionpkg.BindingData{},
		))
		data, ok, err := object.GetOutput(tt.binding)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, tt.expected, data)
	}
}
