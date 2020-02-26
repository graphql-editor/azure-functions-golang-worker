// Package userworker adds support for manual definition of functions without
// using go plugins. This is mostly a drop in replacement for those who need windows support.
// It's main limitation is that there must be one common main.go entrypoint with all
// functions defined and provided as argument to Execute. Function name must be defined as
// {FunctionName}.{FunctionEntryPoint}.
//
// It requires a different language worker.config.json then plugin version.
//
// If your main file is named worker.go and is located in project root then worker.config.json will look something like this.
//  {
//      "description":{
//          "language":"golang",
//          "extensions":[".go"],
//          "defaultExecutablePath":"go",
//          "arguments": ["run", "./worker.go"]
//      }
//  }
//
// Example below has been created with `go mod init function` and function is defined
// in HttpTrigger package:
//
//  package main
//  import (
//    "context"
//    "fmt"
//
//    httpTrigger "function/HttpTrigger"
//
//    "github.com/graphql-editor/azure-functions-golang-worker/api"
//    "github.com/graphql-editor/azure-functions-golang-worker/cmd/userworker"
//  )
//
//  func main() {
//    userworker.Execute(map[string]reflect.Type{
//      "HttpTrigger.Function": reflect.TypeOf(Function),
//    })
//  }
package userworker

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/graphql-editor/azure-functions-golang-worker/api"
	"github.com/graphql-editor/azure-functions-golang-worker/pluginloader"
	"github.com/graphql-editor/azure-functions-golang-worker/worker"
	"github.com/pkg/errors"
)

var (
	host                 = flag.String("host", "", "listend address, required")
	port                 = flag.String("port", "", "listend port, required")
	workerID             = flag.String("workerId", "", "worker id, required")
	requestID            = flag.String("requestId", "", "request id, required")
	grpcMaxMessageLength = flag.Int("grpcMaxMessageLength", 0, "grpc message lenght limit, required")
)

type localLoader map[string]reflect.Type

func (l localLoader) GetFunctionType(fi worker.FunctionInfo, logger api.Logger) (reflect.Type, error) {
	t, ok := l[fi.Name+"."+fi.EntryPoint]
	if !ok {
		return nil, errors.Errorf("could not load function from file %s named %s", fi.ScriptFile, fi.EntryPoint)
	}
	return t, nil
}

// Execute worker with functions defined manually by user.
func Execute(functions map[string]reflect.Type) {
	flag.Parse()
	if *host == "" || *port == "" || *workerID == "" || *requestID == "" || *grpcMaxMessageLength == 0 {
		flag.Usage()
		os.Exit(2)
	}
	loader := pluginloader.NewLoader()
	defer loader.Close()

	w := worker.Worker{
		WorkerID:  *workerID,
		RequestID: *requestID,
		Stream: worker.NewEventStream(
			worker.HostEventStreamOption(*host),
			worker.PortEventStreamOption(*port),
		),
		Loader: worker.Loader{
			TypeLoader:      localLoader(functions),
			LoadedFunctions: make(map[string]worker.Function),
		},
	}

	if err := w.Listen(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
