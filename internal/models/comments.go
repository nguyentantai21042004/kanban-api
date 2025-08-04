package models

import (
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/dbmodels"
)

type Comment struct {
	ID        string     `json:"id"`
	CardID    string     `json:"card_id"`
	UserID    string     `json:"user_id"`
	Content   string     `json:"content"`
	ParentID  *string    `json:"parent_id,omitempty"`
	IsEdited  *bool      `json:"is_edited,omitempty"`
	EditedAt  *time.Time `json:"edited_at,omitempty"`
	EditedBy  *string    `json:"edited_by,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewComment(dbComment dbmodels.Comment) Comment {
	var parentID *string
	if dbComment.ParentID.Valid {
		parentID = &dbComment.ParentID.String
	}

	var isEdited *bool
	if dbComment.IsEdited.Valid {
		isEdited = &dbComment.IsEdited.Bool
	}

	var editedAt *time.Time
	if dbComment.EditedAt.Valid {
		editedAt = &dbComment.EditedAt.Time
	}

	var editedBy *string
	if dbComment.EditedBy.Valid {
		editedBy = &dbComment.EditedBy.String
	}

	var deletedAt *time.Time
	if dbComment.DeletedAt.Valid {
		deletedAt = &dbComment.DeletedAt.Time
	}

	return Comment{
		ID:        dbComment.ID,
		CardID:    dbComment.CardID,
		UserID:    dbComment.UserID,
		Content:   dbComment.Content,
		ParentID:  parentID,
		IsEdited:  isEdited,
		EditedAt:  editedAt,
		EditedBy:  editedBy,
		CreatedAt: dbComment.CreatedAt,
		UpdatedAt: dbComment.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
