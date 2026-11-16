// Package scan walks roots and probes git repositories.
package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RakhaYandra/repos-sync/internal/git"
	"github.com/RakhaYandra/repos-sync/internal/model"
)

func ScanRoots(roots []string) (model.Report, error) {
	if !git.GitAvailable() {
		return model.Report{}, fmt.Errorf("git not found on PATH")
	}
	report := model.Report{
		SchemaVersion: 1,
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Roots:         roots,
		Repos:         []model.RepoInfo{},
	}
	for _, root := range roots {
		fi, err := os.Stat(root)
		if err != nil || !fi.IsDir() {
			return report, fmt.Errorf("root %s: not a directory", root)
		}
		entries, err := os.ReadDir(root)
		if err != nil {
			return report, fmt.Errorf("root %s: %w", root, err)
		}
		// A bare repository (e.g. ~/.dotfiles) is scanned as one entry.
		if git.IsRepoRoot(root) {
			if out, e := git.IsBare(root); e == nil && out {
				report.Repos = append(report.Repos, probeOne(root))
				continue
			}
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			name := entry.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}
			dir := filepath.Join(root, name)
			if !git.IsRepoRoot(dir) {
				continue
			}
			report.Repos = append(report.Repos, probeOne(dir))
		}
	}
	model.SortRepos(report.Repos)
	report.Summary = model.Summarize(report.Repos)
	return report, nil
}

func probeOne(dir string) model.RepoInfo {
	abs, err := filepath.Abs(dir)
	if err == nil {
		dir = abs
	}
	r := model.RepoInfo{Name: filepath.Base(dir), Path: dir}
	p := git.ProbeRepo(dir)

	if p.Err != "" {
		r.Error = &p.Err
		r.Branch = p.Branch
		r.Bare = p.Bare
		return r
	}
	r.Branch = p.Branch
	r.Bare = p.Bare
	if p.HasUpstream {
		r.Upstream = &p.Upstream
		r.Ahead = p.Ahead
		r.Behind = p.Behind
	}
	r.Dirty = p.Dirty
	r.Untracked = p.Untracked
	if p.LastAt != "" {
		at := p.LastAt
		r.LastCommitAt = &at
	}
	if p.LastEpoch != nil {
		age := time.Now().UTC().Unix() - *p.LastEpoch
		r.LastCommitAgeS = &age
	}
	if p.LastSubject != "" {
		s := p.LastSubject
		r.LastCommitSubject = &s
	}
	return r
}

func CheckRepo(dir string) (model.Report, error) {
	if !git.GitAvailable() {
		return model.Report{}, fmt.Errorf("git not found on PATH")
	}
	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		return model.Report{}, fmt.Errorf("%s: not a directory", dir)
	}
	if !git.IsRepoRoot(dir) {
		return model.Report{}, fmt.Errorf("%s: not a git repository", dir)
	}
	abs, err := filepath.Abs(dir)
	if err == nil {
		dir = abs
	}
	r := probeOne(dir)
	report := model.Report{
		SchemaVersion: 1,
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Roots:         []string{dir},
		Repos:         []model.RepoInfo{r},
	}
	report.Summary = model.Summarize(report.Repos)
	return report, nil
}
