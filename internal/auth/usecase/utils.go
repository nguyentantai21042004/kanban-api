package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nguyentantai21042004/kanban-api/pkg/scope"
)

// Helper function to generate JWT tokens
func (uc *implUseCase) generateTokens(ctx context.Context, userID, username string) (string, error) {
	assToken, err := uc.scopeUC.CreateToken(scope.Payload{
		StandardClaims: jwt.StandardClaims{
			Audience: "kanban-api",
			// ExpiresAt: util.Now().Add(accessExpiry).Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer:   "kanban-api",
			// NotBefore: time.Now().Unix(),
			Subject: userID,
		},
		UserID:   userID,
		Username: username,
		Type:     "access",
		Refresh:  false,
	})
	if err != nil {
		uc.l.Errorf(ctx, "internal.auth.usecase.Login.scope.CreateToken: %v", err)
		return "", err
	}

	return assToken, nil
}
