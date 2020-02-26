package main

import (
	"reflect"

	httpTrigger "function/HttpTrigger"

	"github.com/graphql-editor/azure-functions-golang-worker/cmd/userworker"
)

func main() {
	userworker.Execute(map[string]reflect.Type{
		"HttpTrigger.Function": reflect.TypeOf(httpTrigger.Function),
	})
}
