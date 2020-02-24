package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/graphql-editor/azure-functions-golang-worker/pluginloader"
	"github.com/graphql-editor/azure-functions-golang-worker/worker"
)

var (
	host                 = flag.String("host", "", "listend address, required")
	port                 = flag.String("port", "", "listend port, required")
	workerID             = flag.String("workerId", "", "worker id, required")
	requestID            = flag.String("requestId", "", "request id, required")
	grpcMaxMessageLength = flag.Int("grpcMaxMessageLength", 0, "grpc message lenght limit, required")
)

func main() {
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
			TypeLoader:      loader,
			LoadedFunctions: make(map[string]worker.Function),
		},
	}

	if err := w.Listen(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
