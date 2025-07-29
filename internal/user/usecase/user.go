package usecase

import (
	"context"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/pkg/encrypter"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (uc *usecase) Detail(ctx context.Context, sc models.Scope, ID string) (user.UserOutput, error) {
	// Check if user is Super Admin or accessing their own profile
	if sc.Role != "super_admin" && sc.UserID != ID {
		return user.UserOutput{}, user.ErrUnauthorized
	}

	userModel, err := uc.repo.Detail(ctx, sc, ID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.Detail.uc.repo.Detail: %v", err)
		return user.UserOutput{}, err
	}

	return user.UserOutput{User: userModel}, nil
}

func (uc *usecase) DetailMe(ctx context.Context, sc models.Scope) (user.UserOutput, error) {
	userModel, err := uc.repo.Detail(ctx, sc, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.DetailMe.uc.repo.Detail: %v", err)
		return user.UserOutput{}, err
	}

	return user.UserOutput{User: userModel}, nil
}

func (uc *usecase) UpdateProfile(ctx context.Context, sc models.Scope, ip user.UpdateProfileInput) (user.UserOutput, error) {
	// Only allow users to update their own profile
	userModel, err := uc.repo.Detail(ctx, sc, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.UpdateProfile.uc.repo.Detail: %v", err)
		return user.UserOutput{}, err
	}

	// Update fields
	userModel.FullName = ip.FullName
	if ip.AvatarURL != "" {
		userModel.AvatarURL = ip.AvatarURL
	}
	userModel.UpdatedAt = time.Now()

	// Save to database
	updatedUser, err := uc.repo.Update(ctx, sc, user.UpdateOptions{User: userModel})
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.UpdateProfile.uc.repo.Update: %v", err)
		return user.UserOutput{}, err
	}

	return user.UserOutput{User: updatedUser}, nil
}

func (uc *usecase) Create(ctx context.Context, sc models.Scope, ip user.CreateInput) (user.UserOutput, error) {
	// Only Super Admin can create users
	if sc.Role != "super_admin" {
		return user.UserOutput{}, user.ErrUnauthorized
	}

	// Check if user already exists
	existingUser, err := uc.repo.GetByEmail(ctx, sc, ip.Email)
	if err == nil && existingUser.ID != "" {
		return user.UserOutput{}, user.ErrUserExists
	}

	// Hash password
	hashedPassword, err := encrypter.HashPassword(ip.Password)
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.Create.encrypter.HashPassword: %v", err)
		return user.UserOutput{}, err
	}

	// Create user model
	userModel := models.User{
		ID:        postgres.NewUUID(),
		Email:     ip.Email,
		Password:  hashedPassword,
		FullName:  ip.FullName,
		RoleID:    ip.RoleID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	createdUser, err := uc.repo.Create(ctx, sc, user.CreateOptions{User: userModel})
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.Create.uc.repo.Create: %v", err)
		return user.UserOutput{}, err
	}

	return user.UserOutput{User: createdUser}, nil
}

func (uc *usecase) GetByEmail(ctx context.Context, sc models.Scope, email string) (models.User, error) {
	userModel, err := uc.repo.GetByEmail(ctx, sc, email)
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.GetByEmail.uc.repo.GetByEmail: %v", err)
		return models.User{}, err
	}

	return userModel, nil
}
