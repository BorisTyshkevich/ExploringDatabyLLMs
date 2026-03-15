package querylog

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"qforge/internal/execute"
	"qforge/internal/model"
)

func FetchLatest(ctx context.Context, mcpURL, token, logComment string) (*model.QueryLogMetrics, error) {
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
`, escaped)
	_, result, err := execute.ExecuteSQL(ctx, mcpURL, token, sql, "")
	if err != nil {
		return nil, err
	}
	if len(result.Rows) == 0 {
		return nil, nil
	}
	row := result.Rows[0]
	metrics := &model.QueryLogMetrics{
		LogComment:      stringValue(row["log_comment"]),
		QueryID:         stringValue(row["query_id"]),
		QueryDurationMS: int64Value(row["query_duration_ms"]),
		ReadRows:        int64Value(row["read_rows"]),
		ReadBytes:       int64Value(row["read_bytes"]),
		ResultRows:      int64Value(row["result_rows"]),
		ResultBytes:     int64Value(row["result_bytes"]),
		MemoryUsage:     int64Value(row["memory_usage"]),
		PeakThreads:     int64Value(row["peak_threads"]),
		Query:           stringValue(row["query"]),
		EventTime:       stringValue(row["event_time"]),
		Type:            stringValue(row["type"]),
	}
	return metrics, nil
}

func stringValue(value any) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}

func int64Value(value any) int64 {
	switch v := value.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case json.Number:
		n, _ := v.Int64()
		return n
	case string:
		n, _ := strconv.ParseInt(v, 10, 64)
		return n
	default:
		n, _ := strconv.ParseInt(fmt.Sprint(v), 10, 64)
		return n
	}
}
