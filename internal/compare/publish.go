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
	runsRepoPagesBase      = "https://boristyshkevich.github.io/ExploringDatabyLLMs-runs"
	qforgeRepoGitHubBlob   = "https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main"
)

type ArtifactRef struct {
	LocalPath     string `json:"local_path,omitempty"`
	PublishedPath string `json:"published_path,omitempty"`
	URL           string `json:"url,omitempty"`
}

type ArtifactLinks struct {
	QuerySQL   ArtifactRef `json:"query_sql,omitempty"`
	ReportMD   ArtifactRef `json:"report_md,omitempty"`
	ResultJSON ArtifactRef `json:"result_json,omitempty"`
	VisualHTML ArtifactRef `json:"visual_html,omitempty"`
}

func buildRunArtifactLinks(runsRoot, runDir string) ArtifactLinks {
	return ArtifactLinks{
		QuerySQL:   buildArtifactRef(runsRoot, filepath.Join(runDir, "query.sql"), "sql"),
		ReportMD:   buildArtifactRef(runsRoot, filepath.Join(runDir, "report.md"), "md"),
		ResultJSON: buildArtifactRef(runsRoot, filepath.Join(runDir, "result.json"), "json"),
		VisualHTML: buildArtifactRef(runsRoot, filepath.Join(runDir, "visual.html"), "html"),
	}
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

func qforgeRepoURL(repoRoot, path string) string {
	rel := repoRelativePath(repoRoot, path)
	if rel == "" || strings.HasPrefix(rel, "/") || strings.HasPrefix(rel, "..") {
		return ""
	}
	return runsRepoJoin(qforgeRepoGitHubBlob, rel)
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

func runsRepoJoin(base, rel string) string {
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(rel, "/")
}
