package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/terag/kadok/internal/cache"
)

type Header struct {
	Key    string
	Values []string
}

type Request struct {
	Method        string
	Url           url.URL
	CacheKey      string
	CacheDuration time.Duration
	Headers       []Header
	Body          []byte
}

type Response struct {
	StatusCode int
	Headers    []Header
	CacheHit   bool
	Body       []byte
}

type Client interface {
	Execute(Request Request) (Response, error)
}

type HttpClient struct {
	Client http.Client
	Cache  cache.Cache
}

func (hc *HttpClient) Execute(request Request) (Response, error) {
	if request.CacheKey != "" && hc.Cache != nil {
		if v, ok, _ := hc.Cache.Get(request.CacheKey); ok {
			return Response{
				StatusCode: http.StatusOK,
				CacheHit:   true,
				Body:       v,
			}, nil
		}
	}

	httpRequest, err := http.NewRequest(request.Method, request.Url.String(), bytes.NewBuffer(request.Body))
	if err != nil {
		return Response{}, err
	}
	for _, header := range request.Headers {
		for _, hv := range header.Values {
			httpRequest.Header.Add(header.Key, hv)
		}
	}

	httpResponse, err := hc.Client.Do(httpRequest)
	if err != nil {
		return Response{}, err
	}

	responseBody, err := ioutil.ReadAll(httpResponse.Body)

	if httpResponse.StatusCode > 499 {
		return Response{}, HttpError{
			StatusCode: httpResponse.StatusCode,
			Url:        request.Url,
			Message:    string(responseBody),
		}
	}

	if httpResponse.StatusCode == 200 && request.CacheKey != "" && hc.Cache != nil {
		hc.Cache.Put(request.CacheKey, responseBody, time.Duration(5*time.Minute))
	}

	var headers []Header
	for k, v := range httpResponse.Header {
		headers = append(headers, Header{
			Key:    k,
			Values: v,
		})
	}

	return Response{
		StatusCode: httpResponse.StatusCode,
		Headers:    headers,
		CacheHit:   false,
		Body:       responseBody,
	}, nil
}
