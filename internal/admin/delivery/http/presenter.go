package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/internal/admin"
	"gitlab.com/tantai-kanban/kanban-api/pkg/response"
)

type dashboardResp struct {
	Users struct {
		Total  int64   `json:"total"`
		Active int64   `json:"active"`
		Growth float64 `json:"growth"`
	} `json:"users"`
	Boards struct {
		Total  int64 `json:"total"`
		Active int64 `json:"active"`
	} `json:"boards"`
	Cards struct {
		Total     int64 `json:"total"`
		Completed int64 `json:"completed"`
		Overdue   int64 `json:"overdue"`
	} `json:"cards"`
	Activity []admin.ActivityPoint `json:"activity"`
}

func newDashboardResp(o admin.DashboardOutput) dashboardResp {
	var r dashboardResp
	r.Users = o.Users
	r.Boards = o.Boards
	r.Cards = o.Cards
	r.Activity = o.Activity
	return r
}

func respondOK(c *gin.Context, data interface{}) {
	response.OK(c, data)
}

// Users presenter
type usersResp struct {
	Items []admin.UserItem `json:"items"`
	Meta  struct {
		Count       int64 `json:"count"`
		CurrentPage int   `json:"current_page"`
		PerPage     int   `json:"per_page"`
		Total       int64 `json:"total"`
		TotalPages  int   `json:"total_pages"`
	} `json:"meta"`
}

func newUsersResp(o admin.UsersOutput) usersResp {
	var r usersResp
	r.Items = o.Items
	r.Meta = o.Meta
	return r
}
