package metasploit

import (
	"context"

	"github.com/arisinghackers/goxploit/pkg/msfrpc"
)

type AuthService struct {
	rpc *msfrpc.MsfRpcClient
}

type LoginResponse struct {
	Token string
}

func (s *AuthService) Login(userName, userPassword string) (*LoginResponse, error) {
	return s.LoginContext(context.Background(), userName, userPassword)
}

func (s *AuthService) LoginContext(ctx context.Context, userName, userPassword string) (*LoginResponse, error) {
	resp, err := s.rpc.MsfRequestContext(ctx, []interface{}{"auth.login", userName, userPassword})
	if err != nil {
		return nil, err
	}
	if err := checkRPCFailure(resp); err != nil {
		return nil, err
	}

	token, err := readString(resp, "token", true)
	if err != nil {
		return nil, err
	}
	s.rpc.SetToken(token)

	return &LoginResponse{Token: token}, nil
}
