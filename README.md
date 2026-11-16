# repos-sync — workspace git health

Which of my repositories need attention right now, and why.

`repos-sync` scans the direct children of its roots for git repositories and
answers, in one shot, for every repo it finds: branch, how far ahead or
behind upstream you are, whether anything is modified or staged, and what
is untracked.

```
$ repos-sync
repo                   branch       ↑    ↓    status
demo-alpha                 main         0    0    clean
demo-beta                 main         2    0    dirty 2, ahead 2
demo-companion                   master       3    0    no upstream
demo-gamma                   master       0    0    untracked 2

34 repositories · 1 ahead · 0 behind · 1 dirty · 1 untracked
```

## Why

One repo at a time, git tells you things via your prompt. One directory of
repos at a time — nothing. repos-sync is that answer for `~/Projects` (and
anything else you point it at), suitable for a terminal, a script, or a
status bar.

It is **read-only**: no `fetch`, no `gc`, no locks, no writes. It runs
`git status --porcelain` with `GIT_OPTIONAL_LOCKS=0`, so it does not even
refresh the index. It never touches the network.

## How it works

Direct children of each root that are git repositories are inspected with a
handful of `exec.Command("git", ...)` calls — never a shell string — under
`GIT_OPTIONAL_LOCKS=0`: `rev-parse`, `symbolic-ref`, `rev-list`,
`status --porcelain`, `log -1`. Repositories that fail to load are reported
with a per-repo `error` and do not abort the scan. Bare repositories are
supported (reported as one entry).

## Requirements

- Linux
- Go 1.27+ (to build) or a release binary
- Git

## Installation

```sh
go install github.com/RakhaYandra/repos-sync@latest
```

Or from a checkout: `make install` (installs to `~/.local/bin/repos-sync`;
override with `PREFIX`).

## Usage

```sh
repos-sync                      # scan ~/Projects (default root)
repos-sync --root ~/code        # repeatable; replaces the default
repos-sync --root ~/code --root ~/.dotfiles
repos-sync --json               # machine output
repos-sync check .              # inspect one repository
repos-sync check ~/Projects/demo-alpha
```

## Commands

| Command | Effect |
|---|---|
| `repos-sync` | Scan all roots (default `~/Projects`), table + one-line summary |
| `repos-sync --root DIR` | Root to scan (repeatable; replaces the default) |
| `repos-sync --json` | JSON envelope (`schema_version: 1`) |
| `repos-sync check [PATH]` | Inspect one repository (default `.`) |
| `repos-sync --help` / `--version` | Help / version |

Kolom tabel: `repo` (basename), `branch` (atau `(detached <sha>)`), `↑`/`↓`
(ahead/behind vs upstream), `status` — `clean` · `dirty N` · `ahead N` ·
`behind N` · `untracked N` · `no upstream` · `error: <reason>`.


A bare repository passed as `--root` is reported as a single entry. Roots are
scanned at depth 1; repositories in subfolders are not discovered. Empty
repositories are reported loudly (`error: empty repository`), not dropped.

## Output contract (`--json`, `schema_version: 1`)

```json
{
  "schema_version": 1,
  "generated_at": "2026-11-16T12:00:00Z",
  "roots": ["/home/dev/projects"],
  "summary": {
    "repos": 34,
    "with_upstream": 33,
    "ahead_repos": 1,
    "behind_repos": 0,
    "unpushed_commits": 2,
    "dirty_repos": 1,
    "untracked_repos": 1
  },
  "repos": [
    {
      "name": "demo-beta",
      "path": "/home/dev/projects/demo-beta",
      "branch": "main",
      "bare": false,
      "upstream": "origin/main",
      "ahead": 2,
      "behind": 0,
      "dirty": 1,
      "untracked": 0,
      "last_commit_at": "2026-11-15T12:00:00Z",
      "last_commit_age_s": 90000,
      "last_commit_subject": "feat: something",
      "error": null
    }
  ]
}
```

Rules consumers can rely on:

- `schema_version` is `1`, always the first key.
- `ahead`/`behind` are `null` (not `0`) when the repository has no upstream.
- `last_commit_at`, `last_commit_age_s`, `last_commit_subject` are `null`
  for empty repositories.
- `error` carries the human-readable reason when a repository could not be
  inspected; the scan continues past failures.
- `untracked` is reported separately and is **not** included in
  `unpushed_commits`.
- Adding keys is minor; removing or renaming is breaking.
- `last_commit_subject` is arbitrary text from commit messages — quote it.

## Sorting

Severity first, then name:

`error > dirty > ahead > behind > no-upstream > untracked > clean`

## Exit codes

`0` ok · `1` usage error (bad flags, missing root) · `2` operational error
(git missing, or a global failure; partial failures still render with rows
marked `error:`).

## FAQ

**Does it fetch?** No. Ever, by design. It reports what git already knows
locally. If you want fresh numbers, run `git fetch` yourself first.

**Does it write anything in my repos?** No. `GIT_OPTIONAL_LOCKS=0` is set
for every git invocation; the index is not refreshed, no lock is taken.
An integration test asserts the repos are byte-identical before and after
a scan.

**Can it hurt a big repo?** `git status --porcelain` is the most expensive
call; on very large trees that can still be noticeable. It never blocks
your work — worst case it takes a second.

## License

MIT
