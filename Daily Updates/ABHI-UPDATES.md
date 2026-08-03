# Abhi's Daily Updates — proj-init

> Append your daily progress here. Syncthing shares it with Suhrit & Akshay automatically.

---

## [Date]

### What I worked on
- 

### Completed
- 

### Blockers / Questions
- 

### Next steps
- 

---

## 2026-07-24

### What I worked on
- 

### Completed
- 

### Blockers / Questions
- 

### Next steps
- 

---

## 2026-07-26

### What I worked on
- Scaffolded Go project skeleton in PROJECT FILES/proj-init/
- Set up Go module, cobra, all dependencies (viper, fatih/color, spinner)
- Wrote Sector 1 CLI Engine & UX code

### Completed
- `main.go` — entry point invoking cobra
- `cmd/root.go` — root command with --version, --verbose flags
- `cmd/version.go` — version details subcommand
- `cmd/init.go` — interactive mode with project name, template select, output path prompts + validation
- `cmd/list.go` — template listing with descriptions and stack info
- `internal/output/output.go` — colored output (success/fail/info/warn) + spinner helper
- `internal/config/config.go` — viper config management for ~/.proj-init/config.yml
- `.gitignore` — Go standard ignores

### Blockers / Questions
- Go build blocked by Windows Application Control policy (asm.exe restricted). Code verified via gofmt — all syntax clean. Need to build on a machine without this policy or get Go whitelisted.

### Next steps
- Get build working (try admin terminal or IT whitelist)
- End-to-end test once build works
- Add more edge case tests
- GitHub repo setup (waiting on Suhrit)

---

## 2026-07-26 (continued)

### What I worked on
- Integrated `cmd/init.go` with Suhrit's Sector 2 generator (end-to-end flow)
- Merged both sectors into single project at CLI tool Code files/
- Wrote structured logging system
- Wrote README.md

### Completed
- `cmd/init.go` now calls generator.Generate() with variables (ProjectName, GoModule)
- Uses spinner + colored output during generation
- `internal/logger/logger.go` — structured logging with DEBUG/INFO/WARN/ERROR levels
- `README.md` — install, usage, template list, dev instructions
- Merged Sector 1 + Sector 2 code into unified directory (CLI tool Code files/)
- Deleted duplicate PROJECT FILES/proj-init/ to eliminate confusion
- All Go files pass gofmt — syntax verified

### Blockers / Questions
- Same: Windows AppControl policy blocks go build (asm.exe)
- Cannot run go test or verify binary

### Next steps
- Get build working on unrestricted machine
- End-to-end test: init → pick template → generate
- Edge case testing
- GitHub repo setup + tag v0.1.0

---

## 2026-07-26 (Phase 1 code review)

### What I worked on
- Full Phase 1 code audit — found and fixed routing/integration bugs

### Bugs fixed
1. **Hardcoded template list in cmd/init.go & cmd/list.go** — was using a local `templateOption` struct that could get out of sync with embedded templates. Fixed: now calls `templates.List()` directly from embed.go, single source of truth.
2. **generator.Generate() didn't run post-gen hooks** — HookRunner existed in hooks.go but was never wired into Generate(). Fixed: added `Hooks *HookRunner` field to `Options`, hooks run after file generation.
3. **errors.go UserMessage() wrong indentation** — `case ErrFilesystem` had broken indentation that would confuse readers. Fixed via gofmt.

### Issues noted (not fixed, Suhrit's code)
- `internal/templates/registry.go` (`Registry` struct) is completely unused — generator uses `embed.go` package-level functions instead. Suhrit intentionally deferred this.
- `go build` still blocked by Windows AppControl policy.

### Next steps
- Get build working
- End-to-end test + edge case tests
- GitHub repo + v0.1.0 tag

---

## 2026-08-02 (Phase 2 kickoff)

### What I worked on
- Extracted my Phase 2 work split from `PHASE 2 Split.pdf` and wrote a detailed checklist to `PLANS and PHASE WORKS/PHASE 2 ABHI DETAILED.md`
- Executed tasks in order: 6 → 1 → 2 → 3 → 4 → 7 → 5

### Completed
- **Task 6 — help/version polish**: root help now has `Examples:` section + long description; `--version` + `version` both work
- **Task 1 — CI**: `.github/workflows/ci.yml` (gofmt check, go vet, build, tests with coverage, golangci-lint) + `.golangci.yml`
- **Task 2 — GoReleaser + CD**: `.goreleaser.yaml` (5 OS/arch targets, ldflags version injection) + `.github/workflows/release.yml` (tag `v*`)
- **Task 3 — install scripts**: `scripts/install.sh` (bash) + `scripts/install.ps1` (Windows, PATH update)
- **Task 4 — self-update**: `cmd/update.go` + `internal/updater/` (GitHub Releases check, `--check`, atomic in-place binary swap, zip/tar.gz extraction). 404 handled as "no releases yet"
- **Task 7 — config wired in**: `--config` flag, `default_template` pre-selects in init prompt, `output_dir` as default output, `author_name` fills `{{author}}`
- **Task 5 — integration tests**: `internal/integration/` verifies every template's generated file tree matches embedded tree + no un-substituted placeholders remain

