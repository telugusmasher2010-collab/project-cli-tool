---
title: proj-init — Developer CLI Tool Blueprint
created: 2026-07-21
updated: 2026-07-21
type: project
tags: [project, cli-tool, blueprint, go]
sources: []
---
# proj-init — Developer CLI Tool

> **Team:** Abhi + Suhrit + Akshay
> **Timeline:** 3–4 weeks (first 2–5 days: just Abhi + Suhrit, Akshay joins later)
> **Language:** Go 🐹
> **Public name:** `proj-init`
> **Tagline:** Scaffold any project stack in 3 seconds.

---

## What It Does (For the Public)

A single CLI command that scaffolds a complete, working project stack in **3 seconds**.

```bash
# Interactive mode
proj-init init

# Direct mode — pick a template
proj-init tauri-llm
proj-init whatsapp-bot
proj-init expense-splitter

# List available templates
proj-init list

# Version info
proj-init --version
```

**What the user gets:** A ready-to-code folder with `src/`, CI/CD, Dockerfile, README with badges, `.gitignore`, LICENSE — everything compiles and runs first time.

---

## Why This Matters to You (Your Daily Life)

| Before (without proj-init) | After (with proj-init) |
|---|---|
| 30–60 mins of repetitive setup per project | 3 seconds — `proj-init stack-name` |
| Avoid starting new ideas because setup is painful | Experiment freely — zero friction |
| Templates live in your head (forget details) | Templates live in code (always fresh) |
| Friends ask "how do I start?" — you send blog links | Friends run `proj-init` and have a working project |
| Your setup knowledge is trapped in your brain | Your best setups are shared with the team |

**The flex:** "I built a CLI tool that other developers use" = top 1% hiring signal.

---

## Tech Stack

| Layer | Technology | Why |
|---|---|---|
| **Core CLI** | **Go** (`cobra`) | Simple, fast compiles, easy for all 3 to contribute equally |
| **Template Engine** | Go `text/template` + `embed` | Built-in stdlib, no dependencies |
| **Config** | Go `viper` | YAML/TOML/JSON config files |
| **Cross-Compile** | `goreleaser` | One command → binaries for linux/windows/mac/arm64 |
| **CI/CD** | GitHub Actions | Build on push, release on tag |
| **Distribution** | Homebrew / Scoop / AUR / npm | Users install with one command |
| **Docs** | Markdown + GitHub Pages | Simple, no framework needed |

### Why Go Over Rust

| Factor | Go 🐹 | Rust 🦀 |
|---|---|---|
| Learning curve | **Low** — productive in 2 days | High — borrow checker slows first 2 weeks |
| Team ramp-up | Abhi, Suhrit, Akshay all coding **immediately** | 2 weeks fighting the compiler |
| Compile speed | **Seconds** | Minutes |
| Binary size | ~10-15 MB (still single binary) | ~5 MB |
| CLI libs | `cobra` — mature, simple | `clap` — powerful but complex derive macros |
| Release pipeline | `goreleaser` — one-command publish | `cargo-dist` — good but newer |
| Hiring signal | Good (Go is respected) | Better (Rust is elite) |

**Verdict:** Go lets 3 friends ship a working tool in 3-4 weeks instead of spending half that time learning Rust. Ship first, rewrite in Rust later if needed.

---

## Project Structure

```
proj-init/
│
├── main.go                    # Entry point
├── go.mod / go.sum            # Go module
│
├── cmd/
│   ├── root.go                # Root command (cobra)
│   ├── init.go                # proj-init init — interactive mode
│   ├── list.go                # proj-init list — show templates
│   └── version.go             # proj-init --version
│
├── internal/
│   ├── generator/
│   │   ├── generator.go       # Template engine wrapper
│   │   └── vars.go            # Variable resolution
│   ├── templates/
│   │   ├── embed.go           # Embedded template filesystem
│   │   ├── tauri-llm/         # Template files for Tauri + LLM
│   │   ├── whatsapp-bot/      # Template files for WhatsApp bot
│   │   ├── expense-splitter/  # Template files for Flutter expense app
│   │   └── ...                # More templates
│   ├── config/
│   │   └── config.go          # Config file management
│   ├── output/
│   │   └── output.go          # Colored output, spinners
│   └── errors/
│       └── errors.go          # Custom error types
│
├── scripts/
│   ├── install.sh             # Unix install script
│   └── install.ps1            # Windows install script
│
├── .github/workflows/
│   ├── build.yml              # CI — build & test on push
│   └── release.yml            # CD — release on tag via goreleaser
│
├── .goreleaser.yaml           # Goreleaser config
├── README.md                  # Main docs
├── CONTRIBUTING.md            # Contributor guide
├── LICENSE                    # MIT / Apache 2.0
└── .gitignore
```

---

## The 3 Equal Work Sectors ⚖️

Every sector has **~equal complexity, ~equal LoC, ~equal difficulty**. Nobody gets the easy part.

### Sector ① — CLI Engine & UX

| What | Tech |
|---|---|
| Arg parsing & flag handling | `cobra` |
| Interactive `init` prompts | Go `bufio` / `survey` lib |
| Colored output, spinners | `fatih/color`, `briandowns/spinner` |
| Config file management | `viper` |
| Error handling & exit codes | Custom error types |
| `--help`, `--version`, `list` | Cobra built-in |

**Who builds it:** One team member (you or Suhrit)

