package models

import (
	"time"
)

type Upload struct {
	ID            string     `json:"id"`
	BucketName    string     `json:"bucket_name"`
	ObjectName    string     `json:"object_name"`
	OriginalName  string     `json:"original_name"`
	Size          int64      `json:"size"`
	ContentType   string     `json:"content_type"`
	Etag          string     `json:"etag,omitempty"`
	Metadata      string     `json:"metadata,omitempty"`
	URL           string     `json:"url,omitempty"`
	Source        string     `json:"source"`
	PublicID      string     `json:"public_id,omitempty"`
	CreatedUserID string     `json:"created_user_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}
