package worker

import (
	"context"
	"fmt"
	"os"

	"github.com/graphql-editor/azure-functions-golang-worker/api"
	"github.com/graphql-editor/azure-functions-golang-worker/function"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/pkg/errors"
)

type channel struct {
	stream Sender
	loader Loader
	logger api.Logger
}

func (c *channel) StartStream(requestID string, msg *rpc.StartStream) {
	// not yet implemented
}

func (c *channel) InitRequest(requestID string, msg *rpc.WorkerInitRequest) {
	c.stream.Send(&rpc.StreamingMessage{
		RequestId: requestID,
		Content: &rpc.StreamingMessage_WorkerInitResponse{
			WorkerInitResponse: &rpc.WorkerInitResponse{
				Result: c.getStatus(nil),
				Capabilities: Capabilities{
					RPCHttpTriggerMetadataRemoved: "true",
					RPCHttpBodyOnly:               "true",
					RawHTTPBodyBytes:              "true",
				}.ToRPC(),
			},
		},
	})
}

func (c *channel) Heartbeat(requestID string, msg *rpc.WorkerHeartbeat) {
	// not yet implemented
}

func (c *channel) Terminate(requestID string, msg *rpc.WorkerTerminate) {
	// not yet implemented
}

func (c *channel) StatusRequest(requestID string, msg *rpc.WorkerStatusRequest) {
	// not yet implemented
}

func (c *channel) FileChangeEventRequest(requestID string, msg *rpc.FileChangeEventRequest) {
	// not yet implemented
}

func (c *channel) FunctionLoadRequest(requestID string, msg *rpc.FunctionLoadRequest) {
	functionID := msg.GetFunctionId()
	metadata := msg.GetMetadata()
	if functionID != "" && metadata != nil {
		err := c.loader.Load(functionID, metadata, c.logger)
		if err != nil {
			c.logger.Error(
				fmt.Sprintf(
					"Worker was unable to load function %s: %v",
					metadata.GetName(),
					err,
				),
			)
		}
		c.stream.Send(&rpc.StreamingMessage{
			RequestId: requestID,
			Content: &rpc.StreamingMessage_FunctionLoadResponse{
				FunctionLoadResponse: &rpc.FunctionLoadResponse{
					FunctionId: functionID,
					Result:     c.getStatus(err),
				},
			},
		})
	}
}

func (c *channel) InvocationRequest(requestID string, msg *rpc.InvocationRequest) {
	functionID := msg.GetFunctionId()
	info, err := c.loader.Info(functionID)
	outputData := make([]*rpc.ParameterBinding, 0, len(info.OutputBindings))
	var returnValue *rpc.TypedData
	inputData := make([]function.BindingData, 0, len(msg.InputData))
	var triggerData *rpc.TypedData
	var obj function.Object
	if err == nil {
		for _, binding := range msg.InputData {
			if info.TriggerBindingName == binding.GetName() {
				triggerData = binding.Data
			} else {
				inputData = append(inputData, function.BindingData{
					Name: binding.Name,
					Data: binding.Data,
				})
			}
		}
	}
	if err == nil {
		var objType function.ObjectType
		objType, err = c.loader.Func(functionID)
		if err == nil {
			obj = objType.New()
			err = obj.Call(
				context.Background(),
				Logger{
					InvocationID: msg.GetInvocationId(),
					EventID:      requestID,
					Stream:       c.stream,
					Cat:          rpc.RpcLog_User,
				},
				triggerData,
				msg.TriggerMetadata,
				inputData...,
			)
		}
	}
	if err == nil {
		var ok bool
		returnValue, ok, err = obj.ReturnValue()
		if !ok && err == nil {
			returnValue = nil
		}
		for name := range info.OutputBindings {
			var outputValue *rpc.TypedData
			outputValue, ok, err = obj.GetOutput(name)
			if err != nil {
				break
			}
			if ok {
				outputData = append(outputData, &rpc.ParameterBinding{
					Name: name,
					Data: outputValue,
				})
			}
		}
	}
	c.stream.Send(&rpc.StreamingMessage{
		RequestId: requestID,
		Content: &rpc.StreamingMessage_InvocationResponse{
			InvocationResponse: &rpc.InvocationResponse{
				InvocationId: msg.GetInvocationId(),
				OutputData:   outputData,
				ReturnValue:  returnValue,
				Result:       c.getStatus(err),
			},
		},
	})
}

func (c *channel) InvocationCancel(requestID string, msg *rpc.InvocationCancel) {
	// not yet implemented
}

func (c *channel) FunctionEnvironmentReloadRequest(requestID string, msg *rpc.FunctionEnvironmentReloadRequest) {
	c.logger.Info(fmt.Sprintf("Reloading environment variables. Found %d variables to reload", len(msg.EnvironmentVariables)))
	var err error
	for k, v := range msg.EnvironmentVariables {
		os.Setenv(k, v)
	}
	if msg.FunctionAppDirectory != "" {
		c.logger.Info(fmt.Sprintf("Changing current working directory to %s", msg.FunctionAppDirectory))
		err = os.Chdir(msg.FunctionAppDirectory)
	}
	c.stream.Send(&rpc.StreamingMessage{
		RequestId: requestID,
		Content: &rpc.StreamingMessage_FunctionEnvironmentReloadResponse{
			FunctionEnvironmentReloadResponse: &rpc.FunctionEnvironmentReloadResponse{
				Result: c.getStatus(err),
			},
		},
	})
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (c *channel) getStatus(err error) *rpc.StatusResult {
	status := &rpc.StatusResult{
		Status: rpc.StatusResult_Success,
	}
	if err != nil {
		status.Status = rpc.StatusResult_Failure
		status.Exception = &rpc.RpcException{
			Message: err.Error(),
		}
		if st, ok := err.(stackTracer); ok {
			status.Exception.StackTrace = fmt.Sprintf("%+v", st.StackTrace())
		}
	}
	return status
}

func (c *channel) SetEventStream(stream Sender) {
	c.stream = stream
	c.logger = Logger{
		Stream: stream,
		Cat:    rpc.RpcLog_System,
	}
}

func (c *channel) SetLoader(loader Loader) {
	c.loader = loader
}

// NewChannel create new default channel
func NewChannel() Channel {
	return &channel{}
}
