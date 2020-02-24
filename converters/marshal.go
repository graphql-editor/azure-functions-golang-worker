package converters

import "github.com/graphql-editor/azure-functions-golang-worker/rpc"

// Marshaler converts data returned by user to rpc.TypedData
type Marshaler interface {
	// Marshal any type to rpc.TypedData
	Marshal() (*rpc.TypedData, error)
}
