// Package output renders the scan report as human text or JSON.
package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/RakhaYandra/repos-sync/internal/model"
)

func Human(w io.Writer, r model.Report, color bool) {
	red, reset, grey := "", "", ""
	yellow := ""
	if color {
		red, yellow, grey, reset = "\033[31m", "\033[33m", "\033[90m", "\033[0m"
	}
	fmt.Fprintf(w, "%-22s %-12s %-4s %-4s %s\n", "repo", "branch", "↑", "↓", "status")
	for _, repo := range r.Repos {
		up := 0
		down := 0
		if repo.Ahead != nil {
			up = *repo.Ahead
		}
		if repo.Behind != nil {
			down = *repo.Behind
		}
		upStr := fmt.Sprintf("%d", up)
		downStr := fmt.Sprintf("%d", down)
		isErr := repo.Error != nil && *repo.Error != ""
		colorOf := func() string {
			if repo.Dirty > 0 || up > 0 {
				return red
			}
			if down > 0 {
				return yellow
			}
			return grey
		}
		stateParts := []string{}
		switch {
		case isErr:
			stateParts = append(stateParts, "error: "+*repo.Error)
			upStr, downStr = "-", "-"
		case repo.Upstream == nil:
			stateParts = append(stateParts, "no upstream")
			upStr, downStr = "-", "-"
		default:
			if repo.Dirty > 0 {
				stateParts = append(stateParts, fmt.Sprintf("dirty %d", repo.Dirty))
			}
			if up > 0 {
				stateParts = append(stateParts, fmt.Sprintf("ahead %d", up))
			}
			if down > 0 {
				stateParts = append(stateParts, fmt.Sprintf("behind %d", down))
			}
			if repo.Untracked > 0 {
				stateParts = append(stateParts, fmt.Sprintf("untracked %d", repo.Untracked))
			}
			if len(stateParts) == 0 {
				stateParts = append(stateParts, "clean")
			}
		}
		state := colorOf() + strings.Join(stateParts, ", ") + reset
		branch := repo.Branch
		if branch == "" {
			branch = "-"
		}
		fmt.Fprintf(w, "%-22s %-12s %-4s %-4s %s\n", repo.Name, branch, upStr, downStr, state)
	}
	s := r.Summary
	fmt.Fprintf(w, "\n%d repositories · %d ahead · %d behind · %d dirty · %d untracked\n",
		s.Repos, s.AheadRepos, s.BehindRepos, s.DirtyRepos, s.UntrackedRepos)
}
