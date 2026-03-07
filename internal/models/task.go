package models

// Membership represents a task's membership in a project section.
type Membership struct {
	Project Ref `json:"project"`
	Section Ref `json:"section"`
}

// Task represents an Asana task.
type Task struct {
	GID            string                   `json:"gid"`
	Name           string                   `json:"name"`
	Notes          string                   `json:"notes,omitempty"`
	Completed      bool                     `json:"completed"`
	CompletedAt    string                   `json:"completed_at,omitempty"`
	DueOn          string                   `json:"due_on,omitempty"`
	DueAt          string                   `json:"due_at,omitempty"`
	StartOn        string                   `json:"start_on,omitempty"`
	Assignee       *Ref                     `json:"assignee,omitempty"`
	Projects       []Ref                    `json:"projects,omitempty"`
	Memberships    []Membership             `json:"memberships,omitempty"`
	Tags           []Ref                    `json:"tags,omitempty"`
	Dependencies   []Ref                    `json:"dependencies,omitempty"`
	Dependents     []Ref                    `json:"dependents,omitempty"`
	CustomFields   []map[string]interface{} `json:"custom_fields,omitempty"`
	CreatedAt      string                   `json:"created_at,omitempty"`
	ModifiedAt     string                   `json:"modified_at,omitempty"`
	Workspace      *Ref                     `json:"workspace,omitempty"`
	ResourceType   string                   `json:"resource_type,omitempty"`
}

// TaskCreateRequest holds params for creating a task.
type TaskCreateRequest struct {
	Name         string                 `json:"name"`
	Projects     []string               `json:"projects,omitempty"`
	Assignee     string                 `json:"assignee,omitempty"`
	Notes        string                 `json:"notes,omitempty"`
	DueOn        string                 `json:"due_on,omitempty"`
	Workspace    string                 `json:"workspace,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

// TaskUpdateRequest holds params for updating a task.
type TaskUpdateRequest struct {
	Name      *string `json:"name,omitempty"`
	Notes     *string `json:"notes,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
	DueOn     *string `json:"due_on,omitempty"`
	Assignee  *string `json:"assignee,omitempty"`
}
