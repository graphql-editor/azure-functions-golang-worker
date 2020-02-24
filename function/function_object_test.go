package function_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/graphql-editor/azure-functions-golang-worker/api"
	functionpkg "github.com/graphql-editor/azure-functions-golang-worker/function"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	inputRPCHttpData = &rpc.TypedData{
		Data: &rpc.TypedData_Http{
			Http: &rpc.RpcHttp{
				Method: "mockMethod",
				Url:    "http://mockUrl/path",
				Body: &rpc.TypedData{
					Data: &rpc.TypedData_String_{
						String_: "someBody",
					},
				},
				RawBody: &rpc.TypedData{
					Data: &rpc.TypedData_String_{
						String_: "someBody",
					},
				},
				Headers: map[string]string{
					"mockHeader": "mockValue",
				},
				Params: map[string]string{
					"mockParam": "mockValue",
				},
				Query: map[string]string{
					"mockQuery": "mockValue",
				},
			},
		},
	}
	originalBindingData = functionpkg.BindingData{
		Name: "original",
		Data: &rpc.TypedData{
			Data: &rpc.TypedData_Bytes{
				Bytes: []byte("mock-data"),
			},
		},
	}
	expectedMarshaledBody = `{
		"Method": "mockMethod",
		"URL": "http://mockUrl/path",
		"Headers": {"Mockheader": ["mockValue"]},
		"Query": {"mockQuery": ["mockValue"]},
		"Params": {"mockParam": ["mockValue"]},
		"Body": "someBody",
		"RawBody": "someBody"
	}`
)

type loggerMock struct {
	mock.Mock
}

func (m *loggerMock) Trace(msg string) {
	m.Called(msg)
}

func (m *loggerMock) Tracef(msg string, args ...interface{}) {
	args = append([]interface{}{msg}, args...)
	m.Called(args...)
}

func (m *loggerMock) Debug(msg string) {
	m.Called(msg)
}

func (m *loggerMock) Debugf(msg string, args ...interface{}) {
	args = append([]interface{}{msg}, args...)
	m.Called(args...)
}

func (m *loggerMock) Info(msg string) {
	m.Called(msg)
}

func (m *loggerMock) Infof(msg string, args ...interface{}) {
	args = append([]interface{}{msg}, args...)
	m.Called(args...)
}

func (m *loggerMock) Warn(msg string) {
	m.Called(msg)
}

func (m *loggerMock) Warnf(msg string, args ...interface{}) {
	args = append([]interface{}{msg}, args...)
	m.Called(args...)
}

func (m *loggerMock) Error(msg string) {
	m.Called(msg)
}

func (m *loggerMock) Errorf(msg string, args ...interface{}) {
	args = append([]interface{}{msg}, args...)
	m.Called(args...)
}

func (m *loggerMock) Fatal(msg string) {
	m.Called(msg)
}

func (m *loggerMock) Fatalf(msg string, args ...interface{}) {
	args = append([]interface{}{msg}, args...)
	m.Called(args...)
}

type FullBindingStruct struct {
	HTTPTrigger *api.Request
	Original    []byte
	Copy        []byte
	Res         api.Response
}

func (f *FullBindingStruct) Run(ctx context.Context, logger api.Logger) {
	f.Copy = append([]byte{}, f.Original...)
	f.Copy = append(f.Copy, []byte("-mock-copy")...)
	body, err := json.Marshal(f.HTTPTrigger)
	if err != nil {
		panic(err)
	}
	f.Res = api.Response{
		Headers: f.HTTPTrigger.Headers,
		Cookies: api.Cookies{
			api.Cookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: api.Lax,
			},
		},
		StatusCode: http.StatusOK,
		Body:       string(body),
	}
}

