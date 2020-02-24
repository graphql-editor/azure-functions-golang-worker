package converters

import "github.com/graphql-editor/azure-functions-golang-worker/rpc"

var (
	typedDataDecoder TypedDataDecoder
	typedDataEncoder TypedDataEncoder
)

// Unmarshal typed data into anything
func Unmarshal(data *rpc.TypedData) (interface{}, error) {
	return typedDataDecoder.Decode(data)
}

// Marshal anything into typed data
func Marshal(data interface{}) (*rpc.TypedData, error) {
	return typedDataEncoder.Encode(data)
}