### Sector ② — Template System & Generators

| What | Tech |
|---|---|
| Template engine wrapper | Go `text/template` + `embed` |
| Variable substitution system | Custom `vars.go` |
| Template files (all stacks) | Go `embed` FS |
| `proj-init init` generation logic | `generator.go` |
| Post-gen hooks (go mod init, npm install) | Shell exec from Go |
| Template validation | Unit tests |

**Who builds it:** The other early member (you or Suhrit)

### Sector ③ — Distribution & CI/CD

| What | Tech |
|---|---|
| GitHub Actions — build on push | `.github/workflows/build.yml` |
| GitHub Actions — release on tag | `.github/workflows/release.yml` |
| Cross-compilation matrix | `goreleaser` (linux/windows/mac/arm64) |
| Package manager publishing | Homebrew, Scoop, AUR, npm |
| Versioning & changelogs | Conventional commits + auto-changelog |
| Install scripts (bash + PowerShell) | `scripts/install.sh`, `scripts/install.ps1` |
| Self-update mechanism | `update.go` |

**Who builds it:** Akshay (joins later, but this sector is as meaty as the other two)

---

## The Contract Between Sectors (Clean Interfaces)

```
Sector ① (CLI Engine)
    │  Parses flags, collects variables
    │  Calls Sector ② with template name + variables
    ▼
Sector ② (Generator)
    │  Reads embedded templates
    │  Substitutes variables
    │  Writes output to target directory
    │  Returns path or error to Sector ①
    ▼
Sector ① displays result to user

Sector ③ (CI/CD)
    │  Builds the binary from Sector ① + ② combined
    │  Cross-compiles & publishes
    │  Completely independent — no code coupling
```

**No one blocks anyone.** Each sector has clear boundaries.

---

## Sprint Plan (3-4 Weeks)

### Days 1-2 — Foundation (Abhi + Suhrit)

| Task | Owner |
|---|---|
| Decide final project name | Both |
| Create GitHub repo + add collaborators | Both |
| Scaffold Go module + cobra skeleton | Both pair |
| Sector ①: basic `init` command with flag parsing | One person |
| Sector ②: embed first template + basic generator | Other person |
| Sector ③: set up basic GitHub Actions build | Whoever finishes first |
| Set up Syncthing on both laptops | Both |

### Days 3-5 — Core Features (Abhi + Suhrit)

| Task | Owner |
|---|---|
| Sector ①: interactive init prompts, colored output, spinners | Person A |
| Sector ②: 2-3 working templates, variable system, post-gen hooks | Person B |
| Sector ① + ② integration — end-to-end flow works | Both pair |
| Basic error handling & edge cases | Both split |

### Week 2 — Akshay Joins, Polish Core

| Task | Owner |
|---|---|
| Sector ③: `goreleaser` config, cross-compile matrix, install scripts | Akshay |
| Sector ③: GitHub Actions release pipeline + auto-changelog | Akshay |
| Sector ①: config file support (`~/.proj-init/config.yml`) | Person A |
| Sector ②: add 3 more templates | Person B |
| Integration testing on all 3 laptops | All |

### Week 3 — Distribution & Quality

| Task | Owner |
|---|---|
| Sector ③: Homebrew tap, Scoop bucket, AUR package | Akshay |
| Sector ③: npm distribution | Akshay |
| Sector ①: `--help` polish, error messages, output formatting | Person A |
| Sector ②: template validation & unit tests | Person B |
| README, docs site, demo GIF | All split |

### Week 4 — Shipping & Beyond

| Task | Owner |
|---|---|
| Sector ③: First actual release to package managers | Akshay |
| Bug fixes & edge cases | All |
| Dogfood — use it for your own projects | All |
| Collect feedback, plan v2 features | All |

---

## Template Roadmap (Stacks to Ship)

| Template | Stack | Priority |
|---|---|---|
| `tauri-llm` | Tauri v2 + Rust + React + local LLM sidecar | ⭐ High (your ORION stack) |
| `whatsapp-bot` | Node.js + Baileys + SQLite + Fastify | ⭐ High (revenue project) |
| `expense-splitter` | Flutter + Dart + Supabase + UPI | ⭐ High |
| `next-webapp` | Next.js 15 + Prisma + Tailwind + Auth | Medium |
| `react-native-map` | React Native + Expo + MapLibre | Medium |
| `cli-go` | Minimal Go CLI with cobra | Low (meta-template) |

---

## Success Criteria

- [ ] `proj-init list` shows available templates with descriptions
- [ ] `proj-init tauri-llm` creates a complete, compilable project in < 3 seconds
- [ ] Auto-generates README with badges, LICENSE, .gitignore, CI/CD
- [ ] Cross-compiles for Windows, macOS, Linux (amd64 + arm64)
- [ ] Published on at least 2 package managers (Homebrew + Scoop)
- [ ] All 3 team members contributed meaningful code
- [ ] Each sector has equal LoC and complexity (±20%)

---

## Open Questions (Decide With Team)

- [ ] **Project name** — finalize `proj-init` or something else?
- [ ] **License** — MIT or Apache 2.0?
- [ ] **GitHub org or personal account?**
- [ ] **First template to ship** — prioritize ORION stack or something simpler?
- [ ] **Do we dogfood immediately?** — use proj-init to scaffold its own next version?

---

> **Next action:** Pick Go, create GitHub repo, scaffold the skeleton. Let's go. 🚀
