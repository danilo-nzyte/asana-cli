package api

import (
	"encoding/json"
	"fmt"

	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
)

// CustomFieldsAPI provides methods for Asana custom field operations.
type CustomFieldsAPI struct {
	client *client.Client
}

// NewCustomFieldsAPI creates a new CustomFieldsAPI.
func NewCustomFieldsAPI(c *client.Client) *CustomFieldsAPI {
	return &CustomFieldsAPI{client: c}
}

// Create creates a new custom field.
func (a *CustomFieldsAPI) Create(req *models.CustomFieldCreateRequest) (*models.CustomField, error) {
	body, err := a.client.Post("/custom_fields", req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.CustomField]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Get retrieves a custom field by GID.
func (a *CustomFieldsAPI) Get(gid string) (*models.CustomField, error) {
	body, err := a.client.Get(fmt.Sprintf("/custom_fields/%s", gid))
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.CustomField]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// List returns custom fields for a workspace.
func (a *CustomFieldsAPI) List(workspace string) ([]models.CustomField, error) {
	path := fmt.Sprintf("/workspaces/%s/custom_fields", workspace)

	body, err := a.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaListResponse[models.CustomField]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return resp.Data, nil
}

// Update updates a custom field.
func (a *CustomFieldsAPI) Update(gid string, req *models.CustomFieldUpdateRequest) (*models.CustomField, error) {
	body, err := a.client.Put(fmt.Sprintf("/custom_fields/%s", gid), req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.CustomField]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Delete deletes a custom field.
func (a *CustomFieldsAPI) Delete(gid string) error {
	_, err := a.client.Delete(fmt.Sprintf("/custom_fields/%s", gid))
	return err
}
