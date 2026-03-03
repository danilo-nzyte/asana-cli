package models

// AsanaResponse wraps the standard Asana API response envelope.
type AsanaResponse[T any] struct {
	Data     T         `json:"data"`
	NextPage *NextPage `json:"next_page,omitempty"`
}

// AsanaListResponse wraps paginated list responses.
type AsanaListResponse[T any] struct {
	Data     []T       `json:"data"`
	NextPage *NextPage `json:"next_page,omitempty"`
}

// AsanaErrorResponse represents an error from the Asana API.
type AsanaErrorResponse struct {
	Errors []AsanaError `json:"errors"`
}

// AsanaError is a single error entry.
type AsanaError struct {
	Message string `json:"message"`
	Help    string `json:"help,omitempty"`
	Phrase  string `json:"phrase,omitempty"`
}

// Ref is a minimal reference to an Asana object (GID + name).
type Ref struct {
	GID  string `json:"gid"`
	Name string `json:"name,omitempty"`
}

// NextPage holds pagination cursor info.
type NextPage struct {
	Offset string `json:"offset"`
	Path   string `json:"path"`
	URI    string `json:"uri"`
}
