package models

import (
	"time"

	"github.com/nguyentantai21042004/kanban-api/internal/dbmodels"
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
	return Comment{
		ID:        dbComment.ID,
		CardID:    dbComment.CardID,
		UserID:    dbComment.UserID,
		Content:   dbComment.Content,
		ParentID:  dbComment.ParentID.Ptr(),
		IsEdited:  dbComment.IsEdited.Ptr(),
		EditedAt:  dbComment.EditedAt.Ptr(),
		EditedBy:  dbComment.EditedBy.Ptr(),
		CreatedAt: dbComment.CreatedAt,
		UpdatedAt: dbComment.UpdatedAt,
		DeletedAt: dbComment.DeletedAt.Ptr(),
	}
}
