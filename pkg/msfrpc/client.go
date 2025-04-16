package msfrpc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

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
	}
}

func (c *MsfRpcClient) MsfAuth() (string, error) {
	payload := []interface{}{"auth.login", c.UserName, c.UserPassword}
	resp, err := c.MsfRequest(payload)
	if err != nil {
		return "", err
	}

	tokenBytes, ok := resp["token"].([]byte)
	token := string(tokenBytes)
	if !ok {
		return "", errors.New("missing 'token' in response – check credentials or msfrpcd settings")
	}
	c.Token = &token
	return token, nil
}

func (c *MsfRpcClient) AuthenticatedRequest(payload []any) (map[string]interface{}, error) {
	if c.Token == nil {
		return nil, errors.New("token is nil – you must authenticate first")
	}
	return c.MsfRequest(append(payload, *c.Token))
}

// Sends a generic RPC request and decodes the response
func (c *MsfRpcClient) MsfRequest(clientRequest []interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	encoder := msgpack.NewEncoder(&buf)
	if err := encoder.Encode(clientRequest); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "binary/message-pack")
	req.Header.Set("Host", "RPC Server")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var decoded map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

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
