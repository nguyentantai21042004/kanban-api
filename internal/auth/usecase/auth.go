package usecase

import (
	"context"
	"sync"

	"github.com/nguyentantai21042004/kanban-api/internal/auth"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/internal/user"
	"github.com/nguyentantai21042004/kanban-api/pkg/encrypter"
)

func (uc *implUseCase) Login(ctx context.Context, sc models.Scope, ip auth.LoginInput) (auth.LoginOutput, error) {
	// Get user by email
	u, err := uc.userUC.GetOne(ctx, sc, user.GetOneInput{Username: ip.Username})
	if err != nil {
		uc.l.Errorf(ctx, "internal.auth.usecase.Login.uc.userUC.GetOne: %v", err)
		return auth.LoginOutput{}, auth.ErrInvalidCredentials
	}

	// Check if user is active
	if !u.IsActive {
		uc.l.Warnf(ctx, "internal.auth.usecase.Login.user_inactive: %v", "user is inactive")
		return auth.LoginOutput{}, auth.ErrUnauthorized
	}

	// Verify password
	if !encrypter.CheckPasswordHash(ip.Password, u.PasswordHash) {
		uc.l.Warnf(ctx, "internal.auth.usecase.Login.password_mismatch: %v", "password does not match")
		return auth.LoginOutput{}, auth.ErrInvalidCredentials
	}

	var (
		rl       models.Role
		assToken string
		errChan  = make(chan error, 2)
		wg       sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		assToken, err = uc.generateTokens(ctx, u.ID, u.Username)
		if err != nil {
			uc.l.Errorf(ctx, "internal.auth.usecase.Login.generateTokens: %v", err)
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		rl, err = uc.roleUC.Detail(ctx, sc, u.RoleID)
		if err != nil {
			uc.l.Errorf(ctx, "internal.auth.usecase.Login.roleUC.Detail: %v", err)
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)
	if err := <-errChan; err != nil {
		return auth.LoginOutput{}, err
	}

	return auth.LoginOutput{
		AssToken: assToken,
		User:     u,
		Role:     rl,
	}, nil
}

func (uc *implUseCase) RefreshToken(ctx context.Context, sc models.Scope, ip auth.RefreshTokenInput) (auth.RefreshTokenOutput, error) {
	// TODO: Implement refresh token validation and generation
	// For now, return error
	return auth.RefreshTokenOutput{}, auth.ErrInvalidToken
}

func (uc *implUseCase) Logout(ctx context.Context, sc models.Scope) error {
	// TODO: Implement logout (invalidate tokens)
	// For now, just return success
	return nil
}
