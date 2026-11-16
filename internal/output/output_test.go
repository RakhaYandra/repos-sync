package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/RakhaYandra/repos-sync/internal/model"
)

func mkReport() model.Report {
	strPtr := func(s string) *string { return &s }
	intPtr := func(n int) *int { return &n }
	r := model.Report{
		SchemaVersion: 1,
		GeneratedAt:   "2026-11-16T00:00:00Z",
		Roots:         []string{"/x"},
		Repos: []model.RepoInfo{
			{Name: "a", Path: "/x/a", Branch: "main", Upstream: strPtr("origin/main"), Ahead: intPtr(0), Behind: intPtr(0)},
			{Name: "b", Path: "/x/b", Branch: "main", Upstream: strPtr("origin/main"), Ahead: intPtr(2), Behind: intPtr(0), Dirty: 1, LastCommitSubject: strPtr("subject line")},
		},
	}
	model.SortRepos(r.Repos)
	r.Summary = model.Summarize(r.Repos)
	return r
}

func TestJSONSchemaVersionFirstKeyAndShape(t *testing.T) {
	var buf bytes.Buffer
	if err := JSON(&buf, mkReport()); err != nil {
		t.Fatal(err)
	}
	first := strings.SplitN(buf.String(), "\n", 3)[1]
	if !strings.Contains(first, `"schema_version": 1`) {
		t.Fatalf("schema_version not first: %q", first)
	}
	var v struct {
		SchemaVersion int    `json:"schema_version"`
		GeneratedAt   string `json:"generated_at"`
		Summary       struct {
			Repos           int `json:"repos"`
			UnpushedCommits int `json:"unpushed_commits"`
		} `json:"summary"`
		Repos []struct {
			Name   string  `json:"name"`
			Behind *int    `json:"behind"`
			Dirty  int     `json:"dirty"`
			Error  *string `json:"error"`
		} `json:"repos"`
	}
	if err := json.Unmarshal(buf.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	if v.SchemaVersion != 1 || v.Summary.Repos != 2 || v.Summary.UnpushedCommits != 2 {
		t.Fatalf("got %+v", v)
	}
	if v.Repos[0].Name != "b" {
		t.Fatalf("severity sort: first=%s", v.Repos[0].Name)
	}
}

func TestHumanContainsCounts(t *testing.T) {
	var buf bytes.Buffer
	Human(&buf, mkReport(), false)
	s := buf.String()
	for _, want := range []string{"b", "ahead 2", "dirty 1", "2 repositories"} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %q in:\n%s", want, s)
		}
	}
}
