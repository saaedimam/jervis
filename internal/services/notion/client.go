package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	baseURL = "https://api.notion.com/v1"
)

const (
	apiVersion = "2022-06-28"
)

// Client is a lightweight Notion API client.
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a new Notion API client.
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body any) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Notion-Version", apiVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("notion api error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// UpdatePageBlocks updates the blocks of a page.
func (c *Client) UpdatePageBlocks(ctx context.Context, pageID string, children []any) error {
	_, err := c.do(ctx, "PATCH", "/blocks/"+pageID+"/children", map[string]any{
		"children": children,
	})
	return err
}

// GetBlockChildren retrieves the children blocks of a block.
func (c *Client) GetBlockChildren(ctx context.Context, blockID string) ([]map[string]any, error) {
	resp, err := c.do(ctx, "GET", "/blocks/"+blockID+"/children", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Results []map[string]any `json:"results"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

// DeleteBlock deletes a block.
func (c *Client) DeleteBlock(ctx context.Context, blockID string) error {
	_, err := c.do(ctx, "DELETE", "/blocks/"+blockID, nil)
	return err
}

// CreatePageInDatabase creates a new page in a database.
func (c *Client) CreatePageInDatabase(ctx context.Context, databaseID string, properties map[string]any) (string, error) {
	resp, err := c.do(ctx, "POST", "/pages", map[string]any{
		"parent":     map[string]any{"database_id": databaseID},
		"properties": properties,
	})
	if err != nil {
		return "", err
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	return result.ID, nil
}

// UpdateDatabaseProperties updates the properties of a database.
func (c *Client) UpdateDatabaseProperties(ctx context.Context, databaseID string, properties map[string]any) error {
	_, err := c.do(ctx, "PATCH", "/databases/"+databaseID, map[string]any{
		"properties": properties,
	})
	return err
}

// QueryDatabase queries a database.
func (c *Client) QueryDatabase(ctx context.Context, databaseID string, filter map[string]any) ([]map[string]any, error) {
	resp, err := c.do(ctx, "POST", "/databases/"+databaseID+"/query", filter)
	if err != nil {
		return nil, err
	}

	var result struct {
		Results []map[string]any `json:"results"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

// UpdatePageProperties updates the properties of a page.
func (c *Client) UpdatePageProperties(ctx context.Context, pageID string, properties map[string]any) error {
	_, err := c.do(ctx, "PATCH", "/pages/"+pageID, map[string]any{
		"properties": properties,
	})
	return err
}

// GetPage retrieves a page.
func (c *Client) GetPage(ctx context.Context, pageID string) (map[string]any, error) {
	resp, err := c.do(ctx, "GET", "/pages/"+pageID, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}
