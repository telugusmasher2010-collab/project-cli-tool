# proj-init — Developer CLI Tool

> **Production-grade CLI tool.** Not a learning project. Ships to real users.

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

## File Sync
- Syncthing shared folder: all 3 laptops sync project files
- Abhi's device ID: GUOUM26-NOEMMTU-AYNEGL5-GDHTDKB-YVHPJAO-XPYHMQD-HNPU2O4-UMO5NQ5

## Daily Updates
Each team member logs daily progress in their own file (synced via Syncthing):
- `ABHI-UPDATES.md`
- `SUHRIT-UPDATES.md`
- `AKSHAY-UPDATES.md`

Check these before starting work to know what everyone's doing.

## Full Blueprint
See `BLUEPRINT.md` in this directory for the complete document with sprint plans, Go vs Rust comparison, and detailed tables.
