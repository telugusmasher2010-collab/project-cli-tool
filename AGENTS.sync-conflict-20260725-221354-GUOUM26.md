# proj-init — Developer CLI Tool

> Scaffold any project stack in 3 seconds.

## Team
- **Abhi** + **Suhrit** + **Akshay**
- Timeline: 3-4 weeks
- First 2-5 days: only Abhi + Suhrit (Akshay joins later)

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

## File Sync
- Syncthing shared folder: all 3 laptops sync project files
- Abhi's device ID: GUOUM26-NOEMMTU-AYNEGL5-GDHTDKB-YVHPJAO-XPYHMQD-HNPU2O4-UMO5NQ5

## Full Blueprint
See `BLUEPRINT.md` in this directory for the complete document with sprint plans, Go vs Rust comparison, and detailed tables.
