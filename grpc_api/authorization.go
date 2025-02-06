package grpc_api

import (
	"context"
	"fmt"
	"github.com/komron-dev/bank/token"
	"github.com/komron-dev/bank/util"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

var AccessibleRoles = []string{util.DepositorRole, util.BankerRole}

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)

	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)

	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported auth type")
	}
	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	if !hasPermission(payload.Role) {
		return nil, fmt.Errorf("forbidden role: permission denied")
	}
	return payload, nil
}

func hasPermission(role string) bool {
	for _, val := range AccessibleRoles {
		if val == role {
			return true
		}
	}

	return false
}
