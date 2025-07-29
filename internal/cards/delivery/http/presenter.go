package http

import (
	"errors"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

type cardItem struct {
	ID          string              `json:"id"`
	ListID      string              `json:"list_id"`
	Title       string              `json:"title"`
	Description string              `json:"description,omitempty"`
	Position    int                 `json:"position"`
	DueDate     *time.Time          `json:"due_date,omitempty"`
	Priority    models.CardPriority `json:"priority"`
	Labels      []string            `json:"labels,omitempty"`
	IsArchived  bool                `json:"is_archived"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	DeletedAt   *time.Time          `json:"deleted_at,omitempty"`
}

// Get
type getReq struct {
	IDs       []string `form:"ids[]"`
	ListID    string   `form:"list_id"`
	Keyword   string   `form:"keyword"`
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

	return nil
}

func (req getReq) toInput() cards.GetInput {
	return cards.GetInput{
		Filter: cards.Filter{
			IDs:     req.IDs,
			ListID:  req.ListID,
			Keyword: req.Keyword,
		},
		PagQuery: req.PageQuery,
	}
}

type getCardResp struct {
	Items []cardItem                  `json:"items"`
	Meta  paginator.PaginatorResponse `json:"meta"`
}

func (h handler) newGetResp(o cards.GetOutput) getCardResp {
	items := make([]cardItem, len(o.Cards))
	for i, c := range o.Cards {
		items[i] = cardItem{
			ID:          c.ID,
			ListID:      c.ListID,
			Title:       c.Title,
			Description: c.Description,
			Position:    c.Position,
			DueDate:     c.DueDate,
			Priority:    c.Priority,
			Labels:      c.Labels,
			IsArchived:  c.IsArchived,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
			DeletedAt:   c.DeletedAt,
		}
	}
	return getCardResp{
		Items: items,
		Meta:  o.Pagination.ToResponse(),
	}
}

// Create
type createReq struct {
	ListID      string              `json:"list_id"`
	Title       string              `json:"title"`
	Description string              `json:"description,omitempty"`
	Priority    models.CardPriority `json:"priority,omitempty"`
	Labels      []string            `json:"labels,omitempty"`
	DueDate     *time.Time          `json:"due_date,omitempty"`
}

func (req createReq) toInput() cards.CreateInput {
	return cards.CreateInput{
		ListID:      req.ListID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Labels:      req.Labels,
		DueDate:     req.DueDate,
	}
}

func (h handler) newItem(o cards.DetailOutput) cardItem {
	item := cardItem{
		ID:          o.Card.ID,
		ListID:      o.Card.ListID,
		Title:       o.Card.Title,
		Description: o.Card.Description,
		Position:    o.Card.Position,
		DueDate:     o.Card.DueDate,
		Priority:    o.Card.Priority,
		Labels:      o.Card.Labels,
		IsArchived:  o.Card.IsArchived,
		CreatedAt:   o.Card.CreatedAt,
		UpdatedAt:   o.Card.UpdatedAt,
		DeletedAt:   o.Card.DeletedAt,
	}
	return item
}

// Update
type updateReq struct {
	ID          string               `json:"id"`
	Title       *string              `json:"title,omitempty"`
	Description *string              `json:"description,omitempty"`
	Priority    *models.CardPriority `json:"priority,omitempty"`
	Labels      *[]string            `json:"labels,omitempty"`
	DueDate     **time.Time          `json:"due_date,omitempty"`
}

func (req updateReq) toInput() cards.UpdateInput {
	return cards.UpdateInput{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Labels:      req.Labels,
		DueDate:     req.DueDate,
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

// Move
type moveReq struct {
	ID       string `json:"id"`
	ListID   string `json:"list_id"`
	Position int    `json:"position"`
}

func (req moveReq) validate() error {
	if err := postgres.IsUUID(req.ID); err != nil {
		return errors.New("invalid id")
	}
	if err := postgres.IsUUID(req.ListID); err != nil {
		return errors.New("invalid list_id")
	}
	if req.Position < 0 {
		return errors.New("invalid position")
	}
	return nil
}

func (req moveReq) toInput() cards.MoveInput {
	return cards.MoveInput{
		ID:       req.ID,
		ListID:   req.ListID,
		Position: float64(req.Position),
	}
}

// GetActivities
type getActivitiesReq struct {
	CardID    string `form:"card_id"`
	PageQuery paginator.PaginateQuery
}

func (req getActivitiesReq) validate() error {
	if err := postgres.IsUUID(req.CardID); err != nil {
		return errors.New("invalid card_id")
	}
	return nil
}

func (req getActivitiesReq) toInput() cards.GetActivitiesInput {
	return cards.GetActivitiesInput{
		CardID: req.CardID,
	}
}

type cardActivityItem struct {
	ID         string                `json:"id"`
	CardID     string                `json:"card_id"`
	ActionType models.CardActionType `json:"action_type"`
	OldData    map[string]any        `json:"old_data,omitempty"`
	NewData    map[string]any        `json:"new_data,omitempty"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
	DeletedAt  *time.Time            `json:"deleted_at,omitempty"`
}

type getCardActivitiesResp struct {
	Items []cardActivityItem          `json:"items"`
	Meta  paginator.PaginatorResponse `json:"meta"`
}

func (h handler) newGetActivitiesResp(o cards.GetActivitiesOutput) getCardActivitiesResp {
	items := make([]cardActivityItem, len(o.Activities))
	for i, a := range o.Activities {
		items[i] = cardActivityItem{
			ID:         a.ID,
			CardID:     a.CardID,
			ActionType: a.ActionType,
			OldData:    a.OldData,
			NewData:    a.NewData,
			CreatedAt:  a.CreatedAt,
			UpdatedAt:  a.UpdatedAt,
			DeletedAt:  a.DeletedAt,
		}
	}
	return getCardActivitiesResp{
		Items: items,
		Meta:  paginator.PaginatorResponse{}, // Activities không có pagination
	}
}