func TestFunctionStructWithoutReturn(t *testing.T) {
	var function *FullBindingStruct
	objectType, err := functionpkg.NewObjectType(
		reflect.TypeOf(function),
		functionpkg.HTTPTrigger,
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "original",
			},
		},
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "copy",
			},
			functionpkg.Binding{
				Name: "res",
				Type: "http",
			},
		},
	)
	assert.NoError(t, err)
	object := objectType.New()
	assert.NoError(t, object.Call(
		context.Background(),
		&loggerMock{},
		inputRPCHttpData,
		nil,
		originalBindingData,
	))
	returnValue, ok, err := object.ReturnValue()
	assert.False(t, ok)
	assert.NoError(t, err)
	assert.Nil(t, returnValue)
	resOutput, ok, err := object.GetOutput("res")
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.IsType(t, &rpc.TypedData_Http{}, resOutput.Data)
	rpcHTTP := resOutput.Data.(*rpc.TypedData_Http).Http
	assert.Equal(
		t,
		map[string]string{
			"Mockheader": "mockValue",
		},
		rpcHTTP.Headers,
	)
	assert.Equal(t, "200", rpcHTTP.StatusCode)
	assert.Equal(
		t,
		[]*rpc.RpcHttpCookie{
			&rpc.RpcHttpCookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: rpc.RpcHttpCookie_Lax,
			},
		},
		rpcHTTP.Cookies,
	)
	copyBinding, ok, err := object.GetOutput("copy")
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, copyBinding, &rpc.TypedData{
		Data: &rpc.TypedData_Bytes{
			Bytes: []byte("mock-data-mock-copy"),
		},
	})
	assert.JSONEq(t, expectedMarshaledBody, rpcHTTP.Body.Data.(*rpc.TypedData_String_).String_)
}

type FullBindingStructWithReturn FullBindingStruct

func (f *FullBindingStructWithReturn) Run(ctx context.Context, logger api.Logger) interface{} {
	f.Copy = append([]byte{}, f.Original...)
	f.Copy = append(f.Copy, []byte("-mock-copy")...)
	body, err := json.Marshal(f.HTTPTrigger)
	if err != nil {
		panic(err)
	}
	return api.Response{
		Headers: f.HTTPTrigger.Headers,
		Cookies: api.Cookies{
			api.Cookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: api.Lax,
			},
		},
		StatusCode: http.StatusOK,
		Body:       string(body),
	}
}

func TestFunctionStructWithReturn(t *testing.T) {
	var function *FullBindingStructWithReturn
	objectType, err := functionpkg.NewObjectType(
		reflect.TypeOf(function),
		functionpkg.HTTPTrigger,
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "original",
			},
		},
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "copy",
			},
			functionpkg.Binding{
				Name: "$return",
				Type: "http",
			},
		},
	)
	assert.NoError(t, err)
	object := objectType.New()
	assert.NoError(t, object.Call(
		context.Background(),
		&loggerMock{},
		inputRPCHttpData,
		nil,
		originalBindingData,
	))
	returnValue, ok, err := object.ReturnValue()
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.IsType(t, &rpc.TypedData_Http{}, returnValue.Data)
	rpcHTTP := returnValue.Data.(*rpc.TypedData_Http).Http
	assert.Equal(
		t,
		map[string]string{
			"Mockheader": "mockValue",
		},
		rpcHTTP.Headers,
	)
	assert.Equal(t, "200", rpcHTTP.StatusCode)
	assert.Equal(
		t,
		[]*rpc.RpcHttpCookie{
			&rpc.RpcHttpCookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: rpc.RpcHttpCookie_Lax,
			},
		},
		rpcHTTP.Cookies,
	)
	copyBinding, ok, err := object.GetOutput("copy")
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, copyBinding, &rpc.TypedData{
		Data: &rpc.TypedData_Bytes{
			Bytes: []byte("mock-data-mock-copy"),
		},
	})
	assert.JSONEq(t, expectedMarshaledBody, rpcHTTP.Body.Data.(*rpc.TypedData_String_).String_)
}

type FullBindingStructWithTags struct {
	Trigger  *api.Request `azfunc:"httpTrigger"`
	Blob     []byte       `azfunc:"original"`
	BlobCopy []byte       `azfunc:"copy"`
	Response api.Response `azfunc:"res"`
}

