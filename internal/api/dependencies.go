package api

import (
	"encoding/json"
	"fmt"

	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
)

// DependenciesAPI provides methods for Asana task dependency operations.
type DependenciesAPI struct {
	client *client.Client
}

// NewDependenciesAPI creates a new DependenciesAPI.
func NewDependenciesAPI(c *client.Client) *DependenciesAPI {
	return &DependenciesAPI{client: c}
}

// Add adds dependencies to a task.
func (a *DependenciesAPI) Add(taskGID string, dependsOnGIDs []string) error {
	payload := map[string]interface{}{
		"dependencies": dependsOnGIDs,
	}
	_, err := a.client.Post(fmt.Sprintf("/tasks/%s/addDependencies", taskGID), payload)
	return err
}

// Remove removes a dependency from a task.
func (a *DependenciesAPI) Remove(taskGID string, dependsOnGIDs []string) error {
	payload := map[string]interface{}{
		"dependencies": dependsOnGIDs,
	}
	_, err := a.client.Post(fmt.Sprintf("/tasks/%s/removeDependencies", taskGID), payload)
	return err
}

// List returns dependencies of a task.
func (a *DependenciesAPI) List(taskGID string) ([]models.Task, error) {
	body, err := a.client.Get(fmt.Sprintf("/tasks/%s/dependencies", taskGID))
	if err != nil {
		return nil, err
	}
	var resp models.AsanaListResponse[models.Task]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return resp.Data, nil
}
