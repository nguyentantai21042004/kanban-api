package http

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type getReq struct {
	Page  int    `form:"page" binding:"omitempty,min=1"`
	Limit int64  `form:"limit" binding:"omitempty,min=1,max=100"`
	Name  string `form:"name"`
	Code  string `form:"code"`
}

func (req getReq) toInput() role.GetInput {
	filter := role.Filter{}
	if req.Name != "" {
		filter.Name = &req.Name
	}
	if req.Code != "" {
		filter.Code = &req.Code
	}

	return role.GetInput{
		Filter: filter,
		PagQuery: paginator.PaginateQuery{
			Page:  req.Page,
			Limit: req.Limit,
		},
	}
}

type getRoleResp struct {
	Data []roleItem                  `json:"data"`
	Meta paginator.PaginatorResponse `json:"meta"`
}

type roleItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (h handler) newGetResp(o role.GetOutput) getRoleResp {
	items := make([]roleItem, len(o.Roles))
	for i, role := range o.Roles {
		items[i] = roleItem{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Alias:       role.Alias,
			Description: role.Description,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return getRoleResp{
		Data: items,
		Meta: paginator.PaginatorResponse{
			Total:      o.Paginator.Total,
			TotalPages: o.Paginator.TotalPages(),
		},
	}
}

func (h handler) newItem(o role.DetailOutput) roleItem {
	return roleItem{
		ID:          o.Role.ID,
		Name:        o.Role.Name,
		Code:        o.Role.Code,
		Alias:       o.Role.Alias,
		Description: o.Role.Description,
		CreatedAt:   o.Role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   o.Role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
