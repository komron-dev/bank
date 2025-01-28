package grpc_api

import (
	"context"
	db "github.com/komron-dev/bank/db/sqlc"
	"github.com/komron-dev/bank/pb"
	"github.com/komron-dev/bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) VerifyEmail(ctx context.Context, request *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    request.GetEmailId(),
		SecretCode: request.GetSecretCode(),
	})
	if err != nil {
		return nil, err
	}
	response := &pb.VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}

	return response, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}

	if err := val.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}
