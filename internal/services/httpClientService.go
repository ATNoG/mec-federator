package services

import (
	"context"
	"io"
	"net/http"
)

type HttpClientServiceInterface interface {
	DoRequest(ctx context.Context, method string, url string, body io.Reader, headers map[string]string, auth AuthStrategy) (*http.Response, error)
}

type HttpClientService struct {
	httpClient *http.Client
}

func NewHttpClientService(httpClient *http.Client) *HttpClientService {
	return &HttpClientService{
		httpClient: httpClient,
	}
}

// DoRequest performs an HTTP request with flexible auth and headers.
func (s *HttpClientService) DoRequest(
	ctx context.Context,
	method string,
	url string,
	body io.Reader,
	headers map[string]string,
	auth AuthStrategy,
) (*http.Response, error) {

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Apply headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Apply auth strategy if present
	if auth != nil {
		auth.Apply(req)
	}

	return s.httpClient.Do(req)
}

// AuthStrategy is an interface for applying authentication to an HTTP request.
type AuthStrategy interface {
	Apply(req *http.Request)
}

// BearerTokenAuth is an AuthStrategy that applies a Bearer token to an HTTP request.
type BearerTokenAuth struct {
	accessToken string
}

func NewBearerTokenAuth(accessToken string) AuthStrategy {
	return &BearerTokenAuth{accessToken: accessToken}
}

func (a *BearerTokenAuth) Apply(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+a.accessToken)
}
