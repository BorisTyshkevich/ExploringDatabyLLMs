package querylog

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFetchForRunAggregatesMultiQueryRows(t *testing.T) {
	binDir := t.TempDir()
	scriptPath := filepath.Join(binDir, "clickhouse-client")
	script := "#!/bin/sh\nprintf '%s' \"$CLICKHOUSE_CLIENT_OUTPUT\"\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake clickhouse-client: %v", err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("CLICKHOUSE_CLIENT_OUTPUT", strings.Join([]string{
		`{"log_comment":"base|subquestion=q2","query_id":"2","query_duration_ms":200,"read_rows":20,"read_bytes":2000,"result_rows":2,"result_bytes":200,"memory_usage":512,"peak_threads":4,"query":"SELECT 2","event_time":"2026-03-24 12:00:02","type":"QueryFinish"}`,
		`{"log_comment":"base|subquestion=q1","query_id":"1","query_duration_ms":100,"read_rows":10,"read_bytes":1000,"result_rows":1,"result_bytes":100,"memory_usage":1024,"peak_threads":2,"query":"SELECT 1","event_time":"2026-03-24 12:00:01","type":"QueryFinish"}`,
	}, "\n"))

	got, err := FetchForRun(context.Background(), "base", true)
	if err != nil {
		t.Fatalf("FetchForRun returned error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected aggregated metrics, got nil")
	}
	if got.QueryDurationMS != 300 || got.ReadRows != 30 || got.ReadBytes != 3000 || got.ResultRows != 3 || got.ResultBytes != 300 {
		t.Fatalf("unexpected aggregated totals: %+v", got)
	}
	if got.MemoryUsage != 1024 || got.PeakThreads != 4 {
		t.Fatalf("unexpected aggregated maxima: %+v", got)
	}
	if got.Type != "aggregated" {
		t.Fatalf("expected aggregated type, got %+v", got)
	}
}

func TestFetchForRunReturnsNilWhenNoRowsFound(t *testing.T) {
	binDir := t.TempDir()
	scriptPath := filepath.Join(binDir, "clickhouse-client")
	script := "#!/bin/sh\nprintf '%s' \"$CLICKHOUSE_CLIENT_OUTPUT\"\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake clickhouse-client: %v", err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("CLICKHOUSE_CLIENT_OUTPUT", "")

	got, err := FetchForRun(context.Background(), "base", true)
	if err != nil {
		t.Fatalf("FetchForRun returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil metrics, got %+v", got)
	}
}
