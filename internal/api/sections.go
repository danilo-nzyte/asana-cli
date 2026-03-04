package api

import (
	"encoding/json"
	"fmt"

	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
)

// SectionsAPI provides methods for Asana section operations.
type SectionsAPI struct {
	client *client.Client
}

// NewSectionsAPI creates a new SectionsAPI.
func NewSectionsAPI(c *client.Client) *SectionsAPI {
	return &SectionsAPI{client: c}
}

// Create creates a new section in a project.
func (a *SectionsAPI) Create(projectGID string, req *models.SectionCreateRequest) (*models.Section, error) {
	body, err := a.client.Post(fmt.Sprintf("/projects/%s/sections", projectGID), req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Section]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Get retrieves a section by GID.
func (a *SectionsAPI) Get(gid string) (*models.Section, error) {
	body, err := a.client.Get(fmt.Sprintf("/sections/%s", gid))
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Section]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// List returns sections for a project.
func (a *SectionsAPI) List(projectGID string) ([]models.Section, error) {
	body, err := a.client.Get(fmt.Sprintf("/projects/%s/sections", projectGID))
	if err != nil {
		return nil, err
	}
	var resp models.AsanaListResponse[models.Section]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return resp.Data, nil
}

// Update updates a section.
func (a *SectionsAPI) Update(gid string, req *models.SectionUpdateRequest) (*models.Section, error) {
	body, err := a.client.Put(fmt.Sprintf("/sections/%s", gid), req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Section]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Delete deletes a section.
func (a *SectionsAPI) Delete(gid string) error {
	_, err := a.client.Delete(fmt.Sprintf("/sections/%s", gid))
	return err
}

// AddTask adds a task to a section.
func (a *SectionsAPI) AddTask(sectionGID string, taskGID string) error {
	req := &models.AddTaskRequest{Task: taskGID}
	_, err := a.client.Post(fmt.Sprintf("/sections/%s/addTask", sectionGID), req)
	return err
}
