# About

An attempt to implement golang worker for Azure functions host runner.

# Example

```golang
package main

import (
	"context"
	"fmt"

	"github.com/graphql-editor/azure-functions-golang-worker/api"
)

// HTTPTrigger is an example httpTrigger
type HTTPTrigger struct {
	Request  *api.Request `azfunc:"httpTrigger"`
	Response api.Response `azfunc:"res"`
}

// Run implements function behaviour
func (h *HttpTrigger) Run(ctx context.Context, logger api.Logger) {
	logger.Info(fmt.Sprintf("called with %v", h.Request))
	h.Response.Body = []byte("ok")
}

// Function exports function entry point
var Function HttpTrigger
```

# Getting started

Check out example

# Disclaimer

This is not an official Azure Project and as such is not supported by Microsoft.