func (f *FullBindingStructWithTags) Run(ctx context.Context, logger api.Logger) {
	f.BlobCopy = append([]byte{}, f.Blob...)
	f.BlobCopy = append(f.BlobCopy, []byte("-mock-copy")...)
	body, err := json.Marshal(f.Trigger)
	if err != nil {
		panic(err)
	}
	f.Response = api.Response{
		Headers: f.Trigger.Headers,
		Cookies: api.Cookies{
			api.Cookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: api.Lax,
			},
		},
		StatusCode: http.StatusOK,
		Body:       string(body),
	}
}

func TestFunctionStructWithTagsWithoutReturn(t *testing.T) {
	var function *FullBindingStructWithTags
	objectType, err := functionpkg.NewObjectType(
		reflect.TypeOf(function),
		functionpkg.HTTPTrigger,
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "original",
			},
		},
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "copy",
			},
			functionpkg.Binding{
				Name: "res",
				Type: "http",
			},
		},
	)
	assert.NoError(t, err)
	object := objectType.New()
	assert.NoError(t, object.Call(
		context.Background(),
		&loggerMock{},
		inputRPCHttpData,
		nil,
		originalBindingData,
	))
	returnValue, ok, err := object.ReturnValue()
	assert.False(t, ok)
	assert.NoError(t, err)
	assert.Nil(t, returnValue)
	resOutput, ok, err := object.GetOutput("res")
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.IsType(t, &rpc.TypedData_Http{}, resOutput.Data)
	rpcHTTP := resOutput.Data.(*rpc.TypedData_Http).Http
	assert.Equal(
		t,
		map[string]string{
			"Mockheader": "mockValue",
		},
		rpcHTTP.Headers,
	)
	assert.Equal(t, "200", rpcHTTP.StatusCode)
	assert.Equal(
		t,
		[]*rpc.RpcHttpCookie{
			&rpc.RpcHttpCookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: rpc.RpcHttpCookie_Lax,
			},
		},
		rpcHTTP.Cookies,
	)
	copyBinding, ok, err := object.GetOutput("copy")
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, copyBinding, &rpc.TypedData{
		Data: &rpc.TypedData_Bytes{
			Bytes: []byte("mock-data-mock-copy"),
		},
	})
	assert.JSONEq(t, expectedMarshaledBody, rpcHTTP.Body.Data.(*rpc.TypedData_String_).String_)
}

type FullBindingStructWithTagsWithReturn FullBindingStructWithTags

func (f *FullBindingStructWithTagsWithReturn) Run(ctx context.Context, logger api.Logger) interface{} {
	f.BlobCopy = append([]byte{}, f.Blob...)
	f.BlobCopy = append(f.BlobCopy, []byte("-mock-copy")...)
	body, err := json.Marshal(f.Trigger)
	if err != nil {
		panic(err)
	}
	return api.Response{
		Headers: f.Trigger.Headers,
		Cookies: api.Cookies{
			api.Cookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: api.Lax,
			},
		},
		StatusCode: http.StatusOK,
		Body:       string(body),
	}
}

func TestFunctionStructWithTagsWithReturn(t *testing.T) {
	var function *FullBindingStructWithTagsWithReturn
	objectType, err := functionpkg.NewObjectType(
		reflect.TypeOf(function),
		functionpkg.HTTPTrigger,
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "original",
			},
		},
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "copy",
			},
			functionpkg.Binding{
				Name: "$return",
				Type: "http",
			},
		},
	)
	assert.NoError(t, err)
	object := objectType.New()
	assert.NoError(t, object.Call(
		context.Background(),
		&loggerMock{},
		inputRPCHttpData,
		nil,
		originalBindingData,
	))
	returnValue, ok, err := object.ReturnValue()
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.IsType(t, &rpc.TypedData_Http{}, returnValue.Data)
	rpcHTTP := returnValue.Data.(*rpc.TypedData_Http).Http
	assert.Equal(
		t,
		map[string]string{
			"Mockheader": "mockValue",
		},
		rpcHTTP.Headers,
	)
	assert.Equal(t, "200", rpcHTTP.StatusCode)
	assert.Equal(
		t,
		[]*rpc.RpcHttpCookie{
			&rpc.RpcHttpCookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: rpc.RpcHttpCookie_Lax,
			},
		},
		rpcHTTP.Cookies,
	)
	copyBinding, ok, err := object.GetOutput("copy")
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, copyBinding, &rpc.TypedData{
		Data: &rpc.TypedData_Bytes{
			Bytes: []byte("mock-data-mock-copy"),
		},
	})
	assert.JSONEq(t, expectedMarshaledBody, rpcHTTP.Body.Data.(*rpc.TypedData_String_).String_)
}

