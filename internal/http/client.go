package http

import (
	"bytes"
	"fmt"
	"io"
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

type ResponseStream struct {
	StatusCode int
	Headers    []Header
	CacheHit   bool
	Stream     io.ReadCloser
}

type Client interface {
	Execute(request Request) (Response, error)
	OpenStream(request Request) (ResponseStream, error)
}

type HttpClient struct {
	client       http.Client
	streamClient http.Client
	Cache        cache.Cache
}

func NewHttpClient(cache cache.Cache, timeout time.Duration) *HttpClient {
	return &HttpClient{
		Cache: cache,
		client: http.Client{
			Timeout: timeout,
		},
		streamClient: http.Client{
			Timeout: 0,
		},
	}
}

func (hc *HttpClient) Execute(request Request) (Response, error) {
	if request.CacheKey != "" && hc.Cache != nil {
		if v, ok, _ := hc.Cache.Get(request.CacheKey); ok {
			fmt.Println("HTTP Execute - Cache Hit for Key: ", request.CacheKey)
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

	fmt.Println("HTTP Execute - Request: ", request.Method, " ", request.Url)
	httpResponse, err := hc.client.Do(httpRequest)
	if err != nil {
		return Response{}, err
	}
	defer httpResponse.Body.Close()

	responseBody, err := ioutil.ReadAll(httpResponse.Body)

	if httpResponse.StatusCode > 499 {
		fmt.Println("HTTP Execute - Response Error: ", httpResponse.Status, " - ", string(responseBody))
		return Response{}, HttpError{
			StatusCode: httpResponse.StatusCode,
			Url:        request.Url,
			Message:    string(responseBody),
		}
	}

	fmt.Println("HTTP Execute - Response: ", httpResponse.Status)
	if httpResponse.StatusCode == 200 && request.CacheKey != "" && hc.Cache != nil {
		fmt.Println("HTTP Execute - Caching: at ", request.CacheKey, " for ", request.CacheDuration)
		hc.Cache.Put(request.CacheKey, responseBody, request.CacheDuration)
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

func (hc *HttpClient) OpenStream(request Request) (ResponseStream, error) {
	httpRequest, err := http.NewRequest(request.Method, request.Url.String(), bytes.NewBuffer(request.Body))
	if err != nil {
		return ResponseStream{}, err
	}
	for _, header := range request.Headers {
		for _, hv := range header.Values {
			httpRequest.Header.Add(header.Key, hv)
		}
	}

	fmt.Println("HTTP OpenStream - Request: ", request.Method, " ", request.Url)
	httpResponse, err := hc.streamClient.Do(httpRequest)
	if err != nil {
		return ResponseStream{}, err
	}

	if httpResponse.StatusCode > 299 {
		responseBody, _ := ioutil.ReadAll(httpResponse.Body)
		fmt.Println("HTTP OpenStream - Response Error: ", httpResponse.Status, " - ", string(responseBody))
		return ResponseStream{}, HttpError{
			StatusCode: httpResponse.StatusCode,
			Url:        request.Url,
			Message:    string(responseBody),
		}
	}

	var headers []Header
	for k, v := range httpResponse.Header {
		headers = append(headers, Header{
			Key:    k,
			Values: v,
		})
	}

	if httpResponse.StatusCode != http.StatusOK {
		fmt.Println("HTTP OpenStream - Response Not Ok: ", httpResponse.Status)
		return ResponseStream{
			StatusCode: httpResponse.StatusCode,
			Headers:    headers,
			CacheHit:   false,
			Stream:     httpResponse.Body,
		}, nil
	}

	fmt.Println("HTTP OpenStream - Response: ", httpResponse.Status)

	return ResponseStream{
		StatusCode: httpResponse.StatusCode,
		Headers:    headers,
		CacheHit:   false,
		Stream:     httpResponse.Body,
	}, nil
}
