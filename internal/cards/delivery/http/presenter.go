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
	ID             string                 `json:"id"`
	ListID         string                 `json:"list_id"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description,omitempty"`
	Position       float64                `json:"position"`
	DueDate        *time.Time             `json:"due_date,omitempty"`
	Priority       models.CardPriority    `json:"priority"`
	Labels         []string               `json:"labels,omitempty"`
	IsArchived     bool                   `json:"is_archived"`
	AssignedTo     *string                `json:"assigned_to,omitempty"`
	Attachments    []string               `json:"attachments,omitempty"`
	EstimatedHours *float64               `json:"estimated_hours,omitempty"`
	ActualHours    *float64               `json:"actual_hours,omitempty"`
	StartDate      *time.Time             `json:"start_date,omitempty"`
	CompletionDate *time.Time             `json:"completion_date,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Checklist      []models.ChecklistItem `json:"checklist,omitempty"`
	LastActivityAt *time.Time             `json:"last_activity_at,omitempty"`
	CreatedBy      *string                `json:"created_by,omitempty"`
	UpdatedBy      *string                `json:"updated_by,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      *time.Time             `json:"deleted_at,omitempty"`
}

// Get
type getReq struct {
	IDs                []string `form:"ids[]"`
	ListID             string   `form:"list_id"`
	Keyword            string   `form:"keyword"`
	AssignedTo         string   `form:"assigned_to"`
	Priority           string   `form:"priority"`
	Tags               []string `form:"tags[]"`
	DueDateFrom        string   `form:"due_date_from"`
	DueDateTo          string   `form:"due_date_to"`
	StartDateFrom      string   `form:"start_date_from"`
	StartDateTo        string   `form:"start_date_to"`
	CompletionDateFrom string   `form:"completion_date_from"`
	CompletionDateTo   string   `form:"completion_date_to"`
	PageQuery          paginator.PaginateQuery
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
	filter := cards.Filter{
		IDs:        req.IDs,
		ListID:     req.ListID,
		Keyword:    req.Keyword,
		AssignedTo: req.AssignedTo,
		Tags:       req.Tags,
	}

	// Parse priority if provided
	if req.Priority != "" {
		filter.Priority = models.CardPriority(req.Priority)
	}

	// Parse date filters
	if req.DueDateFrom != "" {
		if dueDateFrom, err := time.Parse("2006-01-02", req.DueDateFrom); err == nil {
			filter.DueDateFrom = &dueDateFrom
		}
	}
	if req.DueDateTo != "" {
		if dueDateTo, err := time.Parse("2006-01-02", req.DueDateTo); err == nil {
			filter.DueDateTo = &dueDateTo
		}
	}
	if req.StartDateFrom != "" {
		if startDateFrom, err := time.Parse("2006-01-02", req.StartDateFrom); err == nil {
			filter.StartDateFrom = &startDateFrom
		}
	}
	if req.StartDateTo != "" {
		if startDateTo, err := time.Parse("2006-01-02", req.StartDateTo); err == nil {
			filter.StartDateTo = &startDateTo
		}
	}
	if req.CompletionDateFrom != "" {
		if completionDateFrom, err := time.Parse("2006-01-02", req.CompletionDateFrom); err == nil {
			filter.CompletionDateFrom = &completionDateFrom
		}
	}
	if req.CompletionDateTo != "" {
		if completionDateTo, err := time.Parse("2006-01-02", req.CompletionDateTo); err == nil {
			filter.CompletionDateTo = &completionDateTo
		}
	}

	return cards.GetInput{
		Filter:   filter,
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
			ID:             c.ID,
			ListID:         c.ListID,
			Title:          c.Title,
			Description:    c.Description,
			Position:       c.Position,
			DueDate:        c.DueDate,
			Priority:       c.Priority,
			Labels:         c.Labels,
			IsArchived:     c.IsArchived,
			AssignedTo:     c.AssignedTo,
			Attachments:    c.Attachments,
			EstimatedHours: c.EstimatedHours,
			ActualHours:    c.ActualHours,
			StartDate:      c.StartDate,
			CompletionDate: c.CompletionDate,
			Tags:           c.Tags,
			Checklist:      c.Checklist,
			LastActivityAt: c.LastActivityAt,
			CreatedBy:      c.CreatedBy,
			UpdatedBy:      c.UpdatedBy,
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
			DeletedAt:      c.DeletedAt,
		}
	}
	return getCardResp{
		Items: items,
		Meta:  o.Pagination.ToResponse(),
	}
}

