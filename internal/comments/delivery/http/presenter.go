package http

import (
	"errors"

	"gitlab.com/tantai-kanban/kanban-api/internal/comments"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

type respObj struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type commentItem struct {
	ID        string   `json:"id"`
	CardID    string   `json:"card_id"`
	Content   string   `json:"content"`
	ParentID  *string  `json:"parent_id,omitempty"`
	IsEdited  *bool    `json:"is_edited,omitempty"`
	EditedAt  *string  `json:"edited_at,omitempty"`
	EditedBy  *respObj `json:"edited_by,omitempty"`
	User      respObj  `json:"user"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// Get
type getReq struct {
	IDs       []string `form:"ids[]"`
	Keyword   string   `form:"keyword"`
	CardID    string   `form:"card_id"`
	UserID    string   `form:"user_id"`
	ParentID  string   `form:"parent_id"`
	PageQuery paginator.PaginateQuery
}

func (req getReq) validate() error {
	if len(req.IDs) > 0 {
		for _, id := range req.IDs {
			if err := postgres.IsUUID(id); err != nil {
				return errors.New("invalid id")
			}
		}
	}

	if req.CardID != "" {
		if err := postgres.IsUUID(req.CardID); err != nil {
			return errors.New("invalid card_id")
		}
	}

	if req.UserID != "" {
		if err := postgres.IsUUID(req.UserID); err != nil {
			return errors.New("invalid user_id")
		}
	}

	if req.ParentID != "" {
		if err := postgres.IsUUID(req.ParentID); err != nil {
			return errors.New("invalid parent_id")
		}
	}

	return nil
}

func (req getReq) toInput() comments.GetInput {
	return comments.GetInput{
		Filter: comments.Filter{
			IDs:      req.IDs,
			Keyword:  req.Keyword,
			CardID:   req.CardID,
			UserID:   req.UserID,
			ParentID: req.ParentID,
		},
		PagQuery: req.PageQuery,
	}
}

type getCommentResp struct {
	Items []commentItem               `json:"items"`
	Meta  paginator.PaginatorResponse `json:"meta"`
}

func (h handler) newGetResp(o comments.GetOutput) getCommentResp {
	userMap := make(map[string]models.User)
	for _, u := range o.Users {
		userMap[u.ID] = u
	}

	items := make([]commentItem, len(o.Comments))
	for i, c := range o.Comments {
		items[i] = commentItem{
			ID:       c.ID,
			CardID:   c.CardID,
			Content:  c.Content,
			ParentID: c.ParentID,
			IsEdited: c.IsEdited,
			User: respObj{
				ID:   userMap[c.UserID].ID,
				Name: userMap[c.UserID].FullName,
			},
			CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if c.EditedAt != nil {
			editedAt := c.EditedAt.Format("2006-01-02T15:04:05Z07:00")
			items[i].EditedAt = &editedAt
		}

		if c.EditedBy != nil {
			if user, exists := userMap[*c.EditedBy]; exists {
				items[i].EditedBy = &respObj{
					ID:   user.ID,
					Name: user.FullName,
				}
			}
		}
	}
	return getCommentResp{
		Items: items,
		Meta:  o.Pagination.ToResponse(),
	}
}

// Create
type createReq struct {
	CardID   string  `json:"card_id"`
	Content  string  `json:"content"`
	ParentID *string `json:"parent_id,omitempty"`
}

func (req createReq) toInput() comments.CreateInput {
	return comments.CreateInput{
		CardID:   req.CardID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}
}

func (h handler) newItem(o comments.DetailOutput) commentItem {
	item := commentItem{
		ID:       o.Comment.ID,
		CardID:   o.Comment.CardID,
		Content:  o.Comment.Content,
		ParentID: o.Comment.ParentID,
		IsEdited: o.Comment.IsEdited,
		User: respObj{
			ID:   o.User.ID,
			Name: o.User.FullName,
		},
		CreatedAt: o.Comment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: o.Comment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if o.Comment.EditedAt != nil {
		editedAt := o.Comment.EditedAt.Format("2006-01-02T15:04:05Z07:00")
		item.EditedAt = &editedAt
	}

	if o.Comment.EditedBy != nil {
		item.EditedBy = &respObj{
			ID:   o.User.ID,
			Name: o.User.FullName,
		}
	}

	return item
}

// Update
type updateReq struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func (req updateReq) toInput() comments.UpdateInput {
	return comments.UpdateInput{
		ID:      req.ID,
		Content: req.Content,
	}
}

// Delete
type deleteReq struct {
	IDs []string `json:"ids[]"`
}

func (req deleteReq) validate() error {
	if len(req.IDs) > 0 {
		for _, id := range req.IDs {
			if err := postgres.IsUUID(id); err != nil {
				return errors.New("invalid id")
			}
		}
	}

	return nil
}

// GetByCard
type getByCardReq struct {
	CardID string `uri:"card_id"`
}

func (req getByCardReq) validate() error {
	if err := postgres.IsUUID(req.CardID); err != nil {
		return errors.New("invalid card_id")
	}
	return nil
}
