// Package git wraps the git binary for the queries repos-sync needs.
// Never invokes a shell.
package git

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Probe struct {
	Bare        bool
	Branch      string
	HasUpstream bool
	Upstream    string
	Ahead       *int
	Behind      *int
	Dirty       int
	Untracked   int
	LastAt      string
	LastEpoch   *int64
	LastSubject string
	Err         string
}

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return strings.TrimSpace(out.String()), err
	}
	return strings.TrimSpace(out.String()), nil
}

func IsBare(dir string) (bool, error) {
	out, err := run(dir, "rev-parse", "--is-bare-repository")
	if err != nil {
		return false, err
	}
	return out == "true", nil
}

func IsRepoRoot(dir string) bool {
	_, err := run(dir, "rev-parse", "--git-dir")
	return err == nil
}

func ProbeRepo(dir string) Probe {
	p := Probe{}
	head, err := run(dir, "rev-parse", "--verify", "HEAD")
	if err != nil {
		p.Err = "empty repository"
		// still harvest branch + bare so the row is meaningful
		if out, e := run(dir, "rev-parse", "--is-bare-repository"); e == nil {
			p.Bare = out == "true"
		}
		if out, e := run(dir, "symbolic-ref", "--short", "HEAD"); e == nil {
			p.Branch = out
		}
		return p
	}
	_ = head

	if out, e := run(dir, "rev-parse", "--is-bare-repository"); e == nil {
		p.Bare = out == "true"
	}
	if out, e := run(dir, "symbolic-ref", "--short", "HEAD"); e == nil {
		p.Branch = out
	} else {
		if out, e2 := run(dir, "rev-parse", "--short", "HEAD"); e2 == nil {
			p.Branch = "(detached " + out + ")"
		}
	}
	if out, e := run(dir, "rev-parse", "--abbrev-ref", "@{upstream}"); e == nil && out != "" {
		p.HasUpstream = true
		p.Upstream = out
		ahead, e1 := run(dir, "rev-list", "--count", "@{u}..HEAD")
		behind, e2 := run(dir, "rev-list", "--count", "HEAD..@{u}")
		if e1 == nil {
			if n, x := strconv.Atoi(ahead); x == nil {
				p.Ahead = &n
			}
		}
		if e2 == nil {
			if n, x := strconv.Atoi(behind); x == nil {
				p.Behind = &n
			}
		}
	}
	if !p.Bare {
		if out, e := run(dir, "status", "--porcelain=v1"); e == nil {
			for _, line := range strings.Split(out, "\n") {
				if line == "" {
					continue
				}
				if strings.HasPrefix(line, "??") {
					p.Untracked++
				} else {
					p.Dirty++
				}
			}
		} else {
			p.Err = strings.TrimSpace(out + " (git status failed)")
		}
	}
	if out, e := run(dir, "log", "-1", "--format=%cI %ct %s"); e == nil {
		fields := strings.SplitN(out, " ", 3)
		if len(fields) == 3 {
			p.LastAt = fields[0]
			if ts, x := strconv.ParseInt(fields[1], 10, 64); x == nil {
				p.LastEpoch = &ts
			}
			p.LastSubject = fields[2]
		}
	}
	return p
}

func GitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
