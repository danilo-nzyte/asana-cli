package api

import (
	"encoding/json"
	"fmt"

	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
)

// CommentsAPI provides methods for Asana comment (story) operations.
type CommentsAPI struct {
	client *client.Client
}

// NewCommentsAPI creates a new CommentsAPI.
func NewCommentsAPI(c *client.Client) *CommentsAPI {
	return &CommentsAPI{client: c}
}

// Create adds a comment to a task.
func (a *CommentsAPI) Create(taskGID string, req *models.CommentCreateRequest) (*models.Comment, error) {
	body, err := a.client.Post(fmt.Sprintf("/tasks/%s/stories", taskGID), req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Comment]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Get retrieves a comment by GID.
func (a *CommentsAPI) Get(gid string) (*models.Comment, error) {
	body, err := a.client.Get(fmt.Sprintf("/stories/%s", gid))
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Comment]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// List returns comments for a task.
func (a *CommentsAPI) List(taskGID string) ([]models.Comment, error) {
	body, err := a.client.Get(fmt.Sprintf("/tasks/%s/stories", taskGID))
	if err != nil {
		return nil, err
	}
	var resp models.AsanaListResponse[models.Comment]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return resp.Data, nil
}

// Update updates a comment.
func (a *CommentsAPI) Update(gid string, req *models.CommentUpdateRequest) (*models.Comment, error) {
	body, err := a.client.Put(fmt.Sprintf("/stories/%s", gid), req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Comment]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Delete deletes a comment.
func (a *CommentsAPI) Delete(gid string) error {
	_, err := a.client.Delete(fmt.Sprintf("/stories/%s", gid))
	return err
}
