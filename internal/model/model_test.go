package model

import "testing"

func TestSeverityOrder(t *testing.T) {
	strPtr := func(s string) *string { return &s }
	intPtr := func(n int) *int { return &n }
	repos := []RepoInfo{
		{Name: "clean", Branch: "main", Upstream: strPtr("origin/main"), Ahead: intPtr(0), Behind: intPtr(0)},
		{Name: "untracked", Upstream: strPtr("origin/main"), Ahead: intPtr(0), Behind: intPtr(0), Untracked: 1},
		{Name: "behind", Upstream: strPtr("origin/main"), Ahead: intPtr(0), Behind: intPtr(1)},
		{Name: "ahead", Upstream: strPtr("origin/main"), Ahead: intPtr(1), Behind: intPtr(0)},
		{Name: "dirty", Upstream: strPtr("origin/main"), Ahead: intPtr(0), Behind: intPtr(0), Dirty: 2},
	}
	SortRepos(repos)
	want := []string{"dirty", "ahead", "behind", "untracked", "clean"}
	for i, name := range want {
		if repos[i].Name != name {
			t.Fatalf("order=%v", func() []string {
				out := []string{}
				for _, r := range repos {
					out = append(out, r.Name)
				}
				return out
			}())
		}
	}
}

func TestSummarize(t *testing.T) {
	intPtr := func(n int) *int { return &n }
	strPtr := func(s string) *string { return &s }
	repos := []RepoInfo{
		{Name: "a", Upstream: strPtr("origin/main"), Ahead: intPtr(2), Behind: intPtr(0)},
		{Name: "b", Ahead: nil, Behind: nil},
		{Name: "c", Upstream: strPtr("origin/main"), Ahead: intPtr(0), Behind: intPtr(3), Dirty: 1},
	}
	s := Summarize(repos)
	if s.Repos != 3 || s.WithUpstream != 2 || s.UnpushedCommits != 2 || s.AheadRepos != 1 || s.BehindRepos != 1 || s.DirtyRepos != 1 {
		t.Fatalf("%+v", s)
	}
}