// Create
type createReq struct {
	ListID         string                 `json:"list_id"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description,omitempty"`
	Priority       models.CardPriority    `json:"priority,omitempty"`
	Labels         []string               `json:"labels,omitempty"`
	DueDate        *time.Time             `json:"due_date,omitempty"`
	AssignedTo     *string                `json:"assigned_to,omitempty"`
	EstimatedHours *float64               `json:"estimated_hours,omitempty"`
	StartDate      *time.Time             `json:"start_date,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Checklist      []models.ChecklistItem `json:"checklist,omitempty"`
}

func (req createReq) toInput() cards.CreateInput {
	var assignedTo *string
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		assignedTo = req.AssignedTo
	}

	return cards.CreateInput{
		ListID:         req.ListID,
		Title:          req.Title,
		Description:    req.Description,
		Priority:       req.Priority,
		Labels:         req.Labels,
		DueDate:        req.DueDate,
		AssignedTo:     assignedTo,
		EstimatedHours: req.EstimatedHours,
		StartDate:      req.StartDate,
		Tags:           req.Tags,
		Checklist:      req.Checklist,
	}
}

func (h handler) newItem(o cards.DetailOutput) cardItem {
	item := cardItem{
		ID:             o.Card.ID,
		ListID:         o.Card.ListID,
		Title:          o.Card.Title,
		Description:    o.Card.Description,
		Position:       o.Card.Position,
		DueDate:        o.Card.DueDate,
		Priority:       o.Card.Priority,
		Labels:         o.Card.Labels,
		IsArchived:     o.Card.IsArchived,
		AssignedTo:     o.Card.AssignedTo,
		Attachments:    o.Card.Attachments,
		EstimatedHours: o.Card.EstimatedHours,
		ActualHours:    o.Card.ActualHours,
		StartDate:      o.Card.StartDate,
		CompletionDate: o.Card.CompletionDate,
		Tags:           o.Card.Tags,
		Checklist:      o.Card.Checklist,
		LastActivityAt: o.Card.LastActivityAt,
		CreatedBy:      o.Card.CreatedBy,
		UpdatedBy:      o.Card.UpdatedBy,
		CreatedAt:      o.Card.CreatedAt,
		UpdatedAt:      o.Card.UpdatedAt,
		DeletedAt:      o.Card.DeletedAt,
	}
	return item
}

// Update
type updateReq struct {
	ID             string                  `json:"id"`
	Title          *string                 `json:"title,omitempty"`
	Description    *string                 `json:"description,omitempty"`
	Priority       *models.CardPriority    `json:"priority,omitempty"`
	Labels         *[]string               `json:"labels,omitempty"`
	DueDate        **time.Time             `json:"due_date,omitempty"`
	AssignedTo     *string                 `json:"assigned_to,omitempty"`
	EstimatedHours *float64                `json:"estimated_hours,omitempty"`
	ActualHours    *float64                `json:"actual_hours,omitempty"`
	StartDate      *time.Time              `json:"start_date,omitempty"`
	CompletionDate *time.Time              `json:"completion_date,omitempty"`
	Tags           *[]string               `json:"tags,omitempty"`
	Checklist      *[]models.ChecklistItem `json:"checklist,omitempty"`
}

func (req updateReq) toInput() cards.UpdateInput {
	return cards.UpdateInput{
		ID:             req.ID,
		Title:          req.Title,
		Description:    req.Description,
		Priority:       req.Priority,
		Labels:         req.Labels,
		DueDate:        req.DueDate,
		AssignedTo:     req.AssignedTo,
		EstimatedHours: req.EstimatedHours,
		ActualHours:    req.ActualHours,
		StartDate:      req.StartDate,
		CompletionDate: req.CompletionDate,
		Tags:           req.Tags,
		Checklist:      req.Checklist,
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

// Enhanced functionality request types

// Assign
type assignReq struct {
	CardID     string `json:"card_id"`
	AssignedTo string `json:"assigned_to"`
}

func (req assignReq) toInput() cards.AssignInput {
	return cards.AssignInput{
		CardID:     req.CardID,
		AssignedTo: req.AssignedTo,
	}
}

// Unassign
type unassignReq struct {
	CardID string `json:"card_id"`
}

func (req unassignReq) toInput() cards.UnassignInput {
	return cards.UnassignInput{
		CardID: req.CardID,
	}
}

// AddAttachment
type addAttachmentReq struct {
	CardID       string `json:"card_id"`
	AttachmentID string `json:"attachment_id"`
}

func (req addAttachmentReq) toInput() cards.AddAttachmentInput {
	return cards.AddAttachmentInput{
		CardID:       req.CardID,
		AttachmentID: req.AttachmentID,
	}
}

// RemoveAttachment
type removeAttachmentReq struct {
	CardID       string `json:"card_id"`
	AttachmentID string `json:"attachment_id"`
}

func (req removeAttachmentReq) toInput() cards.RemoveAttachmentInput {
	return cards.RemoveAttachmentInput{
		CardID:       req.CardID,
		AttachmentID: req.AttachmentID,
	}
}

// UpdateTimeTracking
type updateTimeTrackingReq struct {
	CardID         string   `json:"card_id"`
	EstimatedHours *float64 `json:"estimated_hours,omitempty"`
	ActualHours    *float64 `json:"actual_hours,omitempty"`
}

func (req updateTimeTrackingReq) toInput() cards.UpdateTimeTrackingInput {
	return cards.UpdateTimeTrackingInput{
		CardID:         req.CardID,
		EstimatedHours: req.EstimatedHours,
		ActualHours:    req.ActualHours,
	}
}

// UpdateChecklist
type updateChecklistReq struct {
	CardID    string                 `json:"card_id"`
	Checklist []models.ChecklistItem `json:"checklist"`
}

func (req updateChecklistReq) toInput() cards.UpdateChecklistInput {
	return cards.UpdateChecklistInput{
		CardID:    req.CardID,
		Checklist: req.Checklist,
	}
}

// AddTag
type addTagReq struct {
	CardID string `json:"card_id"`
	Tag    string `json:"tag"`
}

func (req addTagReq) toInput() cards.AddTagInput {
	return cards.AddTagInput{
		CardID: req.CardID,
		Tag:    req.Tag,
	}
}

// RemoveTag
type removeTagReq struct {
	CardID string `json:"card_id"`
	Tag    string `json:"tag"`
}

func (req removeTagReq) toInput() cards.RemoveTagInput {
	return cards.RemoveTagInput{
		CardID: req.CardID,
		Tag:    req.Tag,
	}
}

// SetStartDate
type setStartDateReq struct {
	CardID    string     `json:"card_id"`
	StartDate *time.Time `json:"start_date,omitempty"`
}

func (req setStartDateReq) toInput() cards.SetStartDateInput {
	return cards.SetStartDateInput{
		CardID:    req.CardID,
		StartDate: req.StartDate,
	}
}

// SetCompletionDate
type setCompletionDateReq struct {
	CardID         string     `json:"card_id"`
	CompletionDate *time.Time `json:"completion_date,omitempty"`
}

func (req setCompletionDateReq) toInput() cards.SetCompletionDateInput {
	return cards.SetCompletionDateInput{
		CardID:         req.CardID,
		CompletionDate: req.CompletionDate,
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
