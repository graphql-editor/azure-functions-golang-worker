package converters

import "github.com/graphql-editor/azure-functions-golang-worker/rpc"

// Unmarshaler interface is implemented by types that know how to handle their own conversion from rpc.TypedData
type Unmarshaler interface {
	// Unmarshal from rpc.TypedData
	Unmarshal(*rpc.TypedData) error
}
