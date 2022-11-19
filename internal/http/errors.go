package http

import (
	"net/http"
	"net/url"
)

var code = http.StatusOK

type HttpError struct {
	StatusCode int
	Url        url.URL
	Message    string
}

func (he HttpError) Error() string {
	return "Http error, url: " + he.Url.String() + ", status: " + string(he.StatusCode) + ", message: " + he.Message
}
