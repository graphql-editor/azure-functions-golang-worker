// Package api available to user functions. User function must be a type
// that implements either Function or ReturnFunction interface.
// Function type for structs MUST be implemented on pointer reciever.
// User's function type MUST be exported using variable named Function (default) or otherwise defined EntryPoint in function.json (according to GoLang export rules, name MUST begin with capital letter). EntryPoint must be a valid GoLang identifier (https://golang.org/ref/spec#Identifiers).
// Inputs and outputs are read and written in similar fashion to a encoding/json package
//
// If a function object is a struct, the field name or a tag named `azfunc` must match trigger type for triggers and binding name for bindings. Field tag takes priority over field name. Field must be exported. If there's no tag, field name is compared using Unicode case-folding.
// If a function object is a map, keys in map match the type of a trigger for function triggers and a name of a binding for the rest of bindings.
//
// It is not an error if binding is missing from Function struct.
//
// For instance a struct object:
//  package main
//  type HTTPTrigger struct {
//  	// HttpTrigger represents function trigger. Function structure
//  	// can have at most one trigger defined.
//  	HttpTrigger *api.Request `azfunc:"httpTrigger"`
//  	// Additional input from, for instance, blob storage, named Original
//  	Original []byte `azfunc:"original"`
//  	// Response object using named output binding `res`
//  	Response api.Response `azfunc:"res"`
//  }
//  func (f *HTTPTrigger) Run(ctx context.Context, logger api.Logger) {
//  	f.Response.Body = "data"
//  }
//  var Function *HTTPTrigger
//
// For triggers that support $return value binding they can be implemented as so
//  package main
//  type HTTPTrigger struct {
//  	// HttpTrigger represents function trigger. Function structure
//  	// can have at most one trigger defined.
//  	HttpTrigger *api.Request `azfunc:"httpTrigger"`
//  }
//  func (f *HTTPTrigger) Run(ctx context.Context, logger api.Logger) interface{} {
//  	return api.Response{
//  		Body: "data",
//  	}
//  }
//  var Function *HTTPTrigger
//
// Worker also supports simple map type definitions as function objects
//  package main
//  type HTTPTrigger map[string]interface{}
//  func (f Function) Run(ctx context.Context, logger api.Logger) {
//  	f.["res"] = api.Response{
//  		Body: "data",
//  	}
//  }
//  var Function HTTPTrigger
//
// If scriptFile in function.json is empty, whole function package is built, similar to `go build .`, otherwise only file indicated by scriptFile is built and other go sources in function directory are ignored.
package api

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/graphql-editor/azure-functions-golang-worker/converters"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/pkg/errors"
)

// Logger interface
type Logger interface {
	Trace(string)
	Tracef(string, ...interface{})
	Debug(string)
	Debugf(string, ...interface{})
	Info(string)
	Infof(string, ...interface{})
	Warn(string)
	Warnf(string, ...interface{})
	Error(string)
	Errorf(string, ...interface{})
	Fatal(string)
	Fatalf(string, ...interface{})
}

// Function interface that must be implemented by user's function object. Function
// does not return a value.
type Function interface {
	Run(context.Context, Logger)
}

// ReturnFunction interface that must be implemented by user's function object. Function
// returns a value.
type ReturnFunction interface {
	Run(context.Context, Logger) interface{}
}

// Request represents httpTrigger in function definition.
type Request struct {
	Method  string
	URL     string
	Headers http.Header
	Query   url.Values
	Params  url.Values
	Body    interface{}
	RawBody interface{}
}

// Unmarshal implements unmarshaler for api.Request
func (r *Request) Unmarshal(data *rpc.TypedData) error {
	req, ok := data.Data.(*rpc.TypedData_Http)
	if !ok {
		return errors.Errorf("not a http request trigger")
	}
	body, rawBody, err := converters.DecodeHTTPBody(req.Http)
	if err == nil {
		*r = Request{
			Method:  req.Http.GetMethod(),
			URL:     req.Http.GetUrl(),
			Headers: converters.DecodeHeaders(req.Http.GetHeaders()),
			Query:   converters.DecodeValues(req.Http.GetQuery()),
			Params:  converters.DecodeValues(req.Http.GetParams()),
			Body:    body,
			RawBody: rawBody,
		}
	}
	return nil
}

// CookiePolicy for cross-site requests
type CookiePolicy string

const (
	// Strict policy
	Strict CookiePolicy = "Strict"
	// Lax policy
	Lax CookiePolicy = "Lax"
)

// Cookie used with http response Set-Cookie
type Cookie struct {
	Name     string
	Value    string
	Domain   *string
	Path     *string
	Expires  *time.Time
	Secure   *bool
	HTTPOnly *bool
	SameSite CookiePolicy
	MaxAge   *float64
}

// Cookies list
type Cookies []Cookie

// Response represents response from function when using HTTP output binding
type Response struct {
	Headers    http.Header
	Cookies    Cookies
	StatusCode int
	Body       interface{}
}

// Marshal implements Marshaler for converters
func (r Response) Marshal() (*rpc.TypedData, error) {
	return encodeResponseObject(&r)
}

type contextKey string

// TriggerMetadataKey is a context key for trigger data
var TriggerMetadataKey contextKey = contextKey("triggerMetadataKey")

// GetTriggerMetadata returns trigger metadata associated with trigger
func GetTriggerMetadata(ctx context.Context) map[string]interface{} {
	v, ok := ctx.Value(TriggerMetadataKey).(map[string]interface{})
	if !ok {
		return nil
	}
	return v
}