type MapFunction map[string]interface{}

func (f MapFunction) Run(ctx context.Context, logger api.Logger) {
	f["copy"] = append([]byte{}, f["original"].([]byte)...)
	f["copy"] = append(f["copy"].([]byte), []byte("-mock-copy")...)
	body, err := json.Marshal(f["httpTrigger"])
	if err != nil {
		panic(err)
	}
	f["res"] = api.Response{
		Headers: f["httpTrigger"].(map[string]interface{})["headers"].(http.Header),
		Cookies: api.Cookies{
			api.Cookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: api.Lax,
			},
		},
		StatusCode: http.StatusOK,
		Body:       string(body),
	}
}

func TestFunctionMapWithoutReturn(t *testing.T) {
	var function *MapFunction
	objectType, err := functionpkg.NewObjectType(
		reflect.TypeOf(function),
		functionpkg.HTTPTrigger,
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "original",
			},
		},
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "copy",
			},
			functionpkg.Binding{
				Name: "res",
				Type: "http",
			},
		},
	)
	assert.NoError(t, err)
	object := objectType.New()
	assert.NoError(t, object.Call(
		context.Background(),
		&loggerMock{},
		inputRPCHttpData,
		nil,
		originalBindingData,
	))
	returnValue, ok, err := object.ReturnValue()
	assert.False(t, ok)
	assert.NoError(t, err)
	assert.Nil(t, returnValue)
	resOutput, ok, err := object.GetOutput("res")
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.IsType(t, &rpc.TypedData_Http{}, resOutput.Data)
	rpcHTTP := resOutput.Data.(*rpc.TypedData_Http).Http
	assert.Equal(
		t,
		map[string]string{
			"Mockheader": "mockValue",
		},
		rpcHTTP.Headers,
	)
	assert.Equal(t, "200", rpcHTTP.StatusCode)
	assert.Equal(
		t,
		[]*rpc.RpcHttpCookie{
			&rpc.RpcHttpCookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: rpc.RpcHttpCookie_Lax,
			},
		},
		rpcHTTP.Cookies,
	)
	copyBinding, ok, err := object.GetOutput("copy")
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, copyBinding, &rpc.TypedData{
		Data: &rpc.TypedData_Bytes{
			Bytes: []byte("mock-data-mock-copy"),
		},
	})
	assert.JSONEq(t, `{
	"method": "mockMethod",
	"url": "http://mockUrl/path",
	"headers": {"Mockheader": ["mockValue"]},
	"query": {"mockQuery": ["mockValue"]},
	"params": {"mockParam": ["mockValue"]},
	"body": "someBody",
	"rawBody": "someBody"
}`, rpcHTTP.Body.Data.(*rpc.TypedData_String_).String_)
}

type MapFunctionWithReturn map[string]interface{}

func (f MapFunctionWithReturn) Run(ctx context.Context, logger api.Logger) interface{} {
	f["copy"] = append([]byte{}, f["original"].([]byte)...)
	f["copy"] = append(f["copy"].([]byte), []byte("-mock-copy")...)
	body, err := json.Marshal(f["httpTrigger"])
	if err != nil {
		panic(err)
	}
	return api.Response{
		Headers: f["httpTrigger"].(map[string]interface{})["headers"].(http.Header),
		Cookies: api.Cookies{
			api.Cookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: api.Lax,
			},
		},
		StatusCode: http.StatusOK,
		Body:       string(body),
	}
}

