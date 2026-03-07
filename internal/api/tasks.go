package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
)

// TasksAPI provides methods for Asana task operations.
type TasksAPI struct {
	client *client.Client
}

// NewTasksAPI creates a new TasksAPI.
func NewTasksAPI(c *client.Client) *TasksAPI {
	return &TasksAPI{client: c}
}

// Create creates a new task.
func (a *TasksAPI) Create(req *models.TaskCreateRequest) (*models.Task, error) {
	body, err := a.client.Post("/tasks", req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Task]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Get retrieves a task by GID.
func (a *TasksAPI) Get(gid string) (*models.Task, error) {
	body, err := a.client.Get(fmt.Sprintf("/tasks/%s", gid))
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Task]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// List returns tasks for a project, optionally filtered.
// Pass optFields="" to use Asana defaults.
func (a *TasksAPI) List(project string, completed *bool, assignee string, optFields string) ([]models.Task, error) {
	params := url.Values{}
	if project != "" {
		params.Set("project", project)
	}
	if completed != nil {
		params.Set("completed_since", "now")
		if *completed {
			params.Del("completed_since")
		}
	}
	if assignee != "" {
		params.Set("assignee", assignee)
	}
	if optFields != "" {
		params.Set("opt_fields", optFields)
	}

	path := "/tasks"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	body, err := a.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaListResponse[models.Task]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return resp.Data, nil
}

// Update updates a task.
func (a *TasksAPI) Update(gid string, req *models.TaskUpdateRequest) (*models.Task, error) {
	body, err := a.client.Put(fmt.Sprintf("/tasks/%s", gid), req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Task]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Delete deletes a task.
func (a *TasksAPI) Delete(gid string) error {
	_, err := a.client.Delete(fmt.Sprintf("/tasks/%s", gid))
	return err
}

// Search searches for tasks in a workspace.
// Pass optFields="" to use Asana defaults.
func (a *TasksAPI) Search(workspace string, query string, project string, assignee string, optFields string) ([]models.Task, error) {
	params := url.Values{}
	if query != "" {
		params.Set("text", query)
	}
	if project != "" {
		params.Set("projects.any", project)
	}
	if assignee != "" {
		params.Set("assignee.any", assignee)
	}
	if optFields != "" {
		params.Set("opt_fields", optFields)
	}

	path := fmt.Sprintf("/workspaces/%s/tasks/search", workspace)
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	body, err := a.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaListResponse[models.Task]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return resp.Data, nil
}

// MyTasks returns incomplete tasks assigned to a user, sorted by due date.
func (a *TasksAPI) MyTasks(workspace, assignee, project, optFields string) ([]models.Task, error) {
	params := url.Values{}
	params.Set("assignee.any", assignee)
	params.Set("is_subtask", "false")
	params.Set("completed", "false")
	params.Set("sort_by", "due_on")
	params.Set("sort_ascending", "true")
	if project != "" {
		params.Set("projects.any", project)
	}
	if optFields != "" {
		params.Set("opt_fields", optFields)
	}

	path := fmt.Sprintf("/workspaces/%s/tasks/search?%s", workspace, params.Encode())

	body, err := a.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaListResponse[models.Task]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return resp.Data, nil
}
