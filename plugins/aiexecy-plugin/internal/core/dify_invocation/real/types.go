package real

import (
	"net/http"
	"net/url"
)

type RealBackwardsInvocation struct {
	aiexecInnerApiKey     string
	aiexecInnerApiBaseurl *url.URL
	client              *http.Client
	writeTimeout        int64
	readTimeout         int64
}

type BaseBackwardsInvocationResponse[T any] struct {
	Data  *T     `json:"data,omitempty"`
	Error string `json:"error"`
}
