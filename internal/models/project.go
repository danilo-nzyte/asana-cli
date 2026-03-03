package models

// Project represents an Asana project.
type Project struct {
	GID          string `json:"gid"`
	Name         string `json:"name"`
	Notes        string `json:"notes,omitempty"`
	Archived     bool   `json:"archived"`
	Color        string `json:"color,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	ModifiedAt   string `json:"modified_at,omitempty"`
	Workspace    *Ref   `json:"workspace,omitempty"`
	Team         *Ref   `json:"team,omitempty"`
	Owner        *Ref   `json:"owner,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
}

// ProjectCreateRequest holds params for creating a project.
type ProjectCreateRequest struct {
	Name      string `json:"name"`
	Workspace string `json:"workspace,omitempty"`
	Team      string `json:"team,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

// ProjectUpdateRequest holds params for updating a project.
type ProjectUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	Notes    *string `json:"notes,omitempty"`
	Archived *bool   `json:"archived,omitempty"`
	Color    *string `json:"color,omitempty"`
}
