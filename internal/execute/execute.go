package execute

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"qforge/internal/model"
)

type rawExecuteResponse struct {
	Columns []string `json:"columns"`
	Rows    []any    `json:"rows"`
}

func ExecuteSQL(ctx context.Context, mcpURL, token, sql, logComment string) ([]byte, model.CanonicalResult, error) {
	endpoint := strings.TrimSuffix(mcpURL, "/http") + "/openapi/execute_query"
	finalQuery := withLogComment(sql, logComment)
	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, model.CanonicalResult{}, err
	}
	query := reqURL.Query()
	query.Set("query", finalQuery)
	reqURL.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, model.CanonicalResult{}, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Altinity-MCP-Key", token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, model.CanonicalResult{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, model.CanonicalResult{}, err
	}
	if resp.StatusCode >= 300 {
		return body, model.CanonicalResult{}, fmt.Errorf("execute query: %s", strings.TrimSpace(string(body)))
	}
	var raw rawExecuteResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return body, model.CanonicalResult{}, fmt.Errorf("parse execution response: %w", err)
	}
	rows, err := canonicalRows(raw.Columns, raw.Rows)
	if err != nil {
		return body, model.CanonicalResult{}, err
	}
	result := model.CanonicalResult{
		Columns:           raw.Columns,
		Rows:              rows,
		RowCount:          len(rows),
		GeneratedAt:       time.Now().UTC(),
		SourceQuerySHA256: sha256String(sql),
		LogComment:        logComment,
	}
	return body, result, nil
}

func CanonicalizeResult(rawBody []byte, sql, logComment string) (model.CanonicalResult, error) {
	var raw rawExecuteResponse
	if err := json.Unmarshal(rawBody, &raw); err != nil {
		return model.CanonicalResult{}, fmt.Errorf("parse execution response: %w", err)
	}
	rows, err := canonicalRows(raw.Columns, raw.Rows)
	if err != nil {
		return model.CanonicalResult{}, err
	}
	return model.CanonicalResult{
		Columns:           raw.Columns,
		Rows:              rows,
		RowCount:          len(rows),
		GeneratedAt:       time.Now().UTC(),
		SourceQuerySHA256: sha256String(sql),
		LogComment:        logComment,
	}, nil
}

func WriteJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func withLogComment(sql, logComment string) string {
	trimmed := strings.TrimSpace(sql)
	trimmed = strings.TrimSuffix(trimmed, ";")
	if logComment == "" {
		return trimmed
	}
	escaped := strings.ReplaceAll(logComment, "'", "\\'")
	return trimmed + " SETTINGS log_comment = '" + escaped + "'"
}

func canonicalRows(columns []string, rawRows []any) ([]map[string]any, error) {
	out := make([]map[string]any, 0, len(rawRows))
	for _, rawRow := range rawRows {
		values, ok := rawRow.([]any)
		if !ok {
			return nil, fmt.Errorf("unexpected row shape")
		}
		row := make(map[string]any, len(columns))
		for i, column := range columns {
			if i < len(values) {
				row[column] = values[i]
			} else {
				row[column] = nil
			}
		}
		out = append(out, row)
	}
	return out, nil
}

func sha256String(input string) string {
	sum := sha256Bytes([]byte(input))
	return fmt.Sprintf("%x", sum)
}

func sha256Bytes(input []byte) [32]byte {
	return sha256.Sum256(input)
}
