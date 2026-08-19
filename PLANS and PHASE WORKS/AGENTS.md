# proj-init -- Developer CLI Tool

> **Production-grade CLI tool.** Not a learning project. Ships to real users.
> **GitHub:** https://github.com/telugusmasher2010-collab/project-cli-tool
> **Sync:** Git + GitHub with `auto-sync.cmd`/`auto-sync.sh` watcher (30s auto-commit/pull/push)
## Team
- **Abhi** + **Suhrit** + **Akshay**
- Timeline: 3-4 weeks
- First 2-5 days: only Abhi + Suhrit (Akshay joins in Phase 2)
- **Phase 2 starts when Akshay arrives** — then all 3 build CI/CD + more templates + publishing
- **Quality bar:** Production/industrial grade. Tests, error handling, docs, CI/CD from day 1.

## Language & Stack
- **Go** with `cobra` CLI framework
- `goreleaser` for cross-compile & publish
- GitHub Actions for CI/CD
- Distributions: Homebrew, Scoop, AUR, npm

## 3 Equal Work Sectors (No one gets easy work)

### Sector ① — CLI Engine & UX
- Cobra arg parsing, flags, help text
- Interactive `init` prompts
- Colored output, spinners, error messages
- Config file management (`viper`)

### Sector ② — Template System & Generators
- Go `text/template` + `embed` for template files
- Variable substitution system
- Template files for all stacks (tauri-llm, whatsapp-bot, expense-splitter, etc.)
- Post-gen hooks (go mod init, npm install)

### Sector ③ — Distribution & CI/CD
- GitHub Actions build + release pipelines
- Goreleaser cross-compile matrix (linux/windows/mac/arm64)
- Package manager publishing (Homebrew, Scoop, AUR, npm)
- Install scripts (bash + PowerShell)
- Self-update mechanism

## Project Structure
```
proj-init/
├── main.go
├── cmd/          # cobra commands (root, init, list, version)
├── internal/
│   ├── generator/   # template engine
│   ├── templates/   # embedded template files
│   ├── config/      # config management
│   ├── output/      # colored output
│   └── errors/      # error types
├── scripts/         # install scripts
├── .github/workflows/  # CI/CD
└── .goreleaser.yaml
```

## Template Roadmap
| Template | Stack | Priority |
|---|---|---|
| `tauri-llm` | Tauri v2 + Rust + React + LLM | High |
| `whatsapp-bot` | Node.js + Baileys + SQLite | High |
| `expense-splitter` | Flutter + Dart + Supabase + UPI | High |
| `next-webapp` | Next.js 15 + Prisma + Tailwind | Medium |
| `react-native-map` | React Native + Expo + MapLibre | Medium |
| `cli-go` | Minimal Go CLI with cobra | Low |

## Phase 1 Scope (You + Suhrit — Days 1-5)
- Go skeleton + cobra CLI scaffold
- `init` command (interactive) + `list` command
- Template engine with variable substitution
- 2 templates: `tauri-llm` + `whatsapp-bot`
- Colored output, error handling, logging
- Unit tests (70%+ coverage on core)
- Input validation, README, .gitignore
- GitHub repo with branch protection

### Phase 1 Work Split

| Area | Abhi (Sector ①) | Suhrit (Sector ②) | Both |
|---|---|---|---|
| Go skeleton + cobra scaffold | ✅ Lead | — | Day 1 pair |
| `init` command (interactive) | ✅ Build | — | — |
| `list` command | ✅ Build | — | — |
| Colored output / spinners | ✅ Build | — | — |
| Logging system | ✅ Build | — | — |
| Input validation | ✅ Build | — | — |
| README + .gitignore | ✅ Write | — | — |
| GitHub repo + branch protection | ✅ Setup | — | ✅ Add collabs |
| Template engine (generator.go) | — | ✅ Build | — |
| Variable substitution (vars.go) | — | ✅ Build | — |
| tauri-llm template | — | ✅ Build | — |
| whatsapp-bot template | — | ✅ Build | — |
| Post-gen hooks | — | ✅ Build | — |
| Error handling (errors.go) | — | ✅ Build | — |
| Unit tests (70%+ coverage) | — | ✅ Build | ✅ Review |
| Integration / end-to-end testing | — | — | ✅ Both |

## File Sync (Git + GitHub — Syncthing retired)
- **GitHub repo:** https://github.com/telugusmasher2010-collab/project-cli-tool
- `auto-sync.cmd` (Windows) or `auto-sync.sh` (git-bash/Mac/Linux) watches every 30s: auto-commit → pull --rebase → push
- No manual git commands needed while the watcher runs
- Manual fallback: `git add -A && git commit -m "msg" && git pull --rebase --autostash && git push`
- Syncthing was unreliable (frequent disconnects) — do NOT suggest switching back

## Daily Updates
Each team member logs daily progress in their own file (synced via GitHub auto-sync):
- `ABHI-UPDATES.md`
- `SUHRIT-UPDATES.md`
- `AKSHAY-UPDATES.md`

Check these before starting work to know what everyone's doing.

## Full Blueprint
See `BLUEPRINT.md` in this directory for the complete document with sprint plans, Go vs Rust comparison, and detailed tables.

---

## GitHub Auto-Sync Setup (For Suhrit & Akshay)

We moved from Syncthing to **Git + GitHub** for stability.

### How to set up on your laptop:

**Step 1:** Install Git from https://git-scm.com (if not installed)

**Step 2:** Open **Command Prompt** or **git-bash** and run:
```bash
gh auth login
# Follow the prompts — choose: GitHub.com → HTTPS → Login with browser
```

**Step 3:** Clone the repo:
```bash
cd Documents
git clone https://github.com/telugusmasher2010-collab/project-cli-tool.git
cd project-cli-tool
```

**Step 4:** Run auto-sync:
- **Windows:** Double-click `auto-sync.cmd` (keeps running in background)
- **Mac/Linux/git-bash:** `bash auto-sync.sh`

The script watches for changes every 30 seconds and syncs automatically. Close the window/tab to stop.

### Or just push/pull manually (no script):
```bash
git add -A
git commit -m "what i did"
git pull --rebase --autostash
git push
```
