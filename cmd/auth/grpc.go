package main

import (
	"context"
	"log"

	"github.com/marwan562/fintech-ecosystem/internal/auth"
	pb "github.com/marwan562/fintech-ecosystem/proto/auth"
)

type AuthGRPCServer struct {
	pb.UnimplementedAuthServiceServer
	repo *auth.Repository
}

func NewAuthGRPCServer(repo *auth.Repository) *AuthGRPCServer {
	return &AuthGRPCServer{repo: repo}
}

func (s *AuthGRPCServer) ValidateKey(ctx context.Context, req *pb.ValidateKeyRequest) (*pb.ValidateKeyResponse, error) {
	key, err := s.repo.GetAPIKeyByHash(ctx, req.KeyHash)
	if err != nil {
		log.Printf("GRPC ValidateKey error: %v", err)
		return &pb.ValidateKeyResponse{Valid: false}, nil
	}

	if key == nil || key.RevokedAt != nil {
		return &pb.ValidateKeyResponse{Valid: false}, nil
	}

	return &pb.ValidateKeyResponse{
		Valid:       true,
		UserId:      key.UserID,
		Environment: key.Environment,
		Scopes:      key.Scopes,
	}, nil
}

func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	token, err := s.repo.ValidateOAuthToken(ctx, req.AccessToken)
	if err != nil {
		// Log error only if it's not "not found" or "expired" to avoid noise?
		// Repo returns error on expiry.
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}
	if token == nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:     true,
		ClientId:  token.ClientID,
		UserId:    token.UserID,
		Scope:     token.Scope,
		ExpiresAt: token.ExpiresAt.Unix(),
	}, nil
}
