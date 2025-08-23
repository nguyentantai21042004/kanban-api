package usecase

import (
	"context"
	"strings"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/admin"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
)

func (uc implUsecase) Users(ctx context.Context, sc models.Scope, ip admin.UsersInput) (admin.UsersOutput, error) {
	// Defaults
	page := ip.Page
	perPage := ip.PerPage
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	// Reuse existing user list
	users, err := uc.userUC.List(ctx, sc, user.ListInput{})
	if err != nil {
		return admin.UsersOutput{}, err
	}

	// simple search by email/full_name (using Username as email field per current model)
	items := make([]admin.UserItem, 0)
	for _, u := range users {
		if ip.Search != "" {
			hay := strings.ToLower(u.Username + " " + u.FullName)
			if !strings.Contains(hay, strings.ToLower(ip.Search)) {
				continue
			}
		}
		// map role
		r := admin.RoleItem{}
		if u.RoleID != "" {
			if rl, rlErr := uc.roleUC.Detail(ctx, sc, u.RoleID); rlErr == nil {
				r = admin.RoleItem{ID: rl.ID, Name: rl.Name, Alias: rl.Alias}
			} else {
				r = admin.RoleItem{ID: u.RoleID}
			}
		}
		var lastLogin *string
		items = append(items, admin.UserItem{
			ID:          u.ID,
			Email:       u.Username,
			FullName:    u.FullName,
			Role:        r,
			IsActive:    u.IsActive,
			CreatedAt:   u.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   u.UpdatedAt.UTC().Format(time.RFC3339),
			LastLoginAt: lastLogin,
		})
	}

	total := int64(len(items))
	start := (page - 1) * perPage
	end := start + perPage
	if start > len(items) {
		start = len(items)
	}
	if end > len(items) {
		end = len(items)
	}
	pageItems := items[start:end]

	var out admin.UsersOutput
	out.Items = pageItems
	out.Meta.Count = int64(len(pageItems))
	out.Meta.CurrentPage = page
	out.Meta.PerPage = perPage
	out.Meta.Total = total
	out.Meta.TotalPages = int((total + int64(perPage) - 1) / int64(perPage))
	return out, nil
}

func (uc implUsecase) CreateUser(ctx context.Context, sc models.Scope, ip admin.CreateUserInput) (admin.UserItem, error) {
	// If no role is specified, default to "user" role
	roleID := ip.RoleID
	if roleID == "" {
		roles, _ := uc.roleUC.List(ctx, sc, role.ListInput{})
		for _, rr := range roles {
			if strings.EqualFold(rr.Alias, "user") {
				roleID = rr.ID
				break
			}
		}
	}
	// Set default password if not provided
	password := ip.Password
	if password == "" {
		password = "password123" // Default password - user should change it on first login
	}
	
	uo, err := uc.userUC.Create(ctx, sc, user.CreateInput{Username: ip.Email, Password: password, FullName: ip.FullName, RoleID: roleID})
	if err != nil {
		return admin.UserItem{}, err
	}
	// Map role
	r := admin.RoleItem{}
	if roleID != "" {
		if rl, err := uc.roleUC.Detail(ctx, sc, roleID); err == nil {
			r = admin.RoleItem{ID: rl.ID, Name: rl.Name, Alias: rl.Alias}
		} else {
			r = admin.RoleItem{ID: roleID}
		}
	}
	return admin.UserItem{
		ID:        uo.User.ID,
		Email:     uo.User.Username,
		FullName:  uo.User.FullName,
		Role:      r,
		IsActive:  uo.User.IsActive,
		CreatedAt: uo.User.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: uo.User.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

func (uc implUsecase) UpdateUser(ctx context.Context, sc models.Scope, id string, ip admin.UpdateUserInput) (admin.UserItem, error) {
	// Load user
	uo, err := uc.userUC.Detail(ctx, sc, id)
	if err != nil {
		return admin.UserItem{}, err
	}
	// Apply changes
	if ip.FullName != nil {
		uo.User.FullName = *ip.FullName
	}
	if ip.IsActive != nil {
		uo.User.IsActive = *ip.IsActive
	}
	// Resolve role if provided
	if ip.RoleID != nil || ip.RoleAlias != nil {
		var roleID string
		if ip.RoleID != nil {
			roleID = *ip.RoleID
		}
		if roleID == "" && ip.RoleAlias != nil && *ip.RoleAlias != "" {
			roles, _ := uc.roleUC.List(ctx, sc, role.ListInput{})
			for _, rr := range roles {
				if strings.EqualFold(rr.Alias, *ip.RoleAlias) {
					roleID = rr.ID
					break
				}
			}
		}
		uo.User.RoleID = roleID
	}
	// Persist minimal: use UpdateProfile for name
	_, err = uc.userUC.UpdateProfile(ctx, sc, user.UpdateProfileInput{FullName: uo.User.FullName})
	if err != nil {
		return admin.UserItem{}, err
	}
	// Reload
	uo, err = uc.userUC.Detail(ctx, sc, id)
	if err != nil {
		return admin.UserItem{}, err
	}
	// Map role
	r := admin.RoleItem{}
	if uo.User.RoleID != "" {
		if rl, err := uc.roleUC.Detail(ctx, sc, uo.User.RoleID); err == nil {
			r = admin.RoleItem{ID: rl.ID, Name: rl.Name, Alias: rl.Alias}
		} else {
			r = admin.RoleItem{ID: uo.User.RoleID}
		}
	}
	return admin.UserItem{
		ID:        uo.User.ID,
		Email:     uo.User.Username,
		FullName:  uo.User.FullName,
		Role:      r,
		IsActive:  uo.User.IsActive,
		CreatedAt: uo.User.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: uo.User.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}
