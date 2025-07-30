package usecase

import (
	"context"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/internal/user/repository"
	"gitlab.com/tantai-kanban/kanban-api/pkg/encrypter"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

func (uc *usecase) Detail(ctx context.Context, sc models.Scope, ID string) (user.UserOutput, error) {
	// Check if user is Super Admin or accessing their own profile
	u, err := uc.repo.Detail(ctx, sc, ID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.Detail.uc.repo.Detail: %v", err)
		return user.UserOutput{}, err
	}

	return user.UserOutput{User: u}, nil
}

func (uc *usecase) DetailMe(ctx context.Context, sc models.Scope) (user.UserOutput, error) {
	u, err := uc.repo.Detail(ctx, sc, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.DetailMe.uc.repo.Detail: %v", err)
		return user.UserOutput{}, err
	}

	return user.UserOutput{User: u}, nil
}

func (uc *usecase) List(ctx context.Context, sc models.Scope, ip user.ListInput) ([]models.User, error) {
	qr, err := uc.repo.List(ctx, sc, repository.ListOptions(ip))
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.List.uc.repo.List: %v", err)
		return nil, err
	}

	return qr, nil
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
	updatedUser, err := uc.repo.Update(ctx, sc, repository.UpdateOptions{User: userModel})
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.UpdateProfile.uc.repo.Update: %v", err)
		return user.UserOutput{}, err
	}

	return user.UserOutput{User: updatedUser}, nil
}

func (uc *usecase) Create(ctx context.Context, sc models.Scope, ip user.CreateInput) (user.UserOutput, error) {
	// Check if user already exists
	existingUser, err := uc.repo.GetOne(ctx, sc, repository.GetOneOptions{Username: ip.Username})
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
		ID:           postgres.NewUUID(),
		Username:     ip.Username,
		PasswordHash: hashedPassword,
		FullName:     ip.FullName,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to database
	createdUser, err := uc.repo.Create(ctx, sc, repository.CreateOptions{User: userModel})
	if err != nil {
		uc.l.Errorf(ctx, "internal.user.usecase.Create.uc.repo.Create: %v", err)
		return user.UserOutput{}, err
	}

	return user.UserOutput{User: createdUser}, nil
}

func (uc *usecase) GetOne(ctx context.Context, sc models.Scope, ip user.GetOneInput) (models.User, error) {
	u, err := uc.repo.GetOne(ctx, sc, repository.GetOneOptions{Username: ip.Username})
	if err != nil {
		if err == repository.ErrNotFound {
			uc.l.Warnf(ctx, "internal.user.usecase.GetOne.uc.repo.GetOne: %v", err)
			return models.User{}, user.ErrUserNotFound
		}
		uc.l.Errorf(ctx, "internal.user.usecase.GetOne.uc.repo.GetOne: %v", err)
		return models.User{}, err
	}

	return u, nil
}
