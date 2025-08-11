package http

import (
	"errors"

	"gitlab.com/tantai-kanban/kanban-api/internal/lists"
	"gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
	"gitlab.com/tantai-kanban/kanban-api/pkg/postgres"
)

type listItem struct {
	ID       string `json:"id"`
	BoardID  string `json:"board_id"`
	Name     string `json:"name"`
	Position string `json:"position"`
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

func (req getReq) toInput() lists.GetInput {
	return lists.GetInput{
		Filter: lists.Filter{
			IDs:     req.IDs,
			BoardID: req.BoardID,
			Keyword: req.Keyword,
		},
		PagQuery: req.PageQuery,
	}
}

type getListResp struct {
	Items []listItem                  `json:"items"`
	Meta  paginator.PaginatorResponse `json:"meta"`
}

func (h handler) newGetResp(o lists.GetOutput) getListResp {
	items := make([]listItem, len(o.Lists))
	for i, l := range o.Lists {
		items[i] = listItem{
			ID:       l.ID,
			BoardID:  l.BoardID,
			Name:     l.Name,
			Position: l.Position,
		}
	}
	return getListResp{
		Items: items,
		Meta:  o.Pagination.ToResponse(),
	}
}

// Create
type createReq struct {
	BoardID  string `json:"board_id"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

func (req createReq) toInput() lists.CreateInput {
	return lists.CreateInput{
		BoardID:  req.BoardID,
		Name:     req.Name,
		Position: req.Position,
	}
}

func (h handler) newItem(o lists.DetailOutput) listItem {
	item := listItem{
		ID:       o.List.ID,
		BoardID:  o.List.BoardID,
		Name:     o.List.Name,
		Position: o.List.Position,
	}
	return item
}

// Update
type updateReq struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

func (req updateReq) toInput() lists.UpdateInput {
	return lists.UpdateInput{
		ID:       req.ID,
		Name:     req.Name,
		Position: req.Position,
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
