package llmrouter

import (
	"context"
	"net/http"
)

// GetUsageSummary retrieves an aggregated usage summary.
func (c *Client) GetUsageSummary(ctx context.Context, q *UsageQuery) (*UsageSummaryResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodGet, "/api/usage/summary", nil)
	if err != nil {
		return nil, err
	}
	applyUsageQuery(httpReq, q)

	var resp UsageSummaryResponse
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUsageTimeSeries retrieves time-bucketed usage data.
func (c *Client) GetUsageTimeSeries(ctx context.Context, q *UsageQuery) (*UsageTimeSeriesResponse, error) {
	httpReq, err := c.newRequest(ctx, http.MethodGet, "/api/usage/timeseries", nil)
	if err != nil {
		return nil, err
	}
	applyUsageQuery(httpReq, q)

	var resp UsageTimeSeriesResponse
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func applyUsageQuery(req *http.Request, q *UsageQuery) {
	if q == nil {
		return
	}
	params := req.URL.Query()
	if !q.StartTime.IsZero() {
		params.Set("start_time", q.StartTime.Format("2006-01-02T15:04:05Z"))
	}
	if !q.EndTime.IsZero() {
		params.Set("end_time", q.EndTime.Format("2006-01-02T15:04:05Z"))
	}
	if q.ProviderID != "" {
		params.Set("provider_id", q.ProviderID)
	}
	if q.Model != "" {
		params.Set("model", q.Model)
	}
	if q.KeyID != "" {
		params.Set("key_id", q.KeyID)
	}
	if q.Granularity != "" {
		params.Set("granularity", q.Granularity)
	}
	req.URL.RawQuery = params.Encode()
}
