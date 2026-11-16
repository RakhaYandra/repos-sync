// Package model holds the JSON-visible data model (schema_version 1).
package model

import "sort"

type RepoInfo struct {
	Name              string  `json:"name"`
	Path              string  `json:"path"`
	Branch            string  `json:"branch"`
	Bare              bool    `json:"bare"`
	Upstream          *string `json:"upstream"`
	Ahead             *int    `json:"ahead"`
	Behind            *int    `json:"behind"`
	Dirty             int     `json:"dirty"`
	Untracked         int     `json:"untracked"`
	LastCommitAt      *string `json:"last_commit_at"`
	LastCommitAgeS    *int64  `json:"last_commit_age_s"`
	LastCommitSubject *string `json:"last_commit_subject"`
	Error             *string `json:"error"`
}

type Summary struct {
	Repos           int `json:"repos"`
	WithUpstream    int `json:"with_upstream"`
	AheadRepos      int `json:"ahead_repos"`
	BehindRepos     int `json:"behind_repos"`
	UnpushedCommits int `json:"unpushed_commits"`
	DirtyRepos      int `json:"dirty_repos"`
	UntrackedRepos  int `json:"untracked_repos"`
}

type Report struct {
	SchemaVersion int        `json:"schema_version"`
	GeneratedAt   string     `json:"generated_at"`
	Roots         []string   `json:"roots"`
	Summary       Summary    `json:"summary"`
	Repos         []RepoInfo `json:"repos"`
}

func SeverityBucket(r RepoInfo) int {
	if r.Error != nil && *r.Error != "" {
		return 0
	}
	if r.Dirty > 0 {
		return 0
	}
	if r.Ahead != nil && *r.Ahead > 0 {
		return 1
	}
	if r.Behind != nil && *r.Behind > 0 {
		return 2
	}
	if r.Upstream == nil {
		return 3
	}
	if r.Untracked > 0 {
		return 4
	}
	return 5
}

func SortRepos(repos []RepoInfo) {
	sort.SliceStable(repos, func(i, j int) bool {
		si, sj := SeverityBucket(repos[i]), SeverityBucket(repos[j])
		if si != sj {
			return si < sj
		}
		return repos[i].Name < repos[j].Name
	})
}

func Summarize(repos []RepoInfo) Summary {
	s := Summary{Repos: len(repos)}
	for _, r := range repos {
		if r.Upstream != nil {
			s.WithUpstream++
		}
		if r.Ahead != nil {
			if *r.Ahead > 0 {
				s.AheadRepos++
			}
			s.UnpushedCommits += *r.Ahead
		}
		if r.Behind != nil && *r.Behind > 0 {
			s.BehindRepos++
		}
		if r.Dirty > 0 && r.Error == nil {
			s.DirtyRepos++
		}
		if r.Untracked > 0 {
			s.UntrackedRepos++
		}
	}
	return s
}
