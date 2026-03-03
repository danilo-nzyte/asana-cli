package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
)

// PortfoliosAPI provides methods for Asana portfolio operations.
type PortfoliosAPI struct {
	client *client.Client
}

// NewPortfoliosAPI creates a new PortfoliosAPI.
func NewPortfoliosAPI(c *client.Client) *PortfoliosAPI {
	return &PortfoliosAPI{client: c}
}

// Create creates a new portfolio.
func (a *PortfoliosAPI) Create(req *models.PortfolioCreateRequest) (*models.Portfolio, error) {
	body, err := a.client.Post("/portfolios", req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Portfolio]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Get retrieves a portfolio by GID.
func (a *PortfoliosAPI) Get(gid string) (*models.Portfolio, error) {
	body, err := a.client.Get(fmt.Sprintf("/portfolios/%s", gid))
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Portfolio]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// List returns portfolios, optionally filtered by owner and workspace.
func (a *PortfoliosAPI) List(workspace, owner string) ([]models.Portfolio, error) {
	params := url.Values{}
	if workspace != "" {
		params.Set("workspace", workspace)
	}
	if owner != "" {
		params.Set("owner", owner)
	}

	path := "/portfolios"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	body, err := a.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaListResponse[models.Portfolio]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return resp.Data, nil
}

// Update updates a portfolio.
func (a *PortfoliosAPI) Update(gid string, req *models.PortfolioUpdateRequest) (*models.Portfolio, error) {
	body, err := a.client.Put(fmt.Sprintf("/portfolios/%s", gid), req)
	if err != nil {
		return nil, err
	}
	var resp models.AsanaResponse[models.Portfolio]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp.Data, nil
}

// Delete deletes a portfolio.
func (a *PortfoliosAPI) Delete(gid string) error {
	_, err := a.client.Delete(fmt.Sprintf("/portfolios/%s", gid))
	return err
}

// AddItem adds a project to a portfolio.
func (a *PortfoliosAPI) AddItem(gid, item string) error {
	_, err := a.client.Post(fmt.Sprintf("/portfolios/%s/addItem", gid), &models.PortfolioItemRequest{Item: item})
	return err
}

// RemoveItem removes a project from a portfolio.
func (a *PortfoliosAPI) RemoveItem(gid, item string) error {
	_, err := a.client.Post(fmt.Sprintf("/portfolios/%s/removeItem", gid), &models.PortfolioItemRequest{Item: item})
	return err
}
