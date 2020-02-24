package worker_test

import (
	"testing"

	"github.com/graphql-editor/azure-functions-golang-worker/worker"
	"github.com/stretchr/testify/assert"
)

func TestCapabilitiesRPCMarshal(t *testing.T) {
	assert.Equal(
		t,
		map[string]string{
			"RpcHttpTriggerMetadataRemoved": "true",
			"RpcHttpBodyOnly":               "true",
			"RawHttpBodyBytes":              "true",
		},
		worker.Capabilities{
			worker.RPCHttpTriggerMetadataRemoved: "true",
			worker.RPCHttpBodyOnly:               "true",
			worker.RawHTTPBodyBytes:              "true",
		}.ToRPC(),
	)
}
