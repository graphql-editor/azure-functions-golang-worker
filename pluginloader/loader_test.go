package pluginloader_test

import (
	"testing"

	"github.com/graphql-editor/azure-functions-golang-worker/pluginloader"
	"github.com/graphql-editor/azure-functions-golang-worker/worker"
	"github.com/stretchr/testify/assert"
)

func TestPluginLoader(t *testing.T) {
	l := pluginloader.NewLoader()
	rt, err := l.GetFunctionType(worker.FunctionInfo{
		ScriptFile: "./testdata/function.go",
		EntryPoint: "Function",
	}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "main.HTTPTrigger", rt.String())
	assert.NoError(t, l.Close())
}
