package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mankings/mec-federator/internal/models"
)

type FederationHttpClientManager struct {
	mu      sync.RWMutex
	clients map[string]*HttpClient // federation-context-id → HttpClient
}

func NewFederationHttpClientManager() *FederationHttpClientManager {
	return &FederationHttpClientManager{
		clients: make(map[string]*HttpClient),
	}
}

func (m *FederationHttpClientManager) Register(federationContextId string, client *HttpClient) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[federationContextId] = client
}

func (m *FederationHttpClientManager) Get(federationContextId string) (*HttpClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, ok := m.clients[federationContextId]
	if !ok {
		return nil, fmt.Errorf("client not found for federation context id: %s", federationContextId)
	}
	return client, nil
}

type HttpClientConfig struct {
	BaseUrl       string
	TokenEndpoint string
	ClientId      string
	ClientSecret  string
}

type HttpClient struct {
	config     *HttpClientConfig
	httpClient *http.Client

	mu          sync.Mutex
	accessToken models.AccessToken
}

func NewHttpClient(config HttpClientConfig) *HttpClient {
	return &HttpClient{
		config:     &config,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// get access token from cache or fetch new one
func (f *HttpClient) GetAccessToken(ctx context.Context) (models.AccessToken, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.accessToken.Token != "" && time.Now().Before(f.accessToken.ExpiresAt) {
		return f.accessToken, nil
	}

	token, err := f.fetchAccessToken(ctx)
	if err != nil {
		return models.AccessToken{}, fmt.Errorf("failed to fetch access token: %w", err)
	}

	f.accessToken = token

	return f.accessToken, nil
}

// fetch access token from auth server
func (f *HttpClient) fetchAccessToken(ctx context.Context) (models.AccessToken, error) {
	payload := map[string]string{
		"clientId":     f.config.ClientId,
		"clientSecret": f.config.ClientSecret,
	}

	resp, err := f.PostJSON(ctx, f.config.TokenEndpoint, payload)
	if err != nil {
		return models.AccessToken{}, fmt.Errorf("failed to post json: %w", err)
	}
	defer resp.Body.Close()

	var token models.AccessToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return models.AccessToken{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return token, nil
}

// simple json post wrapper method
func (f *HttpClient) PostJSON(ctx context.Context, url string, payload interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return f.httpClient.Do(req)
}

// json post with auth wrapper method
func (f *HttpClient) PostJSONWithAuth(ctx context.Context, url string, payload interface{}) (*http.Response, error) {
	token, err := f.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Token)

	return f.httpClient.Do(req)
}

// json post with auth and unmarshal wrapper method
func (f *HttpClient) PostJSONWithAuthAndUnmarshal(ctx context.Context, url string, payload interface{}, response interface{}) error {
	resp, err := f.PostJSONWithAuth(ctx, url, payload)
	if err != nil {
		return fmt.Errorf("failed to post json with auth: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
