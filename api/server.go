// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get Well-Known
	// (GET /.well-known)
	GetWellKnown(ctx echo.Context) error
	// Get Well-Known Info
	// (GET /.well-known/info)
	GetWellKnownInfo(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetWellKnown converts echo context to params.
func (w *ServerInterfaceWrapper) GetWellKnown(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetWellKnown(ctx)
	return err
}

// GetWellKnownInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetWellKnownInfo(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetWellKnownInfo(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/.well-known", wrapper.GetWellKnown)
	router.GET(baseURL+"/.well-known/info", wrapper.GetWellKnownInfo)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