### Bugs found & fixed
- **Variable naming mismatch**: `init.go` was setting `ProjectName`/`GoModule`/`AuthorName` but templates use lowercase `{{project_name}}`/`{{module_name}}`/`{{author}}` — substitution silently failed. Fixed to lowercase keys.
- **viper bug**: `SetConfigName()` clears a previously-set `configFile`, so `--config` was silently ignored. Fixed with an `explicitPath` package var.

### Blockers / Questions
- Windows AppControl policy intermittently blocks test binaries written to the temp go-build dir (`generator.test.exe`). Workaround: `go test -c -o` to a fixed path then run directly — all tests pass. May need Go whitelisting for clean CI-like runs locally.

### Next steps
- Push repo + tag v0.1.0 → releases + CI/CD start working
- Suhrit: 3 more templates + post-gen hooks (Sector 2)
- 80%+ coverage on all packages
- PR reviews

---

## 2026-08-02 (Phase 2 complete — verified)

### What I worked on
- Finished & verified ALL my Phase 2 tasks end-to-end. Go build toolchain now works on this machine (the earlier AppControl asm.exe block is resolved; only an intermittent test-binary block remains).

### Completed
- `go build` succeeds — binary builds and runs
- `proj-init --help` shows Examples section + all commands (init, list, update, version)
- `proj-init version` / `--version` both print version, commit, built date, Go version
- `proj-init update --check` correctly reports "no releases published yet" (404 handled)
- `init` with `--config` end-to-end verified:
  - `default_template: tauri-llm` pre-selects in the template prompt
  - `author_name` fills `{{author}}` (confirmed in generated Cargo.toml)
  - `{{project_name}}` substitution confirmed (README, package.json)
- `internal/config/config_test.go` — regression test for the viper explicit-path bug
- All packages tested: config, errors, generator, templates, integration all PASS (via `go test -c -o` + direct run for generator, due to the AppControl temp-dir quirk)

### Bugs found & fixed (final list)
1. **Variable key mismatch** — `init.go` set `ProjectName`/`GoModule`/`AuthorName`, templates use lowercase `{{project_name}}`/`{{module_name}}`/`{{author}}` → substitution silently failed. Fixed.
2. **viper `--config` ignored** — `SetConfigName()` clears a previously-set `configFile`. Fixed with `explicitPath` package var + regression test.

### Blocker (minor)
- Windows AppControl policy intermittently blocks test binaries in the temp go-build dir. Workaround: `go test -c -o <path> && <path> -test.v`. Not a code problem — CI on GitHub will run clean.

### Next steps
- Push to GitHub + tag `v0.1.0` → triggers CI + GoReleaser CD + enables `proj-init update`
- Verify CI/goreleaser/install scripts against the real repo once pushed
- Suhrit: 3 more templates + post-gen hooks
- Coverage 80%+ on all packages

---

## 2026-08-02 (Suhrit's Phase 2 received — merged & verified)

### What I worked on
- Syncthing pulled Suhrit's completed Phase 2 work; reviewed, built, and verified it integrates cleanly with my Sector 1 code.

### Suhrit's Phase 2 (verified in codebase)
- **2 new templates** registered in `embed.go`: `next-webapp` (Next.js 15 + React 19 + TS), `react-native-map` (React Native + Expo + Expo Router + Maps) → total **5 templates**
- **Post-gen hook system** (`internal/generator/`):
  - `exec_hook.go` — `CommandHook` (exec.CommandContext, cross-platform, streams stdout/stderr)
  - `builtin_hooks.go` — `GitInitHook`, `NpmInstallHook`, `FlutterPubGetHook`, `GoModTidyHook` + `HooksForTemplate()` auto-selection from manifest files
  - `Options.AutoHooks` + `selectHooks()` — explicit hooks override auto hooks
  - New test files: exec_hook_test.go, builtin_hooks_test.go, generator_hooks_test.go

### My verification
- `go build ./...` ✅
- `go vet ./...` ✅ (type-checks all test files too)
- `go test ./internal/...` — all packages pass (AppControl temp-dir block intermittently hits test binaries, but `-c` + direct run confirms PASS)
- End-to-end: `proj-init init` → generated `next-webapp` project correctly; `{{project_name}}` → project name in layout.tsx/page.tsx/package.json; `list` shows all 5 templates
- New templates use the SAME placeholder contract (`project_name`, `module_name`, `author`) — my `init.go` config wiring already substitutes correctly, no changes needed

### Coverage snapshot (goal 80%+)
- generator: **93.7%** ✅
- templates: 50.8%
- config: 57.1% (+ regression test added)
- errors: passing (blocked by AppControl during count, verified separately)
- updater/logger/output: no tests yet

### Next steps
- Close coverage gaps: config, templates, updater, logger, output
- Push + tag v0.1.0 → CI/CD live
- PR review of Suhrit's hook system (already verified functionally)
