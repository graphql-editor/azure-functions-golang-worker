package converters_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/graphql-editor/azure-functions-golang-worker/converters"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/stretchr/testify/assert"
)

func skipCompareAssertion(assert.TestingT, interface{}, interface{}, ...interface{}) bool {
	return true
}

type stringLike string
type int64Like int64
type uintLike uint32
type floatLike float64
type sliceOfStrings []string
type sliceOfInts []int64
type sliceOfUints []uint32
type sliceOfFloats []float64
type bytes []byte
type sliceOfBytes [][]byte

type stringImplementingMarshaler string

func (m stringImplementingMarshaler) Marshal() (*rpc.TypedData, error) {
	return &rpc.TypedData{
		Data: &rpc.TypedData_Bytes{
			Bytes: []byte("mocked marshaler"),
		},
	}, nil
}

func TestTypedData(t *testing.T) {
	data := []struct {
		data          interface{}
		rpcData       *rpc.TypedData
		decodeCompare assert.ComparisonAssertionFunc
		decodeErr     assert.ErrorAssertionFunc
		encodeCompare assert.ComparisonAssertionFunc
		encodeErr     assert.ErrorAssertionFunc
	}{
		{
			data: "data",
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: "data",
				},
			},
		},
		{
			data: map[string]interface{}{
				"arbitrary": map[string]interface{}{
					"json": "object",
				},
				"key1": 1,
				"key2": []string{"array", "of", "values"},
			},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: `{
	"arbitrary": {
		"json": "object"
	},
	"key1": 1,
	"key2": ["array", "of", "values"]
}`,
				},
			},
			decodeCompare: func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
				expectedJSON, err := json.Marshal(expected)
				actualJSON, actualErr := json.Marshal(actual)
				return assert.NoError(t, err) &&
					assert.NoError(t, actualErr) &&
					assert.JSONEq(t, string(expectedJSON), string(actualJSON))
			},
			encodeCompare: func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
				return assert.IsType(t, &rpc.TypedData{}, expected) &&
					assert.IsType(t, expected, actual) &&
					assert.IsType(
						t,
						&rpc.TypedData_Json{},
						expected.(*rpc.TypedData).Data,
					) &&
					assert.IsType(
						t,
						expected.(*rpc.TypedData).Data,
						actual.(*rpc.TypedData).Data,
					) &&
					assert.JSONEq(
						t,
						expected.(*rpc.TypedData).Data.(*rpc.TypedData_Json).Json,
						actual.(*rpc.TypedData).Data.(*rpc.TypedData_Json).Json,
					)
			},
		},
		{
			data: []byte("raw bytes"),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Bytes{
					Bytes: []byte("raw bytes"),
				},
			},
		},
		{
			data: map[string]interface{}{
				"body": "data",
				"headers": http.Header{
					"Header": []string{"value"},
				},
				"method": "GET",
				"params": url.Values{
					"key": []string{"value"},
				},
				"query": url.Values{
					"key": []string{"value"},
				},
				"rawBody": []byte("data"),
				"url":     "www.example.com",
			},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Http{
					Http: &rpc.RpcHttp{
						Body: &rpc.TypedData{
							Data: &rpc.TypedData_String_{
								String_: "data",
							},
						},
						Headers: map[string]string{
							"header": "value",
						},
						Method: "GET",
						Params: map[string]string{
							"key": "value",
						},
						Query: map[string]string{
							"key": "value",
						},
						RawBody: &rpc.TypedData{
							Data: &rpc.TypedData_Bytes{
								Bytes: []byte("data"),
							},
						},
						Url: "www.example.com",
					},
				},
			},
			encodeCompare: skipCompareAssertion,
		},
		{
			data: int64(1),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: int64(1),
				},
			},
		},
		{
			data: float64(1),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Double{
					Double: float64(1),
				},
			},
		},
		{
			data: [][]byte{[]byte("first"), []byte("second")},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionBytes{
					CollectionBytes: &rpc.CollectionBytes{
						Bytes: [][]byte{[]byte("first"), []byte("second")},
					},
				},
			},
		},
		{
			data: []string{"first", "second"},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionString{
					CollectionString: &rpc.CollectionString{
						String_: []string{"first", "second"},
					},
				},
			},
		},
		{
			data: []int64{int64(1), int64(2)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
		},
		{
			data: []float64{1.0, 2.0},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionDouble{
					CollectionDouble: &rpc.CollectionDouble{
						Double: []float64{1.0, 2.0},
					},
				},
			},
		},
		// decoder reports error on error, for example bad json
		{
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Json{
					Json: `{key: value}`,
				},
			},
			decodeErr:     assert.Error,
			encodeCompare: skipCompareAssertion,
		},
		// decoder reports error on unsuppported stream
		{
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Stream{
					Stream: []byte{},
				},
			},
			decodeErr:     assert.Error,
			encodeCompare: skipCompareAssertion,
		},
		// Test encoder conversions from types that can not apear
		// from typeddata, but can apear from go code.
		{
			data: int(1),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: int64(1),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: int8(1),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: int64(1),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: int16(1),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: int64(1),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: int32(1),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: int64(1),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: uint8(1),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: int64(1),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: uint16(1),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: int64(1),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: uint32(1),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: int64(1),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: float32(1.0),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Double{
					Double: float64(1.0),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: []int{1, 2},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: []int8{int8(1), int8(2)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: []int16{int16(1), int16(2)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: []int32{int32(1), int32(2)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: []uint16{uint16(1), uint16(2)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: []uint32{uint32(1), uint32(2)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: []float32{float32(1.0), float32(2.0)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionDouble{
					CollectionDouble: &rpc.CollectionDouble{
						Double: []float64{float64(1.0), float64(2.0)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		// Encoder's Laws of Reflection shall be tested hereafter
		{
			data: strptr("data"),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: "data",
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: stringLike("data"),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_String_{
					String_: "data",
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: int64Like(1),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: int64(1),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: uintLike(1),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Int{
					Int: int64(1),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: floatLike(1.0),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Double{
					Double: float64(1.0),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: [3]byte{'a', 'b', 'c'},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Bytes{
					Bytes: []byte("abc"),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: bytes("abc"),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Bytes{
					Bytes: []byte("abc"),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: [2]string{"abc", "def"},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionString{
					CollectionString: &rpc.CollectionString{
						String_: []string{"abc", "def"},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: sliceOfStrings{"abc", "def"},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionString{
					CollectionString: &rpc.CollectionString{
						String_: []string{"abc", "def"},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: []stringLike{stringLike("abc"), stringLike("def")},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionString{
					CollectionString: &rpc.CollectionString{
						String_: []string{"abc", "def"},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: [2]int64{int64(1), int64(2)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: sliceOfInts{int64(1), int64(2)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: []int64Like{int64Like(1), int64Like(2)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: [2]uint32{uint32(1), uint32(2)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: sliceOfUints{uint32(1), uint32(2)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionSint64{
					CollectionSint64: &rpc.CollectionSInt64{
						Sint64: []int64{int64(1), int64(2)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: [2]float64{float64(1.0), float64(2.0)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionDouble{
					CollectionDouble: &rpc.CollectionDouble{
						Double: []float64{float64(1.0), float64(2.0)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: sliceOfFloats{float64(1.0), float64(2.0)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionDouble{
					CollectionDouble: &rpc.CollectionDouble{
						Double: []float64{float64(1.0), float64(2.0)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: []floatLike{floatLike(1.0), floatLike(2.0)},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionDouble{
					CollectionDouble: &rpc.CollectionDouble{
						Double: []float64{float64(1.0), float64(2.0)},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: [2][]byte{[]byte("abc"), []byte("def")},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionBytes{
					CollectionBytes: &rpc.CollectionBytes{
						Bytes: [][]byte{[]byte("abc"), []byte("def")},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: sliceOfBytes{[]byte("abc"), []byte("def")},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionBytes{
					CollectionBytes: &rpc.CollectionBytes{
						Bytes: [][]byte{[]byte("abc"), []byte("def")},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: [2][2]byte{[2]byte{'a', 'b'}, [2]byte{'c', 'd'}},
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_CollectionBytes{
					CollectionBytes: &rpc.CollectionBytes{
						Bytes: [][]byte{[]byte("ab"), []byte("cd")},
					},
				},
			},
			decodeCompare: skipCompareAssertion,
		},
		{
			data: stringImplementingMarshaler("abc"),
			rpcData: &rpc.TypedData{
				Data: &rpc.TypedData_Bytes{
					Bytes: []byte("mocked marshaler"),
				},
			},
			decodeCompare: skipCompareAssertion,
		},
	}
	for i := range data {
		if data[i].decodeErr == nil {
			data[i].decodeErr = assert.NoError
		}
		if data[i].encodeErr == nil {
			data[i].encodeErr = assert.NoError
		}
		if data[i].decodeCompare == nil {
			data[i].decodeCompare = assert.Equal
		}
		if data[i].encodeCompare == nil {
			data[i].encodeCompare = assert.Equal
		}
	}

	t.Run("DecodeTypedData", func(t *testing.T) {
		dec := &converters.TypedDataDecoder{}
		for _, tc := range data {
			data, err := dec.Decode(tc.rpcData)
			tc.decodeErr(t, err)
			tc.decodeCompare(t, tc.data, data)
		}
	})
	t.Run("EncodeTypedData", func(t *testing.T) {
		enc := &converters.TypedDataEncoder{}
		for _, tc := range data {
			data, err := enc.Encode(tc.data)
			tc.encodeErr(t, err)
			tc.encodeCompare(t, tc.rpcData, data)
		}
	})
}
