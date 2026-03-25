package querylog

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"qforge/internal/model"
)

const demoConnectionName = "demo"

func FetchLatest(ctx context.Context, logComment string) (*model.QueryLogMetrics, error) {
	rows, err := fetchRows(ctx, fmt.Sprintf("WHERE log_comment = '%s'", escapeLiteral(logComment)))
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

func FetchForRun(ctx context.Context, logComment string, multiQuery bool) (*model.QueryLogMetrics, error) {
	whereClause := fmt.Sprintf("WHERE log_comment = '%s'", escapeLiteral(logComment))
	if multiQuery {
		whereClause = fmt.Sprintf(
			"WHERE log_comment = '%s' OR startsWith(log_comment, '%s|section=')",
			escapeLiteral(logComment),
			escapeLiteral(logComment),
		)
	}
	rows, err := fetchRows(ctx, whereClause)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	if !multiQuery {
		return &rows[0], nil
	}
	return aggregateRows(rows), nil
}

func fetchRows(ctx context.Context, whereClause string) ([]model.QueryLogMetrics, error) {
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
%s
ORDER BY event_time_microseconds DESC
FORMAT JSONEachRow
`, whereClause)

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

	var rows []model.QueryLogMetrics
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var row model.QueryLogMetrics
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			return nil, fmt.Errorf("parse query_log row: %w", err)
		}
		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan query_log rows: %w", err)
	}
	return rows, nil
}

func aggregateRows(rows []model.QueryLogMetrics) *model.QueryLogMetrics {
	agg := rows[0]
	for _, row := range rows[1:] {
		agg.QueryDurationMS += row.QueryDurationMS
		agg.ReadRows += row.ReadRows
		agg.ReadBytes += row.ReadBytes
		agg.ResultRows += row.ResultRows
		agg.ResultBytes += row.ResultBytes
		if row.MemoryUsage > agg.MemoryUsage {
			agg.MemoryUsage = row.MemoryUsage
		}
		if row.PeakThreads > agg.PeakThreads {
			agg.PeakThreads = row.PeakThreads
		}
	}
	agg.QueryID = ""
	agg.Query = ""
	agg.Type = "aggregated"
	return &agg
}

func escapeLiteral(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
