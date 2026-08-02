# PHASE 2 — ABHI DETAILED WORK SPLIT

Derived from `PHASE 2 Split.pdf`. Abhi owns Sector 1 (CLI Engine & UX). Suhrit owns Sector 2 (templates + hooks).

Suggested execution order: 6 -> 1 -> 2 -> 3 -> 4 -> 7 -> 5

---

## Task 6 — `--help` polish + `--version`
**Files:** `cmd/root.go`, `cmd/version.go`

- [ ] Custom help template: usage, commands, flags grouped, `Examples:` section
- [ ] `proj-init --version` and `proj-init version` both show commit / date / Go version
- [ ] Verify ldflags wiring in release build (`Version`, `Commit`, `Date`, `GoVersion`)

## Task 1 — GitHub Actions CI (test on push, lint)
**Files:** `.github/workflows/ci.yml`, `.golangci.yml`

- [ ] Trigger on `push` + `pull_request`
- [ ] Job: setup Go 1.26 -> `go mod download` -> `go vet ./...` -> `gofmt` check -> `go build ./...` -> `go test ./... -cover`
- [ ] `golangci-lint` action with `.golangci.yml` (errcheck, govet, ineffassign, staticcheck)
- [ ] Cache Go modules

## Task 2 — Goreleaser + CD (release on tag)
**Files:** `.goreleaser.yaml`, `.github/workflows/release.yml`

- [ ] Build matrix: windows/amd64, linux/amd64, darwin/amd64, linux/arm64, darwin/arm64
- [ ] Release triggered by tag `v*`
- [ ] `ldflags` injection: version, commit, date, Go version
- [ ] Publish binaries + checksums to GitHub Releases

## Task 3 — Install Scripts (bash + PowerShell)
**Files:** `scripts/install.sh`, `scripts/install.ps1`

- [ ] `install.sh`: detect OS/arch, download latest release from GitHub, install to `/usr/local/bin`, chmod +x
- [ ] `install.ps1`: download `.exe`, install to `%LOCALAPPDATA%\proj-init`, add to PATH
- [ ] Idempotent — safe to re-run

## Task 4 — `proj-init update` (self-update)
**Files:** `cmd/update.go`, `internal/updater/updater.go`

- [ ] New cobra subcommand `proj-init update`
- [ ] Check latest version via GitHub Releases API
- [ ] Compare against current `Version`
- [ ] Download + replace binary in-place (backup + atomic rename)
- [ ] `proj-init update --check` prints whether update is available

## Task 7 — Config file support (`~/.proj-init/config.yml`)
**Files:** `internal/config/config.go`, `cmd/root.go`, `cmd/init.go`

- [ ] Load config at startup (`initConfig`) via `config.Load()`
- [ ] `default_template` -> pre-select in `init` if set
- [ ] `output_dir` -> default output path
- [ ] `author_name` -> auto-set `AuthorName` variable in generator vars
- [ ] Support override flag `--config <path>`

## Task 5 — Integration Tests (scaffold + build in CI)
**Files:** `internal/integration/integration_test.go`

- [ ] Table-driven: for each embedded template run `Generate()` into temp dir, assert output tree matches, spot-check variable substitution
- [ ] CI job: scaffold with CLI, verify exit code 0 + files exist
- [ ] Test `init` command end-to-end with piped stdin (name + template selection)

---

## SUHRIT — Phase 2 (reference)
- 3 more templates: `expense-splitter`, `next-webapp`, `react-native-map`
- Post-gen hooks: `go mod init`, `npm install`, `git init`

## BOTH
- Coverage 80%+ on all packages
- Code reviews on every PR
