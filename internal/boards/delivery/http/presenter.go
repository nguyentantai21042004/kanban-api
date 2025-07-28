package http

import (
	"errors"

	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

type boardItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Alias       string `json:"alias"`
}

// Get
type getReq struct {
	IDs       []string `form:"ids[]"`
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

func (req getReq) toInput() boards.GetInput {
	return boards.GetInput{
		Filter: boards.Filter{
			IDs:     req.IDs,
			Keyword: req.Keyword,
		},
		PagQuery: req.PageQuery,
	}
}

type getBoardResp struct {
	Items []boardItem                 `json:"items"`
	Meta  paginator.PaginatorResponse `json:"meta"`
}

func (h handler) newGetResp(o boards.GetOutput) getBoardResp {
	items := make([]boardItem, len(o.Boards))
	for i, b := range o.Boards {
		items[i] = boardItem{
			ID:    b.ID,
			Name:  b.Name,
			Alias: b.Alias,
		}
		if b.Description != nil {
			items[i].Description = *b.Description
		}
	}
	return getBoardResp{
		Items: items,
		Meta:  o.Pagination.ToResponse(),
	}
}

// Create
type createReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (req createReq) toInput() boards.CreateInput {
	return boards.CreateInput{
		Name:        req.Name,
		Description: req.Description,
	}
}

func (h handler) newItem(o boards.DetailOutput) boardItem {
	item := boardItem{
		ID:    o.Board.ID,
		Name:  o.Board.Name,
		Alias: o.Board.Alias,
	}
	if o.Board.Description != nil {
		item.Description = *o.Board.Description
	}
	return item
}

// Update
type updateReq struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (req updateReq) toInput() boards.UpdateInput {
	return boards.UpdateInput{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
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
