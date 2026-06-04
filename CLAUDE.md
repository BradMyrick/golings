# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Fork of `golings` (Rustlings-style Go exercises) by BradMyrick. Module path `github.com/bradmyrick/golings`. Go 1.26.1.

## Commands

Run from repo root:

- `make test` — full test suite with coverage (excludes `fixtures/error1`, which is intentionally broken).
- `make watch` — run `golings watch` from source.
- `make run` — run `golings run` from source.
- Single package test: `cd golings && go test -v ./cmd/...` (or `./exercises/...`, `./ui/...`).
- Single test: `cd golings && go test -v ./exercises -run TestName`.
- Build binary: `cd golings && go build -o ../golings`.

Tests use Ginkgo/Gomega (`onsi/ginkgo/v2`). `cmd_suite_test.go` / `exercises_suite_test.go` are the suite entry points.

## Architecture

Three layers under `golings/`:

- `golings/cmd/` — Cobra CLI. `root.go` wires subcommands: `run`, `verify`, `watch`, `list`, `hint`, `print`. Entry is `golings/golings.go` → `cmd.Execute`.
- `golings/exercises/` — exercise model and runner.
  - `exercise.go` defines an exercise (name, path, mode, hint). `mode` is `compile` or `test` — `compile` mode runs `go build`, `test` mode runs `go test`. This determines how `runner.go` invokes the toolchain on the file.
  - `list.go` parses `info.toml` (repo root) which is the **ordered** source of truth for exercises. Order in `info.toml` = progression order.
- `golings/ui/` — `lipgloss`-based TUI for watch mode. `layout.go` renders the main view; `list.go` is the scrollable interactive exercise picker.

State: progress is tracked in a `.golings-state` file at repo root (this fork replaced the upstream `// I AM NOT DONE` marker convention — exercises in this fork generally have no such marker).

Watch loop: `fsnotify` watches `exercises/`, re-runs the current exercise on save, single-key keybinds (`n`/`h`/`l`/`q`) handled in the TUI.

`golings/fixtures/` contains canned exercises used by tests; `fixtures/error1` is intentionally non-compiling and is excluded from `go test ./...` via the Makefile's `grep -v`.

## Adding an exercise

1. Append `[[exercises]]` block to `info.toml` at desired position (order matters).
2. Create file under `exercises/<topic>/`.
3. Pick `mode = "compile"` (just needs to build) or `"test"` (must pass `go test`).

## Repo conventions

- Pinned release in README is `@v0.0.1` — `@latest` was buggy.
- Upstream `mauricioabreu` metadata was scrubbed; this fork is the canonical one.
