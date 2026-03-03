package client

import (
	"encoding/json"
	"fmt"

	"github.com/danilodrobac/asana-cli/internal/models"
)

// PaginatedResponse holds a page of results with optional next page info.
type PaginatedResponse struct {
	Data     json.RawMessage  `json:"data"`
	NextPage *models.NextPage `json:"next_page,omitempty"`
}

// CollectAll fetches all pages for a given path, returning combined raw JSON arrays.
func (c *Client) CollectAll(path string, params map[string]string) ([]json.RawMessage, error) {
	var all []json.RawMessage
	offset := ""

	for {
		p := make(map[string]string)
		for k, v := range params {
			p[k] = v
		}
		if offset != "" {
			p["offset"] = offset
		}

		query := ""
		if len(p) > 0 {
			parts := ""
			for k, v := range p {
				if parts != "" {
					parts += "&"
				}
				parts += k + "=" + v
			}
			query = "?" + parts
		}

		body, err := c.Get(fmt.Sprintf("%s%s", path, query))
		if err != nil {
			return nil, err
		}

		var page PaginatedResponse
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("failed to parse paginated response: %w", err)
		}

		// Parse the data array
		var items []json.RawMessage
		if err := json.Unmarshal(page.Data, &items); err != nil {
			return nil, fmt.Errorf("failed to parse data array: %w", err)
		}
		all = append(all, items...)

		if page.NextPage == nil || page.NextPage.Offset == "" {
			break
		}
		offset = page.NextPage.Offset
	}

	return all, nil
}
