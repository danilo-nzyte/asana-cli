package models

// CustomField represents an Asana custom field.
type CustomField struct {
	GID          string       `json:"gid"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	EnumOptions  []EnumOption `json:"enum_options,omitempty"`
	Workspace    *Ref         `json:"workspace,omitempty"`
	ResourceType string       `json:"resource_type,omitempty"`
}

// EnumOption is a single option in an enum custom field.
type EnumOption struct {
	GID     string `json:"gid,omitempty"`
	Name    string `json:"name"`
	Color   string `json:"color,omitempty"`
	Enabled bool   `json:"enabled"`
}

// CustomFieldCreateRequest holds params for creating a custom field.
type CustomFieldCreateRequest struct {
	Name        string       `json:"name"`
	Workspace   string       `json:"workspace"`
	Type        string       `json:"resource_subtype"`
	EnumOptions []EnumOption `json:"enum_options,omitempty"`
}

// CustomFieldUpdateRequest holds params for updating a custom field.
type CustomFieldUpdateRequest struct {
	Name *string `json:"name,omitempty"`
}
