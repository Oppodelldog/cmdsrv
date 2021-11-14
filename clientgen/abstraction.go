package clientgen

import (
	"fmt"
	"strings"

	"github.com/Oppodelldog/cmdsrv/rpcserver"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/usecase"
)

type Def struct {
	Methods []MethodDef
}

type MethodDef struct {
	Name    string
	Inputs  map[string]string
	Outputs map[string]string
}

func ApiDefinition(interactors []usecase.Interactor) Def {
	const (
		schemaPrefix = "#/components/schemas/"
		contentType  = "application/json"
		responseCode = "200"
	)

	var h = rpcserver.Handler(openapi3.Info{})
	for _, interactor := range interactors {
		h.Add(interactor)
	}

	var methods []MethodDef
	var getRefFields = func(schemaRef string) map[string]string {
		var (
			parameters = map[string]string{}
			refName    = strings.Replace(schemaRef, schemaPrefix, "", -1)
			refSchema  = h.OpenAPI.Reflector().Spec.Components.Schemas.MapOfSchemaOrRefValues[refName]
		)

		for propName, prop := range refSchema.Schema.Properties {
			parameters[propName] = string(*prop.Schema.Type)
		}

		return parameters
	}

	for methodName, methodSpec := range h.OpenAPI.Reflector().Spec.Paths.MapOfPathItemValues {
		var method = MethodDef{Name: methodName}
		for _, operation := range methodSpec.MapOfOperationValues {
			var (
				request  = operation.RequestBody.RequestBody
				response = operation.Responses.MapOfResponseOrRefValues[responseCode].Response
			)

			method.Inputs = getRefFields(request.Content[contentType].Schema.SchemaReference.Ref)
			method.Outputs = getRefFields(response.Content[contentType].Schema.SchemaReference.Ref)
		}

		methods = append(methods, method)
	}

	validateTypes(methods)
	validateOutput(methods)

	return Def{Methods: methods}
}

func validateTypes(methods []MethodDef) {
	for _, method := range methods {
		for _, t := range method.Inputs {
			validateType("input", t)
		}
		for _, t := range method.Outputs {
			validateType("output", t)
		}
	}
}

func validateType(s, t string) {
	if t != "string" &&
		t != "integer" &&
		t != "number" &&
		t != "array" {
		panic(fmt.Sprintf("%s type '%s' is not supported", s, t))
	}
}

func validateOutput(methods []MethodDef) {
	for _, method := range methods {
		if len(method.Outputs) != 1 {
			panic(fmt.Sprintf("output must have exactly one field, but it has %v", len(method.Outputs)))
		}

		for name := range method.Outputs {
			if name != "value" {
				panic(fmt.Sprintf("output field must be named 'value', it was '%s'", name))
			}
		}
	}
}
