package models

// Portfolio represents an Asana portfolio.
type Portfolio struct {
	GID          string `json:"gid"`
	Name         string `json:"name"`
	Color        string `json:"color,omitempty"`
	Owner        *Ref   `json:"owner,omitempty"`
	Workspace    *Ref   `json:"workspace,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
}

// PortfolioCreateRequest holds params for creating a portfolio.
type PortfolioCreateRequest struct {
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
	Color     string `json:"color,omitempty"`
}

// PortfolioUpdateRequest holds params for updating a portfolio.
type PortfolioUpdateRequest struct {
	Name  *string `json:"name,omitempty"`
	Color *string `json:"color,omitempty"`
}

// PortfolioItemRequest holds params for adding/removing portfolio items.
type PortfolioItemRequest struct {
	Item string `json:"item"`
}
