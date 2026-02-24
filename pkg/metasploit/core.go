package metasploit

import (
	"context"

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
	return s.VersionContext(context.Background())
}

func (s *CoreService) VersionContext(ctx context.Context) (*CoreVersion, error) {
	resp, err := s.rpc.AuthenticatedRequestContext(ctx, []any{"core.version"})
	if err != nil {
		return nil, err
	}
	if err := checkRPCFailure(resp); err != nil {
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
