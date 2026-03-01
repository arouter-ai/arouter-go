package llmrouter

import (
	"context"
	"fmt"
	"net/http"
)

// --- OpenAPI-Gateway key management (uses API key auth) ---

// CreateSubKey creates a sub-key under the current API key.
// This calls the openapi-gateway directly — no JWT or dashboard login needed.
func (c *Client) CreateSubKey(ctx context.Context, req *CreateSubKeyRequest) (*CreateSubKeyResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/v1/keys/subkeys", req)
	if err != nil {
		return nil, err
	}

	var resp CreateSubKeyResponse
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListSubKeys lists sub-keys belonging to the current API key.
// This calls the openapi-gateway directly — no JWT or dashboard login needed.
func (c *Client) ListSubKeys(ctx context.Context, opts *ListKeysOptions) (*ListKeysResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodGet, "/v1/keys/subkeys", nil)
	if err != nil {
		return nil, err
	}

	if opts != nil {
		q := httpReq.URL.Query()
		if opts.PageSize > 0 {
			q.Set("page_size", fmt.Sprintf("%d", opts.PageSize))
		}
		if opts.PageToken != "" {
			q.Set("page_token", opts.PageToken)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	var resp ListKeysResponse
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RevokeSubKey revokes a sub-key by its ID. Only the parent key can revoke.
func (c *Client) RevokeSubKey(ctx context.Context, subKeyID string) error {
	httpReq, err := c.newRequest(ctx, http.MethodDelete, "/v1/keys/subkeys/"+subKeyID, nil)
	if err != nil {
		return err
	}
	return c.do(httpReq, nil)
}

// --- Dashboard-Gateway key management (uses JWT auth) ---

// CreateKey creates a new API key via the dashboard-gateway (requires JWT).
func (c *Client) CreateKey(ctx context.Context, req *CreateKeyRequest) (*CreateKeyResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/keys", req)
	if err != nil {
		return nil, err
	}

	var resp CreateKeyResponse
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListKeys returns a paginated list of API keys via the dashboard-gateway (requires JWT).
func (c *Client) ListKeys(ctx context.Context, opts *ListKeysOptions) (*ListKeysResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodGet, "/api/keys", nil)
	if err != nil {
		return nil, err
	}

	if opts != nil {
		q := httpReq.URL.Query()
		if opts.PageSize > 0 {
			q.Set("page_size", fmt.Sprintf("%d", opts.PageSize))
		}
		if opts.PageToken != "" {
			q.Set("page_token", opts.PageToken)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	var resp ListKeysResponse
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RevokeKey revokes an API key by its ID via the dashboard-gateway (requires JWT).
func (c *Client) RevokeKey(ctx context.Context, keyID string) error {
	httpReq, err := c.newRequest(ctx, http.MethodDelete, "/api/keys/"+keyID, nil)
	if err != nil {
		return err
	}
	return c.do(httpReq, nil)
}

// RotateKey rotates an API key via the dashboard-gateway (requires JWT).
func (c *Client) RotateKey(ctx context.Context, keyID string) (*CreateKeyResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/keys/"+keyID+"/rotate", nil)
	if err != nil {
		return nil, err
	}

	var resp CreateKeyResponse
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
