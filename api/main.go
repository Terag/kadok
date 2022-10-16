package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Properties struct {
	Enable             bool
	ListeningInterface string
	HttpPort           int
	Hostname           string
	BasePath           string
}

// Run Kadok HTTP API server.
// In case of error when launching the server, `onError` is called
// Once the server is launched, call onReady and wait on its return to close the server.
func Run(onReady func(), onError func(error), p Properties) {

	e := echo.New()
	// Cleanly close down the API
	defer func() {
		if err := e.Close(); err != nil {
			fmt.Println("Kadok API error: ", err)
		}
	}()
	kas := KadokApiServer{
		Hostname: p.Hostname,
		BasePath: p.BasePath,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = httpErrorHandler

	RegisterHandlersWithBaseURL(e, kas, p.BasePath)

	s := http.Server{
		Addr: p.Hostname + ":" + fmt.Sprint(p.HttpPort),
	}

	go func() {
		if err := e.StartServer(&s); err != nil {
			onError(errors.New("Kadok API error: " + err.Error()))
			return
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	onReady()
}

// KadokApiServer
type KadokApiServer struct {
	Hostname string
	BasePath string
	HttpPort int
}

// nullable is used by the API domain to get a pointer on a variable. This is used to manage nullable fields in API structures.
func nullable[T any](v T) *T {
	return &v
}

// httpStatusCodeToMessage Convert a http status code to human readable message
func httpStatusCodeToMessage(code int) string {
	// StringToRoleType Convert a string to a RoleType
	if message, ok := map[int]string{
		http.StatusOK:                           "Ok",                            // 200
		http.StatusCreated:                      "Created",                       // 201
		http.StatusAccepted:                     "Accepted",                      // 202
		http.StatusNoContent:                    "No Content",                    // 204
		http.StatusPartialContent:               "Partial Content",               // 206
		http.StatusMultipleChoices:              "Multiple Choices",              // 300
		http.StatusMovedPermanently:             "Moved Permanently",             // 301
		http.StatusFound:                        "Found",                         // 302
		http.StatusSeeOther:                     "See other",                     // 303
		http.StatusNotModified:                  "Not Modified",                  // 304
		http.StatusTemporaryRedirect:            "Temporary Redirect",            // 307
		http.StatusPermanentRedirect:            "Permanent Redirect",            // 308
		http.StatusBadRequest:                   "Bad Request",                   // 400
		http.StatusUnauthorized:                 "Unauthorized",                  // 401
		http.StatusForbidden:                    "Forbidden",                     // 403
		http.StatusNotFound:                     "Not Found",                     // 404
		http.StatusMethodNotAllowed:             "Method Not Allowed",            // 405
		http.StatusNotAcceptable:                "Not Acceptable",                // 406
		http.StatusProxyAuthRequired:            "Proxy Authentication Required", // 407
		http.StatusRequestTimeout:               "Request Timeout",               // 408
		http.StatusConflict:                     "Conflict",                      // 409
		http.StatusPreconditionFailed:           "Precondition Failed",           // 412
		http.StatusUnsupportedMediaType:         "Unsupported Media Type",        // 415
		http.StatusRequestedRangeNotSatisfiable: "Range Not Satisfiable",         // 416
		http.StatusInternalServerError:          "Internal Server Error",         // 500
		http.StatusNotImplemented:               "Not Implemented",               // 501
	}[code]; ok {
		return message
	}
	return "Internal Server Error"
}
