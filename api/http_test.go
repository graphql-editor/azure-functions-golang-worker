package api_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/graphql-editor/azure-functions-golang-worker/api"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	NullableTypes "github.com/graphql-editor/azure-functions-golang-worker/rpc/shared"
	"github.com/stretchr/testify/assert"
)

func strptr(s string) *string        { return &s }
func boolptr(b bool) *bool           { return &b }
func float64ptr(f float64) *float64  { return &f }
func timeptr(t time.Time) *time.Time { return &t }

func TestRequestUnmarshal(t *testing.T) {
	data := &rpc.TypedData{
		Data: &rpc.TypedData_Http{
			Http: &rpc.RpcHttp{
				Method: "GET",
				Url:    "https://mock/",
				Headers: map[string]string{
					"Content-Type":   "mock-content-type",
					"Content-Length": "200",
				},
				Params: map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				},
				Query: map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				},
				Body: &rpc.TypedData{
					Data: &rpc.TypedData_String_{
						String_: "data",
					},
				},
				RawBody: &rpc.TypedData{
					Data: &rpc.TypedData_Bytes{
						Bytes: []byte("data"),
					},
				},
			},
		},
	}
	var r api.Request
	r.Unmarshal(data)
}

func TestResponseUnmarshal(t *testing.T) {
	response := api.Response{
		Headers: http.Header{
			"Content-Type":   []string{"mock-content-type"},
			"Content-Length": []string{"200"},
		},
		Cookies: api.Cookies{
			api.Cookie{
				Name:     "cookie",
				Value:    "value",
				Domain:   strptr("mock.domain.com"),
				Path:     strptr("/mock/path"),
				Expires:  timeptr(time.Time{}),
				Secure:   boolptr(true),
				HTTPOnly: boolptr(true),
				MaxAge:   float64ptr(3600.0),
				SameSite: api.Strict,
			},
		},
		StatusCode: http.StatusOK,
		Body:       "data",
	}
	data, err := response.Marshal()
	assert.NoError(t, err)
	timestamp, _ := ptypes.TimestampProto(time.Time{})
	assert.Equal(t, &rpc.TypedData{
		Data: &rpc.TypedData_Http{
			Http: &rpc.RpcHttp{
				Headers: map[string]string{
					"Content-Type":   "mock-content-type",
					"Content-Length": "200",
				},
				Body: &rpc.TypedData{
					Data: &rpc.TypedData_String_{
						String_: "data",
					},
				},
				StatusCode: "200",
				Cookies: []*rpc.RpcHttpCookie{
					&rpc.RpcHttpCookie{
						Name:  "cookie",
						Value: "value",
						Domain: &NullableTypes.NullableString{
							String_: &NullableTypes.NullableString_Value{
								Value: "mock.domain.com",
							},
						},
						Path: &NullableTypes.NullableString{
							String_: &NullableTypes.NullableString_Value{
								Value: "/mock/path",
							},
						},
						Expires: &NullableTypes.NullableTimestamp{
							Timestamp: &NullableTypes.NullableTimestamp_Value{
								Value: timestamp,
							},
						},
						Secure: &NullableTypes.NullableBool{
							Bool: &NullableTypes.NullableBool_Value{
								Value: true,
							},
						},
						HttpOnly: &NullableTypes.NullableBool{
							Bool: &NullableTypes.NullableBool_Value{
								Value: true,
							},
						},
						MaxAge: &NullableTypes.NullableDouble{
							Double: &NullableTypes.NullableDouble_Value{
								Value: 3600.0,
							},
						},
						SameSite: rpc.RpcHttpCookie_Strict,
					},
				},
			},
		},
	}, data)
}

func TestResponseDefaultStatusOK(t *testing.T) {
	response := api.Response{
		Body: "data",
	}
	data, err := response.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, "200", data.Data.(*rpc.TypedData_Http).Http.StatusCode)
}
