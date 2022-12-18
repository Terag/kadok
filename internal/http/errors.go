package http

import (
	"fmt"
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
	return "Http error, url: " + he.Url.String() + ", status: " + fmt.Sprint(he.StatusCode) + ", message: " + he.Message
}
