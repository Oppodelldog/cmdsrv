package rpcserver

import (
	"github.com/swaggest/jsonrpc"
	"github.com/swaggest/openapi-go/openapi3"
)

func Handler(info openapi3.Info) *jsonrpc.Handler {
	var (
		apiSchema = jsonrpc.OpenAPI{}
		validator = jsonrpc.JSONSchemaValidator{}
		handler   = jsonrpc.Handler{}
	)

	apiSchema.Reflector().SpecEns().Info = info

	handler.OpenAPI = &apiSchema
	handler.Validator = &validator

	return &handler
}
