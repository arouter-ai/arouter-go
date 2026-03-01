# LLMRouter Go SDK

Official Go client for the [LLMRouter](https://github.com/llmrouter-ai) API gateway.

## Installation

```bash
go get github.com/llmrouter/llmrouter-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	llmrouter "github.com/llmrouter/llmrouter-go"
)

func main() {
	client := llmrouter.NewClient(
		"https://api.llmrouter.example.com",
		"your-api-key",
	)

	resp, err := client.ChatCompletion(context.Background(), &llmrouter.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []llmrouter.Message{
			{Role: "user", Content: "Hello!"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
```

## Features

- **Chat Completions** — OpenAI-compatible `/v1/chat/completions` with streaming support
- **Provider Proxy** — forward raw requests to any configured LLM provider
- **Sub-Key Management** — create, list, and revoke sub-keys programmatically via API key auth
- **Key Management** — full CRUD + rotation via dashboard gateway (JWT auth)
- **Usage Tracking** — query usage summaries and time-series data
- **Error Handling** — structured `APIError` with sentinel errors for `errors.Is` matching

## Streaming

```go
stream, err := client.ChatCompletionStream(ctx, &llmrouter.ChatCompletionRequest{
	Model:    "gpt-4o",
	Messages: []llmrouter.Message{{Role: "user", Content: "Tell me a story"}},
})
if err != nil {
	log.Fatal(err)
}
defer stream.Close()

for {
	chunk, err := stream.Recv()
	if err == llmrouter.ErrStreamDone {
		break
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(chunk.Choices[0].Delta.Content)
}
```

## Sub-Key Management

Sub-keys inherit permissions from a parent key and can be created directly via the API gateway — no dashboard login required.

```go
sub, err := client.CreateSubKey(ctx, &llmrouter.CreateSubKeyRequest{
	Name:          "worker-1",
	AllowedModels: []string{"gpt-4o-mini"},
	RateLimit: &llmrouter.RateLimitConfig{
		RequestsPerMinute: 60,
	},
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("Sub-key:", sub.RawKey)

// List all sub-keys
keys, _ := client.ListSubKeys(ctx, nil)
for _, k := range keys.Keys {
	fmt.Println(k.ID, k.Name)
}

// Revoke
_ = client.RevokeSubKey(ctx, sub.Key.ID)
```

## Provider Proxy

Forward raw requests to any provider endpoint (embeddings, images, etc.):

```go
body := strings.NewReader(`{"input": "hello", "model": "text-embedding-3-small"}`)
resp, err := client.ProxyRequest(ctx, "openai", "v1/embeddings", body)
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close()
// read resp.Body ...
```

## Usage Tracking

```go
summary, err := client.GetUsageSummary(ctx, &llmrouter.UsageQuery{
	StartTime: time.Now().AddDate(0, 0, -7),
	EndTime:   time.Now(),
})
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Requests: %d, Tokens: %d, Cost: $%.4f\n",
	summary.TotalRequests, summary.TotalTokens, summary.EstimatedCostUSD)
```

## Client Options

```go
client := llmrouter.NewClient(baseURL, apiKey,
	llmrouter.WithTimeout(60 * time.Second),
	llmrouter.WithHTTPClient(customHTTPClient),
)
```

## Error Handling

All API errors are returned as `*llmrouter.APIError` and can be matched with sentinel errors:

```go
_, err := client.ChatCompletion(ctx, req)
if errors.Is(err, llmrouter.ErrRateLimited) {
	// back off and retry
}
if errors.Is(err, llmrouter.ErrQuotaExceeded) {
	// quota exhausted
}

var apiErr *llmrouter.APIError
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.StatusCode, apiErr.Code, apiErr.Message)
}
```

Available sentinel errors: `ErrUnauthorized`, `ErrForbidden`, `ErrNotFound`, `ErrRateLimited`, `ErrQuotaExceeded`, `ErrBadRequest`, `ErrServerError`.

## License

MIT
