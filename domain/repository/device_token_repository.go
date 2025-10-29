package repository

import (
	"context"

	"github.com/suk-chanthea/ezra/domain/entity"
)

type DeviceTokenRepository interface {
	// Save or update a device token
	Save(ctx context.Context, token *entity.DeviceToken) error
	
	// Get all active tokens for a specific user
	GetActiveTokensByUserID(ctx context.Context, userID uint) ([]string, error)
	
	// Get all active tokens for users in a specific band
	GetTokensByBandID(ctx context.Context, bandID uint) ([]string, error)
	
	// Get all active tokens in the system (for broadcast)
	GetAllActiveTokens(ctx context.Context) ([]string, error)
	
	// Get all active tokens excluding specific user (for broadcast excluding sender)
	GetAllActiveTokensExcept(ctx context.Context, excludeUserID uint) ([]string, error)
	
	// Delete a specific token
	DeleteToken(ctx context.Context, token string) error
	
	// Delete multiple tokens (for cleanup of invalid tokens)
	DeleteTokens(ctx context.Context, tokens []string) error
	
	// Deactivate a token
	DeactivateToken(ctx context.Context, token string) error
	
	// Delete all tokens for a user
	DeleteUserTokens(ctx context.Context, userID uint) error
}

