package api

import (
	"encoding/json"
	"fmt"

	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
)

// AttachmentsAPI provides methods for Asana attachment operations.
type AttachmentsAPI struct {
	client *client.Client
}

// NewAttachmentsAPI creates a new AttachmentsAPI.
func NewAttachmentsAPI(c *client.Client) *AttachmentsAPI {
	return &AttachmentsAPI{client: c}
}

// Upload uploads a file as an attachment to a task.
func (a *AttachmentsAPI) Upload(taskGID, filePath string) (*models.Attachment, error) {
	body, err := a.client.PostMultipart(fmt.Sprintf("/tasks/%s/attachments", taskGID), filePath)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Attachment]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Get retrieves an attachment by GID.
func (a *AttachmentsAPI) Get(gid string) (*models.Attachment, error) {
	body, err := a.client.Get(fmt.Sprintf("/attachments/%s", gid))
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Attachment]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// List returns attachments for a task.
func (a *AttachmentsAPI) List(taskGID string) ([]models.Attachment, error) {
	body, err := a.client.Get(fmt.Sprintf("/tasks/%s/attachments", taskGID))
	if err != nil {
		return nil, err
	}
	var resp models.AsanaListResponse[models.Attachment]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return resp.Data, nil
}

// Delete deletes an attachment.
func (a *AttachmentsAPI) Delete(gid string) error {
	_, err := a.client.Delete(fmt.Sprintf("/attachments/%s", gid))
	return err
}
