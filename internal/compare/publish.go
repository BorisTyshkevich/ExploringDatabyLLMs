package compare

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	runsRepoGitHubBlobBase = "https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main"
	runsRepoGitHubTreeBase = "https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/tree/main"
	runsRepoPagesBase      = "https://boristyshkevich.github.io/ExploringDatabyLLMs-runs"
)

type ArtifactRef struct {
	LocalPath     string `json:"local_path,omitempty"`
	PublishedPath string `json:"published_path,omitempty"`
	URL           string `json:"url,omitempty"`
}

type ArtifactLinks struct {
	PromptReportMD ArtifactRef   `json:"prompt_report_md,omitempty"`
	PromptVisualMD ArtifactRef   `json:"prompt_visual_md,omitempty"`
	QuerySQL       ArtifactRef   `json:"query_sql,omitempty"`
	QuerySQLs      []ArtifactRef `json:"query_sqls,omitempty"`
	ReportMD       ArtifactRef   `json:"report_md,omitempty"`
	ReviewMD       ArtifactRef   `json:"review_md,omitempty"`
	ResultJSON     ArtifactRef   `json:"result_json,omitempty"`
	VisualHTML     ArtifactRef   `json:"visual_html,omitempty"`
	VisualSource   ArtifactRef   `json:"visual_source,omitempty"`
	VisualBuild    ArtifactRef   `json:"visual_build,omitempty"`
}

func buildRunArtifactLinks(runsRoot, runDir string) ArtifactLinks {
	querySQLs := buildQueryArtifactRefs(runsRoot, runDir)
	querySQL := ArtifactRef{}
	if len(querySQLs) == 1 {
		querySQL = querySQLs[0]
	}
	return ArtifactLinks{
		PromptReportMD: buildArtifactRef(runsRoot, filepath.Join(runDir, "prompt.report.md"), "md"),
		PromptVisualMD: buildArtifactRef(runsRoot, filepath.Join(runDir, "prompt.visual.md"), "md"),
		QuerySQL:       querySQL,
		QuerySQLs:      querySQLs,
		ReportMD:       buildArtifactRef(runsRoot, filepath.Join(runDir, "report.md"), "md"),
		ReviewMD:       buildArtifactRef(runsRoot, filepath.Join(runDir, "review.md"), "md"),
		ResultJSON:     buildArtifactRef(runsRoot, filepath.Join(runDir, "result.json"), "json"),
		VisualHTML:     buildArtifactRef(runsRoot, filepath.Join(runDir, "visual.html"), "html"),
		VisualSource:   buildArtifactRef(runsRoot, filepath.Join(runDir, "visual_src"), "dir"),
		VisualBuild:    buildArtifactRef(runsRoot, filepath.Join(runDir, "visual_build"), "dir"),
	}
}

func buildQueryArtifactRefs(runsRoot, runDir string) []ArtifactRef {
	single := buildArtifactRef(runsRoot, filepath.Join(runDir, "query.sql"), "sql")
	if single.URL != "" {
		return []ArtifactRef{single}
	}
	entries, err := os.ReadDir(filepath.Join(runDir, "queries"))
	if err != nil {
		return nil
	}
	var refs []ArtifactRef
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		ref := buildArtifactRef(runsRoot, filepath.Join(runDir, "queries", entry.Name()), "sql")
		if ref.URL != "" {
			refs = append(refs, ref)
		}
	}
	return refs
}

func buildArtifactRef(runsRoot, localPath, kind string) ArtifactRef {
	if _, err := os.Stat(localPath); err != nil {
		return ArtifactRef{}
	}
	ref := ArtifactRef{
		LocalPath:     repoRelativePath(runsRoot, localPath),
		PublishedPath: publishedRelativePath(runsRoot, localPath),
	}
	switch kind {
	case "md":
		ref.URL = publishedMarkdownURL(ref.PublishedPath)
	case "html":
		ref.URL = publishedVisualURL(ref.PublishedPath)
	case "dir":
		ref.URL = publishedTreeURL(ref.PublishedPath)
	default:
		ref.URL = publishedBlobURL(ref.PublishedPath)
	}
	return ref
}

func repoRelativePath(repoRoot, path string) string {
	if path == "" {
		return path
	}
	rel, err := filepath.Rel(repoRoot, path)
	if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return filepath.ToSlash(rel)
	}
	return path
}

func publishedRelativePath(repoRoot, path string) string {
	if path == "" {
		return path
	}
	rel, err := filepath.Rel(repoRoot, path)
	if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return filepath.ToSlash(rel)
	}
	return repoRelativePath(repoRoot, path)
}

func publishedMarkdownURL(publishedPath string) string {
	if publishedPath == "" {
		return ""
	}
	return fmt.Sprintf("%s/md.html?file=%s", runsRepoPagesBase, url.QueryEscape(publishedPath))
}

func publishedVisualURL(publishedPath string) string {
	if publishedPath == "" {
		return ""
	}
	return runsRepoJoin(runsRepoPagesBase, publishedPath)
}

func publishedBlobURL(publishedPath string) string {
	if publishedPath == "" {
		return ""
	}
	return runsRepoJoin(runsRepoGitHubBlobBase, publishedPath)
}

func publishedTreeURL(publishedPath string) string {
	if publishedPath == "" {
		return ""
	}
	return runsRepoJoin(runsRepoGitHubTreeBase, publishedPath)
}

func runsRepoJoin(base, rel string) string {
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(rel, "/")
}
