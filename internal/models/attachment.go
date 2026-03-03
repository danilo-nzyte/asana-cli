package models

// Attachment represents an Asana attachment.
type Attachment struct {
	GID          string `json:"gid"`
	Name         string `json:"name"`
	Host         string `json:"host,omitempty"`
	DownloadURL  string `json:"download_url,omitempty"`
	ViewURL      string `json:"view_url,omitempty"`
	Parent       *Ref   `json:"parent,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
}
