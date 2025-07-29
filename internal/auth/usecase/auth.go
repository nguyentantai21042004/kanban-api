package usecase

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/auth"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/encrypter"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (uc *usecase) Login(ctx context.Context, sc models.Scope, ip auth.LoginInput) (auth.LoginOutput, error) {
	// Get user by email
	userModel, err := uc.userUC.GetByEmail(ctx, sc, ip.Email)
	if err != nil {
		uc.l.Errorf(ctx, "internal.auth.usecase.Login.uc.userUC.GetByEmail: %v", err)
		return auth.LoginOutput{}, auth.ErrInvalidCredentials
	}

	// Check if user is active
	if !userModel.IsActive {
		return auth.LoginOutput{}, auth.ErrUnauthorized
	}

	// Verify password
	if !encrypter.CheckPasswordHash(ip.Password, userModel.Password) {
		uc.l.Warnf(ctx, "internal.auth.usecase.Login.password_mismatch: %v", "password does not match")
		return auth.LoginOutput{}, auth.ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, refreshToken, err := uc.generateTokens(userModel.ID, userModel.Email, userModel.RoleID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.auth.usecase.Login.generateTokens: %v", err)
		return auth.LoginOutput{}, err
	}

	return auth.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: auth.UserInfo{
			ID:       userModel.ID,
			Email:    userModel.Email,
			FullName: userModel.FullName,
			Role:     userModel.RoleID,
		},
	}, nil
}

func (uc *usecase) RefreshToken(ctx context.Context, sc models.Scope, ip auth.RefreshTokenInput) (auth.RefreshTokenOutput, error) {
	// TODO: Implement refresh token validation and generation
	// For now, return error
	return auth.RefreshTokenOutput{}, auth.ErrInvalidToken
}

func (uc *usecase) Logout(ctx context.Context, sc models.Scope) error {
	// TODO: Implement logout (invalidate tokens)
	// For now, just return success
	return nil
}

// Helper function to generate JWT tokens
func (uc *usecase) generateTokens(userID, email, roleID string) (string, string, error) {
	// TODO: Implement JWT token generation
	// For now, return placeholder tokens
	accessToken := "access_token_" + postgres.NewUUID()
	refreshToken := "refresh_token_" + postgres.NewUUID()

	return accessToken, refreshToken, nil
}
