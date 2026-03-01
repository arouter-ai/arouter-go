package llmrouter

import (
	"encoding/json"
	"time"
)

// Message represents a chat message in the OpenAI-compatible format.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// ChatCompletionRequest is the request payload for chat completions.
type ChatCompletionRequest struct {
	Model            string         `json:"model"`
	Messages         []Message      `json:"messages"`
	Temperature      *float64       `json:"temperature,omitempty"`
	TopP             *float64       `json:"top_p,omitempty"`
	N                *int           `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	MaxTokens        *int           `json:"max_tokens,omitempty"`
	PresencePenalty  *float64       `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64       `json:"frequency_penalty,omitempty"`
	User             string         `json:"user,omitempty"`
	Extra            map[string]any `json:"-"`
}

// ChatCompletionResponse is the response from a non-streaming chat completion.
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

// Choice represents a single completion choice.
type Choice struct {
	Index        int      `json:"index"`
	Message      *Message `json:"message,omitempty"`
	Delta        *Message `json:"delta,omitempty"`
	FinishReason *string  `json:"finish_reason,omitempty"`
}

// Usage tracks token consumption.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionChunk is a single chunk from a streaming chat completion.
type ChatCompletionChunk struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

// --- Key Management Types ---

// RateLimitConfig configures per-key rate limits.
type RateLimitConfig struct {
	RequestsPerMinute int32 `json:"requests_per_minute,omitempty"`
	RequestsPerDay    int32 `json:"requests_per_day,omitempty"`
	TokensPerMinute   int32 `json:"tokens_per_minute,omitempty"`
}

// QuotaLimitConfig configures per-key quota limits.
type QuotaLimitConfig struct {
	MaxTokensPerDay   int64 `json:"max_tokens_per_day,omitempty"`
	MaxTokensPerMonth int64 `json:"max_tokens_per_month,omitempty"`
	MaxRequestsPerDay int64 `json:"max_requests_per_day,omitempty"`
}

// CreateKeyRequest is the payload for creating a new API key.
type CreateKeyRequest struct {
	Name             string            `json:"name"`
	AllowedProviders []string          `json:"allowed_providers,omitempty"`
	AllowedModels    []string          `json:"allowed_models,omitempty"`
	RateLimit        *RateLimitConfig  `json:"rate_limit,omitempty"`
	QuotaLimit       *QuotaLimitConfig `json:"quota_limit,omitempty"`
	ExpiresAt        *string           `json:"expires_at,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// APIKeyInfo represents the key metadata returned by the service.
type APIKeyInfo struct {
	ID               string            `json:"id"`
	TenantID         string            `json:"tenant_id,omitempty"`
	Prefix           string            `json:"prefix,omitempty"`
	Name             string            `json:"name"`
	Status           int               `json:"status,omitempty"`
	AllowedProviders []string          `json:"allowed_providers,omitempty"`
	AllowedModels    []string          `json:"allowed_models,omitempty"`
	RateLimit        *RateLimitConfig  `json:"rate_limit,omitempty"`
	QuotaLimit       *QuotaLimitConfig `json:"quota_limit,omitempty"`
	ExpiresAt        *string           `json:"expires_at,omitempty"`
	ParentKeyID      string            `json:"parent_key_id,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	CreatedAt        json.RawMessage   `json:"created_at,omitempty"`
}

// CreateKeyResponse is returned when a key is created (includes the raw secret).
type CreateKeyResponse struct {
	Key    APIKeyInfo `json:"key"`
	RawKey string     `json:"raw_key"`
}

// ListKeysOptions contains query parameters for listing keys.
type ListKeysOptions struct {
	PageSize  int    `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}

// ListKeysResponse is the paginated response for key listing.
type ListKeysResponse struct {
	Keys          []APIKeyInfo `json:"keys"`
	NextPageToken string       `json:"next_page_token,omitempty"`
}

// CreateSubKeyRequest is the payload for creating a sub-key.
type CreateSubKeyRequest struct {
	Name             string            `json:"name"`
	AllowedProviders []string          `json:"allowed_providers,omitempty"`
	AllowedModels    []string          `json:"allowed_models,omitempty"`
	RateLimit        *RateLimitConfig  `json:"rate_limit,omitempty"`
	QuotaLimit       *QuotaLimitConfig `json:"quota_limit,omitempty"`
	ExpiresAt        *string           `json:"expires_at,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// CreateSubKeyResponse is returned when a sub-key is created.
type CreateSubKeyResponse struct {
	Key    APIKeyInfo `json:"key"`
	RawKey string     `json:"raw_key"`
}

// --- Usage Types ---

// UsageQuery contains filters for usage queries.
type UsageQuery struct {
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	ProviderID  string    `json:"provider_id,omitempty"`
	Model       string    `json:"model,omitempty"`
	KeyID       string    `json:"key_id,omitempty"`
	Granularity string    `json:"granularity,omitempty"`
}

// UsageSummaryResponse matches the gRPC UsageSummary response.
type UsageSummaryResponse struct {
	TotalRequests     int64           `json:"total_requests"`
	TotalInputTokens  int64           `json:"total_input_tokens"`
	TotalOutputTokens int64           `json:"total_output_tokens"`
	TotalTokens       int64           `json:"total_tokens"`
	EstimatedCostUSD  float64         `json:"estimated_cost_usd"`
	ByProvider        []ProviderUsage `json:"by_provider,omitempty"`
}

// ProviderUsage is per-provider usage within a summary.
type ProviderUsage struct {
	ProviderID       string  `json:"provider_id"`
	Requests         int64   `json:"requests"`
	InputTokens      int64   `json:"input_tokens"`
	OutputTokens     int64   `json:"output_tokens"`
	EstimatedCostUSD float64 `json:"estimated_cost_usd"`
}

// UsageTimeSeriesResponse contains time-bucketed usage data.
type UsageTimeSeriesResponse struct {
	DataPoints []UsageDataPoint `json:"data_points"`
}

// UsageDataPoint is one data point in a usage time series.
type UsageDataPoint struct {
	Timestamp        string  `json:"timestamp"`
	Requests         int64   `json:"requests"`
	InputTokens      int64   `json:"input_tokens"`
	OutputTokens     int64   `json:"output_tokens"`
	EstimatedCostUSD float64 `json:"estimated_cost_usd"`
}
