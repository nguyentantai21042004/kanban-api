package http

import (
	"errors"
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/cards"
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	"github.com/nguyentantai21042004/kanban-api/pkg/paginator"
	"github.com/nguyentantai21042004/kanban-api/pkg/postgres"
	"github.com/nguyentantai21042004/kanban-api/pkg/response"
	"github.com/nguyentantai21042004/kanban-api/pkg/util"
)

type respObj struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type cardItem struct {
	ID             string                 `json:"id"`
	Board          respObj                `json:"board"`
	List           respObj                `json:"list"`
	Name           string                 `json:"name"`
	Alias          string                 `json:"alias"`
	Description    string                 `json:"description,omitempty"`
	Position       string                 `json:"position"`
	DueDate        *response.DateTime     `json:"due_date,omitempty"`
	Priority       models.CardPriority    `json:"priority"`
	Labels         []string               `json:"labels,omitempty"`
	IsArchived     bool                   `json:"is_archived"`
	AssignedTo     *string                `json:"assigned_to,omitempty"`
	Attachments    []string               `json:"attachments,omitempty"`
	EstimatedHours *float64               `json:"estimated_hours,omitempty"`
	ActualHours    *float64               `json:"actual_hours,omitempty"`
	StartDate      *response.DateTime     `json:"start_date,omitempty"`
	CompletionDate *response.DateTime     `json:"completion_date,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Checklist      []models.ChecklistItem `json:"checklist,omitempty"`
	LastActivityAt *response.DateTime     `json:"last_activity_at,omitempty"`
	CreatedBy      *respObj               `json:"created_by,omitempty"`
	UpdatedBy      *respObj               `json:"updated_by,omitempty"`
	CreatedAt      response.DateTime      `json:"created_at"`
	UpdatedAt      response.DateTime      `json:"updated_at"`
}

// Get
type getReq struct {
	IDs                []string `form:"ids[]"`
	ListID             string   `form:"list_id"`
	BoardID            string   `form:"board_id"`
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
		BoardID:    req.BoardID,
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
			Name:           c.Name,
			Alias:          c.Alias,
			Description:    c.Description,
			Position:       c.Position,
			Priority:       c.Priority,
			Labels:         c.Labels,
			IsArchived:     c.IsArchived,
			AssignedTo:     c.AssignedTo,
			Attachments:    c.Attachments,
			EstimatedHours: c.EstimatedHours,
			ActualHours:    c.ActualHours,
			Tags:           c.Tags,
			Checklist:      c.Checklist,
			CreatedAt:      response.DateTime(c.CreatedAt),
			UpdatedAt:      response.DateTime(c.UpdatedAt),
		}

		if c.ListID != "" {
			items[i].List = respObj{
				ID:   c.ListID,
				Name: c.Name,
			}
		}

		if c.DueDate != nil {
			dueDate := response.DateTime(*c.DueDate)
			items[i].DueDate = &dueDate
		}

		if c.StartDate != nil {
			startDate := response.DateTime(*c.StartDate)
			items[i].StartDate = &startDate
		}

		if c.CompletionDate != nil {
			completionDate := response.DateTime(*c.CompletionDate)
			items[i].CompletionDate = &completionDate
		}

		if c.LastActivityAt != nil {
			lastActivityAt := response.DateTime(*c.LastActivityAt)
			items[i].LastActivityAt = &lastActivityAt
		}

		if c.CreatedBy != nil {
			items[i].CreatedBy = &respObj{
				ID:   *c.CreatedBy,
				Name: c.Name,
			}
		}

		if c.UpdatedBy != nil {
			items[i].UpdatedBy = &respObj{
				ID:   *c.UpdatedBy,
				Name: c.Name,
			}
		}
	}
	return getCardResp{
		Items: items,
		Meta:  o.Pagination.ToResponse(),
	}
}

// Create
type checkListItemReq struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	IsCompleted bool   `json:"is_completed"`
}

type createReq struct {
	BoardID        string              `json:"board_id" binding:"required"`
	ListID         string              `json:"list_id" binding:"required"`
	Name           string              `json:"name" binding:"required"`
	Description    string              `json:"description"`
	Priority       models.CardPriority `json:"priority"`
	Labels         []string            `json:"labels"`
	DueDate        string              `json:"due_date"`
	AssignedTo     *string             `json:"assigned_to"`
	EstimatedHours *float64            `json:"estimated_hours"`
	StartDate      string              `json:"start_date"`
	Tags           []string            `json:"tags"`
	Checklist      []checkListItemReq  `json:"checklist"`
}

func (req createReq) validate() error {
	if err := postgres.IsUUID(req.BoardID); err != nil {
		return errors.New("invalid board_id")
	}

	if err := postgres.IsUUID(req.ListID); err != nil {
		return errors.New("invalid list_id")
	}

	if req.DueDate != "" {
		if _, err := util.StrToDate(req.DueDate); err != nil {
			return errors.New("invalid due_date")
		}
	}

	if req.AssignedTo != nil {
		if *req.AssignedTo != "" {
			if err := postgres.IsUUID(*req.AssignedTo); err != nil {
				return errors.New("invalid assigned_to")
			}
		}
	}

	if req.StartDate != "" {
		if _, err := util.StrToDate(req.StartDate); err != nil {
			return errors.New("invalid start_date")
		}
	}

	if req.EstimatedHours != nil {
		if *req.EstimatedHours < 0 {
			return errors.New("estimated_hours must be greater than 0")
		}
	}

	return nil
}

func (req createReq) toInput() cards.CreateInput {
	var assignedTo *string
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		assignedTo = req.AssignedTo
	}

	dueDate, _ := util.StrToDate(req.DueDate)
	startDate, _ := util.StrToDate(req.StartDate)

	checklist := make([]cards.ChecklistItemInput, len(req.Checklist))
	for i, c := range req.Checklist {
		checklist[i] = cards.ChecklistItemInput{
			Content:     c.Content,
			IsCompleted: c.IsCompleted,
		}
	}

	return cards.CreateInput{
		BoardID:        req.BoardID,
		ListID:         req.ListID,
		Name:           req.Name,
		Description:    req.Description,
		Priority:       req.Priority,
		Labels:         req.Labels,
		DueDate:        &dueDate,
		AssignedTo:     assignedTo,
		EstimatedHours: req.EstimatedHours,
		StartDate:      &startDate,
		Tags:           req.Tags,
		Checklist:      checklist,
	}
}

func (h handler) newItem(o cards.DetailOutput) cardItem {
	item := cardItem{
		ID:             o.Card.ID,
		Name:           o.Card.Name,
		Alias:          o.Card.Alias,
		Description:    o.Card.Description,
		Position:       o.Card.Position,
		Priority:       o.Card.Priority,
		Labels:         o.Card.Labels,
		IsArchived:     o.Card.IsArchived,
		AssignedTo:     o.Card.AssignedTo,
		Attachments:    o.Card.Attachments,
		EstimatedHours: o.Card.EstimatedHours,
		ActualHours:    o.Card.ActualHours,
		Tags:           o.Card.Tags,
		Checklist:      o.Card.Checklist,
		CreatedAt:      response.DateTime(o.Card.CreatedAt),
		UpdatedAt:      response.DateTime(o.Card.UpdatedAt),
	}

	if o.Board.ID != "" {
		item.Board = respObj{
			ID:   o.Board.ID,
			Name: o.Board.Name,
		}
	}

	if o.Card.ListID != "" {
		item.List = respObj{
			ID:   o.List.ID,
			Name: o.List.Name,
		}
	}

	if o.Card.DueDate != nil {
		dueDate := response.DateTime(*o.Card.DueDate)
		item.DueDate = &dueDate
	}

	if o.Card.StartDate != nil {
		startDate := response.DateTime(*o.Card.StartDate)
		item.StartDate = &startDate
	}

	if o.Card.CompletionDate != nil {
		completionDate := response.DateTime(*o.Card.CompletionDate)
		item.CompletionDate = &completionDate
	}

	if o.Card.LastActivityAt != nil {
		lastActivityAt := response.DateTime(*o.Card.LastActivityAt)
		item.LastActivityAt = &lastActivityAt
	}

	if o.Card.CreatedBy != nil {
		item.CreatedBy = &respObj{
			ID:   *o.Card.CreatedBy,
			Name: o.Card.Name,
		}
	}

	if o.Card.UpdatedBy != nil {
		item.UpdatedBy = &respObj{
			ID:   *o.Card.UpdatedBy,
			Name: o.Card.Name,
		}
	}

	return item
}

// Update
type updateReq struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	Description    *string                 `json:"description"`
	Priority       *models.CardPriority    `json:"priority"`
	Labels         *[]string               `json:"labels"`
	DueDate        string                  `json:"due_date"`
	AssignedTo     *string                 `json:"assigned_to"`
	EstimatedHours *float64                `json:"estimated_hours"`
	ActualHours    *float64                `json:"actual_hours"`
	StartDate      string                  `json:"start_date"`
	CompletionDate *time.Time              `json:"completion_date"`
	Tags           *[]string               `json:"tags"`
	Checklist      *[]models.ChecklistItem `json:"checklist"`
}

func (req updateReq) toInput() cards.UpdateInput {
	dueDate, _ := util.StrToDate(req.DueDate)
	startDate, _ := util.StrToDate(req.StartDate)

	return cards.UpdateInput{
		ID:             req.ID,
		Name:           req.Name,
		Description:    req.Description,
		Priority:       req.Priority,
		Labels:         req.Labels,
		DueDate:        &dueDate,
		AssignedTo:     req.AssignedTo,
		EstimatedHours: req.EstimatedHours,
		ActualHours:    req.ActualHours,
		StartDate:      &startDate,
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
	AfterID  string `json:"after_id"`
	BeforeID string `json:"before_id"`
}

func (req moveReq) validate() error {
	if err := postgres.IsUUID(req.ID); err != nil {
		return errors.New("invalid id")
	}
	if err := postgres.IsUUID(req.ListID); err != nil {
		return errors.New("invalid list_id")
	}
	if req.AfterID != "" {
		if err := postgres.IsUUID(req.AfterID); err != nil {
			return errors.New("invalid after_id")
		}
	}
	if req.BeforeID != "" {
		if err := postgres.IsUUID(req.BeforeID); err != nil {
			return errors.New("invalid before_id")
		}
	}
	return nil
}

func (req moveReq) toInput() cards.MoveInput {
	return cards.MoveInput{
		ID:       req.ID,
		ListID:   req.ListID,
		AfterID:  req.AfterID,
		BeforeID: req.BeforeID,
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
	CreatedAt  response.DateTime     `json:"created_at"`
	UpdatedAt  response.DateTime     `json:"updated_at"`
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
			CreatedAt:  response.DateTime(a.CreatedAt),
			UpdatedAt:  response.DateTime(a.UpdatedAt),
		}
	}
	return getCardActivitiesResp{
		Items: items,
		Meta:  paginator.PaginatorResponse{},
	}
}
