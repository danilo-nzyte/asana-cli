package models

// Section represents an Asana section.
type Section struct {
	GID          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	Project      *Ref   `json:"project,omitempty"`
}

// SectionCreateRequest holds params for creating a section.
type SectionCreateRequest struct {
	Name string `json:"name"`
}

// SectionUpdateRequest holds params for updating a section.
type SectionUpdateRequest struct {
	Name *string `json:"name,omitempty"`
}

// AddTaskRequest holds params for adding a task to a section.
type AddTaskRequest struct {
	Task string `json:"task"`
}