func TestFunctionMapWithReturn(t *testing.T) {
	var function *MapFunctionWithReturn
	objectType, err := functionpkg.NewObjectType(
		reflect.TypeOf(function),
		functionpkg.HTTPTrigger,
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "original",
			},
		},
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "copy",
			},
			functionpkg.Binding{
				Name: "res",
				Type: "http",
			},
		},
	)
	assert.NoError(t, err)
	object := objectType.New()
	assert.NoError(t, object.Call(
		context.Background(),
		&loggerMock{},
		inputRPCHttpData,
		nil,
		originalBindingData,
	))
	returnValue, ok, err := object.ReturnValue()
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.IsType(t, &rpc.TypedData_Http{}, returnValue.Data)
	rpcHTTP := returnValue.Data.(*rpc.TypedData_Http).Http
	assert.Equal(
		t,
		map[string]string{
			"Mockheader": "mockValue",
		},
		rpcHTTP.Headers,
	)
	assert.Equal(t, "200", rpcHTTP.StatusCode)
	assert.Equal(
		t,
		[]*rpc.RpcHttpCookie{
			&rpc.RpcHttpCookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: rpc.RpcHttpCookie_Lax,
			},
		},
		rpcHTTP.Cookies,
	)
	copyBinding, ok, err := object.GetOutput("copy")
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, copyBinding, &rpc.TypedData{
		Data: &rpc.TypedData_Bytes{
			Bytes: []byte("mock-data-mock-copy"),
		},
	})
	assert.JSONEq(t, `{
	"method": "mockMethod",
	"url": "http://mockUrl/path",
	"headers": {"Mockheader": ["mockValue"]},
	"query": {"mockQuery": ["mockValue"]},
	"params": {"mockParam": ["mockValue"]},
	"body": "someBody",
	"rawBody": "someBody"
}`, rpcHTTP.Body.Data.(*rpc.TypedData_String_).String_)
}

type Bad struct{}

type BadRun struct {
}

func (f BadRun) Run(ctx context.Context, logger api.Logger) {
}

type BadReturnRun struct {
}

func (f BadReturnRun) Run(ctx context.Context, logger api.Logger) interface{} {
	return nil
}

func TestFailBadFunctionTypes(t *testing.T) {
	var bad *Bad
	_, err := functionpkg.NewObjectType(
		reflect.TypeOf(bad),
		functionpkg.HTTPTrigger,
		functionpkg.Bindings{},
		functionpkg.Bindings{},
	)
	assert.Error(t, err)
	var badRun *BadRun
	_, err = functionpkg.NewObjectType(
		reflect.TypeOf(badRun),
		functionpkg.HTTPTrigger,
		functionpkg.Bindings{},
		functionpkg.Bindings{},
	)
	assert.Error(t, err)
	var badReturnRun *BadReturnRun
	_, err = functionpkg.NewObjectType(
		reflect.TypeOf(badReturnRun),
		functionpkg.HTTPTrigger,
		functionpkg.Bindings{},
		functionpkg.Bindings{},
	)
	assert.Error(t, err)
}

type OptionalBinding struct {
	HTTPTrigger *api.Request
	Res         api.Response
}

func (f *OptionalBinding) Run(ctx context.Context, logger api.Logger) {
	f.Res = api.Response{
		Headers: f.HTTPTrigger.Headers,
		Cookies: api.Cookies{
			api.Cookie{
				Name:     "cookie",
				Value:    "value",
				SameSite: api.Lax,
			},
		},
		StatusCode: http.StatusOK,
		Body:       "body",
	}
}

func TestOptionalBindings(t *testing.T) {
	var function *OptionalBinding
	objectType, err := functionpkg.NewObjectType(
		reflect.TypeOf(function),
		functionpkg.HTTPTrigger,
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "original",
			},
		},
		functionpkg.Bindings{
			functionpkg.Binding{
				Name: "copy",
			},
			functionpkg.Binding{
				Name: "res",
				Type: "http",
			},
		},
	)
	assert.NoError(t, err)
	object := objectType.New()
	assert.NoError(t, object.Call(
		context.Background(),
		&loggerMock{},
		inputRPCHttpData,
		nil,
		originalBindingData,
	))
	copyBinding, ok, err := object.GetOutput("copy")
	assert.False(t, ok)
	assert.NoError(t, err)
	assert.Nil(t, copyBinding)
}
