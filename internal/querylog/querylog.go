package querylog

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"qforge/internal/model"
)

const demoConnectionName = "demo"

func FetchLatest(ctx context.Context, logComment string) (*model.QueryLogMetrics, error) {
	escaped := strings.ReplaceAll(logComment, "'", "''")
	sql := fmt.Sprintf(`
SELECT
  log_comment,
  query_id,
  query_duration_ms,
  read_rows,
  read_bytes,
  result_rows,
  result_bytes,
  memory_usage,
  peak_threads_usage AS peak_threads,
  query,
  toString(event_time) AS event_time,
  type
FROM system.query_log
WHERE log_comment = '%s'
ORDER BY event_time_microseconds DESC
LIMIT 1
FORMAT JSONEachRow
`, escaped)

	cmd := exec.CommandContext(ctx, "clickhouse-client", "--connection", demoConnectionName, "--query", sql)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("clickhouse-client --connection %s: %s", demoConnectionName, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, err
	}
	if strings.TrimSpace(string(output)) == "" {
		return nil, nil
	}

	var row struct {
		LogComment      string `json:"log_comment"`
		QueryID         string `json:"query_id"`
		QueryDurationMS int64  `json:"query_duration_ms"`
		ReadRows        int64  `json:"read_rows"`
		ReadBytes       int64  `json:"read_bytes"`
		ResultRows      int64  `json:"result_rows"`
		ResultBytes     int64  `json:"result_bytes"`
		MemoryUsage     int64  `json:"memory_usage"`
		PeakThreads     int64  `json:"peak_threads"`
		Query           string `json:"query"`
		EventTime       string `json:"event_time"`
		Type            string `json:"type"`
	}
	if err := json.Unmarshal(output, &row); err != nil {
		return nil, fmt.Errorf("parse query_log row: %w", err)
	}
	return &model.QueryLogMetrics{
		LogComment:      row.LogComment,
		QueryID:         row.QueryID,
		QueryDurationMS: row.QueryDurationMS,
		ReadRows:        row.ReadRows,
		ReadBytes:       row.ReadBytes,
		ResultRows:      row.ResultRows,
		ResultBytes:     row.ResultBytes,
		MemoryUsage:     row.MemoryUsage,
		PeakThreads:     row.PeakThreads,
		Query:           row.Query,
		EventTime:       row.EventTime,
		Type:            row.Type,
	}, nil
}
