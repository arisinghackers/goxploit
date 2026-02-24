package metasploit

import "github.com/arisinghackers/goxploit/pkg/msfrpc"

type ModuleService struct {
	rpc *msfrpc.MsfRpcClient
}

type ExecuteModuleRequest struct {
	ModuleType string
	ModuleName string
	Options    map[string]interface{}
}

type ExecuteModuleResponse struct {
	JobID int64
	UUID  string
}

func (s *ModuleService) Execute(req ExecuteModuleRequest) (*ExecuteModuleResponse, error) {
	options := req.Options
	if options == nil {
		options = map[string]interface{}{}
	}

	resp, err := s.rpc.AuthenticatedRequest([]any{
		"module.execute",
		req.ModuleType,
		req.ModuleName,
		options,
	})
	if err != nil {
		return nil, err
	}
	if err := checkRPCFailure(resp); err != nil {
		return nil, err
	}

	jobID, err := readInt(resp, "job_id", true)
	if err != nil {
		return nil, err
	}
	uuid, err := readString(resp, "uuid", true)
	if err != nil {
		return nil, err
	}

	return &ExecuteModuleResponse{
		JobID: jobID,
		UUID:  uuid,
	}, nil
}
