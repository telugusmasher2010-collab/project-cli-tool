# PHASE 3 — ABHI DETAILED WORK SPLIT

Derived from `PHASE 3 Split.pdf`. Abhi owns distribution. Suhrit owns validation/demo/marketing.

> Note: most package managers require the GitHub repo to be live + a `v1.0.0` release
> (GoReleaser artifacts) before installers resolve. Create all configs now, publish when tagged.

---

## Task 1 — Homebrew tap & formula
**Files:** `dist/homebrew/proj-init.rb`

- [ ] Formula `class ProjInit < Formula`
- [ ] `url` points at GoReleaser tarball: `proj-init_<ver>_darwin_amd64.tar.gz` (and arm64)
- [ ] `sha256` per bottle (fill after release)
- [ ] Install step: `bin.install "proj-init"`
- [ ] Depends: Go >= 1.21 (build from source fallback) OR use release bottles only
- [ ] README note: `brew install telugusmasher2010-collab/tap/proj-init`

## Task 2 — Scoop bucket
**Files:** `dist/scoop/proj-init.json`

- [ ] Scoop manifest with version, description, license
- [ ] `architecture` → `64bit`/`arm64` URLs to windows zip from GoReleaser
- [ ] `bin: ["proj-init.exe"]`
- [ ] `autoupdate` block using `$version`
- [ ] README note: `scoop bucket add proj-init ...; scoop install proj-init`

## Task 3 — AUR package
**Files:** `dist/aur/PKGBUILD`

- [ ] `pkgname=proj-init`, `pkgver` = release version
- [ ] `source` → linux tar.gz from GoReleaser, `sha256sums` (fill after release)
- [ ] `package()` installs binary to `/usr/bin/proj-init`
- [ ] `depends`: glibc (runtime)
- [ ] README note: `yay -S proj-init`

## Task 4 — npm distribution (npx)
**Files:** `dist/npm/` (separate npm package)

- [ ] `package.json` with bin `proj-init` → `cli.js`, name e.g. `proj-init`
- [ ] `cli.js` — wrapper that downloads the platform binary from GitHub releases on first run (or a `postinstall` install script)
- [ ] Map `process.platform` + `process.arch` → GoReleaser asset name
- [ ] Cache binary in `~/.cache/proj-init/`
- [ ] Support `PROJ_INIT_VERSION` env override + `--version` passthrough
- [ ] README note: `npx proj-init` / `npm i -g proj-init`

## Task 5 — v1.0.0 release
**Files:** release process

- [ ] Bump `Version` in `cmd/root.go` to `1.0.0`
- [ ] Tag `v1.0.0` on main
- [ ] GoReleaser produces archives + checksums
- [ ] Fill hashes into Homebrew/Scoop/AUR configs
- [ ] Publish: Homebrew tap, Scoop bucket, AUR, npm
- [ ] Smoke-test each: `brew install`, `scoop install`, `yay -S`, `npx proj-init`

---

## SUHRIT — Phase 3 (reference)
- Template validation tests (every template compiles)
- Edge cases: existing dirs, overwrite protection
- Demos: scaffold a real project, fix bugs found
- Demo GIF for README
- Marketing: Reddit, Hacker News, Dev.to

## BOTH
- Final docs + README polish
- Bug bash before release
