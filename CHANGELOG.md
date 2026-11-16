# Changelog

All notable changes to this project are documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.1.1] - 2026-11-17

### Fixed

- `repos-sync --version` now prints the version (previously unknown flag, exit 1).
- Colour output no longer leaks into redirected stdout/stderr — terminal detection uses `ioctl(TCGETS)` instead of `os.ModeCharDevice`.
- Duplicate `--root` values are deduplicated (`filepath.Clean` + set), matching the documented behaviour.

### Added

- Tests for `--version` output and root deduplication.

## [0.1.0] - 2026-11-16

### Added

- `repos-sync` — one-shot workspace git health scan for the direct children
  of one or more roots (default `~/Projects`, override with repeatable
  `--root`).
- Per repo: branch (or `(detached <sha>)`), bare detection, upstream ref,
  ahead/behind counts, dirty vs untracked split, last commit timestamp,
  age and subject.
- `--json` output with `schema_version: 1`; severity-sorted repos; tolerant
  partial failure (`repos[].error`) and `GIT_OPTIONAL_LOCKS=0` throughout.
- `repos-sync check [PATH]` for a single repository.
- MIT license.
