package metasploit

import (
	"fmt"

	"github.com/arisinghackers/goxploit/pkg/msfrpc"
)

type CoreService struct {
	rpc *msfrpc.MsfRpcClient
}

// CoreVersion is a typed view of the core.version response.
type CoreVersion struct {
	Version string
	Ruby    string
	API     string
}

func (s *CoreService) Version() (*CoreVersion, error) {
	resp, err := s.rpc.AuthenticatedRequest([]any{"core.version"})
	if err != nil {
		return nil, err
	}

	version, err := readString(resp, "version", true)
	if err != nil {
		return nil, err
	}
	ruby, err := readString(resp, "ruby", false)
	if err != nil {
		return nil, err
	}
	api, err := readString(resp, "api", false)
	if err != nil {
		return nil, err
	}

	return &CoreVersion{
		Version: version,
		Ruby:    ruby,
		API:     api,
	}, nil
}

func readString(resp map[string]interface{}, key string, required bool) (string, error) {
	v, ok := resp[key]
	if !ok {
		if required {
			return "", fmt.Errorf("missing %q in response", key)
		}
		return "", nil
	}

	switch typed := v.(type) {
	case string:
		return typed, nil
	case []byte:
		return string(typed), nil
	default:
		return "", fmt.Errorf("field %q has invalid type %T", key, v)
	}
}
