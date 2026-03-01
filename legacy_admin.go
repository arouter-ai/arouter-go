package llmrouter

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// LegacyAdminClient manages per-instance router tokens via the
// LLMRouter's /admin/keys backward-compatibility API. This replaces
// the RemoteKeyStore previously embedded in BeeOS.
type LegacyAdminClient struct {
	baseURL    string
	adminToken string
	httpClient *http.Client
}

// NewLegacyAdminClient creates a client for the LLMRouter legacy admin API.
// baseURL is the openapi-gateway URL (e.g. "http://localhost:19080").
// adminToken is the ADMIN_TOKEN configured on the LLMRouter (can be empty for
// backward-compatible deployments that don't require auth).
func NewLegacyAdminClient(baseURL, adminToken string) *LegacyAdminClient {
	return &LegacyAdminClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		adminToken: adminToken,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type storeKeysRequest struct {
	Token      string            `json:"token"`
	InstanceID string            `json:"instanceId"`
	OwnerID    string            `json:"ownerId"`
	APIKeys    map[string]string `json:"apiKeys"`
}

// StoreKeys registers a per-instance token → API keys mapping.
func (c *LegacyAdminClient) StoreKeys(ctx context.Context, token, ownerID string, apiKeys map[string]string) error {
	body := storeKeysRequest{
		Token:      token,
		InstanceID: ExtractInstanceID(token),
		OwnerID:    ownerID,
		APIKeys:    apiKeys,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal store request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/admin/keys", bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.adminToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.adminToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("POST /admin/keys: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("POST /admin/keys: status %d: %s", resp.StatusCode, errBody)
	}
	return nil
}

// DeleteByToken removes a specific token mapping.
func (c *LegacyAdminClient) DeleteByToken(ctx context.Context, token string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/admin/keys/"+token, nil)
	if err != nil {
		return err
	}
	if c.adminToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.adminToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE /admin/keys/{token}: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("DELETE /admin/keys/{token}: status %d: %s", resp.StatusCode, errBody)
	}
	return nil
}

// DeleteByInstance removes all tokens for a given instance ID.
func (c *LegacyAdminClient) DeleteByInstance(ctx context.Context, instanceID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/admin/keys/instance/"+instanceID, nil)
	if err != nil {
		return err
	}
	if c.adminToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.adminToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE /admin/keys/instance/{id}: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("DELETE /admin/keys/instance/{id}: status %d: %s", resp.StatusCode, errBody)
	}
	return nil
}

// GenerateRouterToken creates a cryptographically random internal token
// scoped to an instance. Format: "rt_{instanceID}_{random}".
func GenerateRouterToken(instanceID string) string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("rt_%s_fallback", instanceID)
	}
	return fmt.Sprintf("rt_%s_%s", instanceID, hex.EncodeToString(b[:8]))
}

// ExtractInstanceID extracts the instance ID from a router token.
func ExtractInstanceID(token string) string {
	var rest string
	switch {
	case strings.HasPrefix(token, "rt_"):
		rest = token[3:]
	case strings.HasPrefix(token, "inst_"):
		rest = token[5:]
	default:
		return token
	}
	if idx := strings.LastIndex(rest, "_"); idx > 0 {
		return rest[:idx]
	}
	return rest
}
