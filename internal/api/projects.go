package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
)

// ProjectsAPI provides methods for Asana project operations.
type ProjectsAPI struct {
	client *client.Client
}

// NewProjectsAPI creates a new ProjectsAPI.
func NewProjectsAPI(c *client.Client) *ProjectsAPI {
	return &ProjectsAPI{client: c}
}

// Create creates a new project.
func (a *ProjectsAPI) Create(req *models.ProjectCreateRequest) (*models.Project, error) {
	body, err := a.client.Post("/projects", req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Project]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Get retrieves a project by GID.
func (a *ProjectsAPI) Get(gid string) (*models.Project, error) {
	body, err := a.client.Get(fmt.Sprintf("/projects/%s", gid))
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Project]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// List returns projects, optionally filtered by workspace, team, and archived status.
func (a *ProjectsAPI) List(workspace, team string, archived *bool) ([]models.Project, error) {
	params := url.Values{}
	if workspace != "" {
		params.Set("workspace", workspace)
	}
	if team != "" {
		params.Set("team", team)
	}
	if archived != nil {
		params.Set("archived", fmt.Sprintf("%t", *archived))
	}

	path := "/projects"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	body, err := a.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaListResponse[models.Project]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return resp.Data, nil
}

// Update updates a project.
func (a *ProjectsAPI) Update(gid string, req *models.ProjectUpdateRequest) (*models.Project, error) {
	body, err := a.client.Put(fmt.Sprintf("/projects/%s", gid), req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Project]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Delete deletes a project.
func (a *ProjectsAPI) Delete(gid string) error {
	_, err := a.client.Delete(fmt.Sprintf("/projects/%s", gid))
	return err
}
