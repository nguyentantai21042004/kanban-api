package usecase

import (
	"context"
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/admin"
	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	"github.com/nguyentantai21042004/kanban-api/internal/cards"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/internal/user"
)

func (uc implUsecase) Dashboard(ctx context.Context, sc models.Scope, ip admin.DashboardInput) (admin.DashboardOutput, error) {
	// Normalize period
	period := ip.Period
	if period == "" {
		period = "7d"
	}
	days := 7
	switch period {
	case "30d":
		days = 30
	case "90d":
		days = 90
	}

	end := time.Now().UTC()
	start := end.AddDate(0, 0, -days+1)
	prevStart := start.AddDate(0, 0, -days)
	prevEnd := start.AddDate(0, 0, -1)

	// Cards dashboard
	cardsDash, err := uc.cardUC.Dashboard(ctx, sc, cards.DashboardInput{From: start, To: end})
	if err != nil {
		return admin.DashboardOutput{}, err
	}

	// Boards dashboard (totals)
	boardsDash, err := uc.boardUC.Dashboard(ctx, sc, boards.DashboardInput{From: start, To: end})
	if err != nil {
		return admin.DashboardOutput{}, err
	}
	// Active boards from cards activity
	boardsDash.Active = int64(len(cardsDash.ActiveBoardIDs))

	// Users dashboard (total)
	usersDash, err := uc.userUC.Dashboard(ctx, sc, user.DashboardInput{From: start, To: end})
	if err != nil {
		return admin.DashboardOutput{}, err
	}
	// Active/growth from cards dashboard
	usersDash.Active = int64(len(cardsDash.ActiveUserIDs))
	prevCardsDash, err := uc.cardUC.Dashboard(ctx, sc, cards.DashboardInput{From: prevStart, To: prevEnd})
	if err != nil {
		return admin.DashboardOutput{}, err
	}
	prevActive := int64(len(prevCardsDash.ActiveUserIDs))
	growth := 0.0
	if prevActive > 0 {
		growth = (float64(usersDash.Active-prevActive) / float64(prevActive)) * 100.0
	} else if usersDash.Active > 0 {
		growth = 100.0
	}
	usersDash.Growth = growth

	var out admin.DashboardOutput
	out.Users.Total = usersDash.Total
	out.Users.Active = usersDash.Active
	out.Users.Growth = usersDash.Growth
	out.Boards.Total = boardsDash.Total
	out.Boards.Active = boardsDash.Active
	out.Cards.Total = cardsDash.Total
	out.Cards.Completed = cardsDash.Completed
	out.Cards.Overdue = cardsDash.Overdue
	for _, ap := range cardsDash.Activity {
		out.Activity = append(out.Activity, admin.ActivityPoint{
			Date:           ap.Date,
			CardsCreated:   ap.CardsCreated,
			CardsCompleted: ap.CardsCompleted,
		})
	}
	return out, nil
}
