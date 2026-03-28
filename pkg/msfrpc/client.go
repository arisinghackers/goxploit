package msfrpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type MsfRpcClient struct {
	UserPassword string
	Token        *string
	SSL          string
	UserName     string
	IP           string
	Port         int
	WebServerURI string
	BaseURL      string
	HTTPClient   *http.Client
}

func NewMsfRpcClient(userPassword, ssl, userName, ip string, port int, webServerURI string) *MsfRpcClient {
	protocol := "http://"
	if ssl == "true" {
		protocol = "https://"
	}
	baseURL := fmt.Sprintf("%s%s:%d%s", protocol, ip, port, webServerURI)

	return &MsfRpcClient{
		UserPassword: userPassword,
		SSL:          protocol,
		UserName:     userName,
		IP:           ip,
		Port:         port,
		WebServerURI: webServerURI,
		BaseURL:      baseURL,
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *MsfRpcClient) MsfAuth() (string, error) {
	return c.MsfAuthContext(context.Background())
}

func (c *MsfRpcClient) MsfAuthContext(ctx context.Context) (string, error) {
	payload := []interface{}{"auth.login", c.UserName, c.UserPassword}
	resp, err := c.MsfRequestContext(ctx, payload)
	if err != nil {
		return "", err
	}

	rawToken, ok := resp["token"]
	if !ok {
		return "", errors.New("missing 'token' in response; check credentials or msfrpcd settings")
	}

	var token string
	switch v := rawToken.(type) {
	case []byte:
		token = string(v)
	case string:
		token = v
	default:
		return "", fmt.Errorf("invalid token type %T", rawToken)
	}

	c.Token = &token
	return token, nil
}

func (c *MsfRpcClient) AuthenticatedRequest(payload []any) (map[string]interface{}, error) {
	return c.AuthenticatedRequestContext(context.Background(), payload)
}

func (c *MsfRpcClient) AuthenticatedRequestContext(ctx context.Context, payload []any) (map[string]interface{}, error) {
	if c.Token == nil {
		return nil, errors.New("token is nil; you must authenticate first")
	}
	if len(payload) == 0 {
		return nil, errors.New("payload cannot be empty")
	}

	requestPayload := []interface{}{payload[0], *c.Token}
	for _, v := range payload[1:] {
		requestPayload = append(requestPayload, v)
	}

	return c.MsfRequestContext(ctx, requestPayload)
}

// Sends a generic RPC request and decodes the response
func (c *MsfRpcClient) MsfRequest(clientRequest []interface{}) (map[string]interface{}, error) {
	return c.MsfRequestContext(context.Background(), clientRequest)
}

// MsfRequestContext sends a generic RPC request using the provided context.
func (c *MsfRpcClient) MsfRequestContext(ctx context.Context, clientRequest []interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	encoder := msgpack.NewEncoder(&buf)
	if err := encoder.Encode(clientRequest); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "binary/message-pack")

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("rpc request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var decoded map[string]interface{}
	decoder := msgpack.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}

func (c *MsfRpcClient) SetToken(token string) {
	c.Token = &token
}

func (c *MsfRpcClient) GetToken() *string {
	return c.Token
}
