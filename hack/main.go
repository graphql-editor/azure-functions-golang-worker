package main

import (
	"github.com/Dennor/gbtb"
)

func generateProto() error {
	gens := []struct {
		cmd  string
		args []string
	}{
		{
			cmd: "protoc",
			args: []string{
				"-I", "rpc", "--go_out=plugins=grpc,paths=source_relative:rpc/", "rpc/FunctionRpc.proto",
			},
		},
		{
			cmd: "protoc",
			args: []string{
				"-I", "rpc", "--go_out=plugins=grpc,paths=source_relative:rpc/", "rpc/identity/ClaimsIdentityRpc.proto",
			},
		},
		{
			cmd: "protoc",
			args: []string{
				"-I", "rpc", "--go_out=plugins=grpc,paths=source_relative:rpc/", "rpc/shared/NullableTypes.proto",
			},
		},
	}
	var err error
	for err == nil && len(gens) > 1 {
		err = gbtb.CommandJob(gens[0].cmd, gens[0].args...)()
		gens = gens[1:]
	}
	return err
}

func main() {
	gbtb.MustRun(
		&gbtb.Task{
			Name: "generate-proto",
			Job:  generateProto,
		},
	)
}
