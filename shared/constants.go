package shared

import "github.com/google/uuid"

type (
	PostJson struct {
		Title     string                 `json:"title" binding:"required"`
		Body      string                 `json:"body" binding:"required"`
		CreatedAt int64                  `json:"createdAt"`
		ID        uuid.UUID              `json:"id"`
		Comments  map[string]CommentJson `json:"comments"`
	}
	CommentJson struct {
		Body      string `json:"body" binding:"required"`
		CreatedAt int64  `json:"createdAt"`
		ID        string `json:"id"`
	}
	Posts map[uuid.UUID]PostJson
)
