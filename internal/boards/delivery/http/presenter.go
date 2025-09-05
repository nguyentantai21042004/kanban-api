package http

import (
	"errors"

	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
	"github.com/nguyentantai21042004/kanban-api/pkg/postgres"
)

type respObj struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type boardItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Alias       string  `json:"alias,omitempty"`
	CreatedBy   respObj `json:"created_by"`
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
	userMap := make(map[string]models.User)
	for _, u := range o.Users {
		userMap[u.ID] = u
	}

	items := make([]boardItem, len(o.Boards))
	for i, b := range o.Boards {
		items[i] = boardItem{
			ID:    b.ID,
			Name:  b.Name,
			Alias: b.Alias,
			CreatedBy: respObj{
				ID:   userMap[*b.CreatedBy].ID,
				Name: userMap[*b.CreatedBy].FullName,
			},
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
	userMap := make(map[string]models.User)
	for _, u := range o.Users {
		userMap[u.ID] = u
	}

	item := boardItem{
		ID:    o.Board.ID,
		Name:  o.Board.Name,
		Alias: o.Board.Alias,
		CreatedBy: respObj{
			ID:   userMap[*o.Board.CreatedBy].ID,
			Name: userMap[*o.Board.CreatedBy].FullName,
		},
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
