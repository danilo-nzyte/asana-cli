package models

// Comment represents an Asana story (comment).
type Comment struct {
	GID          string `json:"gid"`
	Text         string `json:"text"`
	HTMLText     string `json:"html_text,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	CreatedBy    *Ref   `json:"created_by,omitempty"`
	Target       *Ref   `json:"target,omitempty"`
	Type         string `json:"type,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
}

// CommentCreateRequest holds params for creating a comment.
type CommentCreateRequest struct {
	Text string `json:"text"`
}

// CommentUpdateRequest holds params for updating a comment.
type CommentUpdateRequest struct {
	Text *string `json:"text,omitempty"`
}
