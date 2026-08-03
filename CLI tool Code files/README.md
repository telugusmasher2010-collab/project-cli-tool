# proj-init

Scaffold any project stack in 3 seconds.

## Install

proj-init is available on every major package manager:

```bash
# Homebrew (macOS / Linux)
brew install telugusmasher2010-collab/tap/proj-init

# Scoop (Windows)
scoop bucket add proj-init https://github.com/telugusmasher2010-collab/scoop-proj-init
scoop install proj-init

# AUR (Arch Linux)
yay -S proj-init

# npm / npx (any platform)
npm i -g proj-init
# or run without installing:
npx proj-init

# Go (development)
go install github.com/telugusmasher2010-collab/project-cli-tool@latest
```

Or use the one-line installers (see `dist/`):

```bash
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/telugusmasher2010-collab/project-cli-tool/main/scripts/install.sh | bash

# Windows (PowerShell)
irm https://raw.githubusercontent.com/telugusmasher2010-collab/project-cli-tool/main/scripts/install.ps1 | iex
```

## Usage

```bash
# Interactive mode
proj-init init

# Scaffold into a specific directory (overrides config output_dir)
proj-init init --output ./my-app

# List available templates
proj-init list

# Version info
proj-init version
proj-init --version

# Update proj-init to the latest release
proj-init update
proj-init update --check

# Verbose output
proj-init init --verbose
```

## Config

proj-init reads `~/.proj-init/config.yml` when present:

```yaml
default_template: tauri-llm   # pre-selected template in init
output_dir: ~/projects         # default output location
author_name: "Your Name"       # fills the {{author}} template variable
```

Override the config file path with `--config <path>`.

## Templates

| Template | Stack | Status |
|----------|-------|--------|
| tauri-llm | Tauri v2 + Rust + React + local LLM | ✅ |
| whatsapp-bot | Node.js + Baileys + SQLite + Fastify | ✅ |
| expense-splitter | Flutter + Dart + Supabase + UPI | ✅ |
| next-webapp | Next.js 15 + Prisma + Tailwind + Auth | 🚧 |
| react-native-map | React Native + Expo + MapLibre | 🚧 |
| cli-go | Minimal Go CLI with cobra | 🚧 |

## Development

```bash
go build -o proj-init.exe .
go test ./...
```

## License

MIT
