package http

import (
	"errors"

	"gitlab.com/tantai-kanban/kanban-api/internal/labels"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

type labelItem struct {
	ID      string `json:"id"`
	BoardID string `json:"board_id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

// Get
type getReq struct {
	IDs       []string `form:"ids[]"`
	BoardID   string   `form:"board_id"`
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

func (req getReq) toInput() labels.GetInput {
	return labels.GetInput{
		Filter: labels.Filter{
			IDs:     req.IDs,
			BoardID: req.BoardID,
			Keyword: req.Keyword,
		},
		PagQuery: req.PageQuery,
	}
}

type getLabelResp struct {
	Items []labelItem                 `json:"items"`
	Meta  paginator.PaginatorResponse `json:"meta"`
}

func (h handler) newGetResp(o labels.GetOutput) getLabelResp {
	items := make([]labelItem, len(o.Labels))
	for i, l := range o.Labels {
		items[i] = labelItem{
			ID:      l.ID,
			BoardID: l.BoardID,
			Name:    l.Name,
			Color:   l.Color,
		}
	}
	return getLabelResp{
		Items: items,
		Meta:  o.Pagination.ToResponse(),
	}
}

// Create
type createReq struct {
	BoardID string `json:"board_id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

func (req createReq) toInput() labels.CreateInput {
	return labels.CreateInput{
		BoardID: req.BoardID,
		Name:    req.Name,
		Color:   req.Color,
	}
}

func (h handler) newItem(o labels.DetailOutput) labelItem {
	item := labelItem{
		ID:      o.Label.ID,
		BoardID: o.Label.BoardID,
		Name:    o.Label.Name,
		Color:   o.Label.Color,
	}
	return item
}

// Update
type updateReq struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (req updateReq) toInput() labels.UpdateInput {
	return labels.UpdateInput{
		ID:    req.ID,
		Name:  req.Name,
		Color: req.Color,
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
