# {{project_name}}

> Tauri v2 + React + Rust + Local LLM sidecar

[![Built with Tauri](https://img.shields.io/badge/Built%20with-Tauri_v2-orange)](https://tauri.app)
[![Rust](https://img.shields.io/badge/Rust-stable-blue)](https://www.rust-lang.org)
[![React](https://img.shields.io/badge/React-18-blue)](https://react.dev)

## Overview

{{project_name}} is a desktop application built with Tauri v2, featuring a React frontend and Rust backend with local LLM inference capabilities.

## Tech Stack

| Layer | Technology |
|-------|------------|
| Desktop Framework | Tauri v2 |
| Frontend | React 18 + TypeScript |
| Backend | Rust |
| Build Tool | Vite |
| LLM | Local sidecar integration |

## Prerequisites

- [Rust](https://www.rust-lang.org/tools/install) (stable)
- [Node.js](https://nodejs.org/) (v18+)
- [pnpm](https://pnpm.io/) (recommended) or npm

## Getting Started

```bash
# Install dependencies
pnpm install

# Start development server
pnpm tauri dev

# Build for production
pnpm tauri build
```

## Project Structure

```
{{project_name}}/
├── src/                # React frontend
├── src-tauri/          # Rust backend
│   ├── Cargo.toml      # Rust dependencies
│   ├── tauri.conf.json # Tauri configuration
│   └── src/
│       └── main.rs     # Rust entry point
├── package.json
└── README.md
```

## Development

The frontend runs on Vite dev server, and Tauri spawns a native webview to display the app. During development, hot module replacement works seamlessly.

```bash
pnpm tauri dev
```

## Building

```bash
# Build for current platform
pnpm tauri build

# The output binary will be in src-tauri/target/release/
```

## License

MIT
