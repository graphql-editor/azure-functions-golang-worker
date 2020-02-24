package worker

import (
	"fmt"

	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
	"github.com/pkg/errors"
)

// Logger that sends logs through rpc
type Logger struct {
	InvocationID string
	EventID      string
	Cat          rpc.RpcLog_RpcLogCategory
	Stream       Sender
}

// Trace logs current stack trace with message
func (l Logger) Trace(msg string) {
	st := errors.WithStack(errors.New("")).(stackTracer).StackTrace()
	msg = fmt.Sprintf("%s\n%+v", msg, st[1:])
	l.log(msg, rpc.RpcLog_Trace)
}

// Tracef logs current stack trace with message with formatted string
func (l Logger) Tracef(sfmt string, args ...interface{}) {
	l.Trace(fmt.Sprintf(sfmt, args...))
}

// Debug logs a debug level message
func (l Logger) Debug(msg string) {
	l.log(msg, rpc.RpcLog_Debug)
}

// Debugf logs current stack trace with message with formatted string
func (l Logger) Debugf(sfmt string, args ...interface{}) {
	l.Debug(fmt.Sprintf(sfmt, args...))
}

// Info logs an info level message
func (l Logger) Info(msg string) {
	l.log(msg, rpc.RpcLog_Information)
}

// Infof logs current stack trace with message with formatted string
func (l Logger) Infof(sfmt string, args ...interface{}) {
	l.Info(fmt.Sprintf(sfmt, args...))
}

// Warn logs an warning level message
func (l Logger) Warn(msg string) {
	l.log(msg, rpc.RpcLog_Warning)
}

// Warnf logs current stack trace with message with formatted string
func (l Logger) Warnf(sfmt string, args ...interface{}) {
	l.Warn(fmt.Sprintf(sfmt, args...))
}

// Error logs an error level message
func (l Logger) Error(msg string) {
	l.log(msg, rpc.RpcLog_Error)
}

// Errorf logs current stack trace with message with formatted string
func (l Logger) Errorf(sfmt string, args ...interface{}) {
	l.Error(fmt.Sprintf(sfmt, args...))
}

// Fatal logs an fatal level message
func (l Logger) Fatal(msg string) {
	l.log(msg, rpc.RpcLog_Critical)
}

// Fatalf logs current stack trace with message with formatted string
func (l Logger) Fatalf(sfmt string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(sfmt, args...))
}

func (l *Logger) log(msg string, level rpc.RpcLog_Level) {
	l.Stream.Send(&rpc.StreamingMessage{
		Content: &rpc.StreamingMessage_RpcLog{
			RpcLog: &rpc.RpcLog{
				InvocationId: l.InvocationID,
				EventId:      l.EventID,
				Message:      msg,
				Level:        level,
				Category:     rpc.RpcLog_RpcLogCategory_name[int32(l.Cat)],
			},
		},
	})
}
