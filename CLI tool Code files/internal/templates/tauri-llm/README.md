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
- npm (ships with Node.js)

## Getting Started

```bash
# Install dependencies
npm install

# Start the desktop app in development mode
npm run tauri dev

# Build for production
npm run tauri build
```

## Project Structure

```
{{project_name}}/
├── index.html            # Vite entry point
├── src/                  # React frontend
│   ├── main.tsx          # React entry point
│   ├── App.tsx           # Root component
│   └── App.css
├── vite.config.ts        # Vite configuration (port 1420)
├── tsconfig.json
├── Cargo.toml            # Rust workspace
├── src-tauri/            # Rust backend
│   ├── Cargo.toml        # Rust dependencies
│   ├── build.rs          # tauri-build build script
│   ├── tauri.conf.json   # Tauri configuration
│   ├── capabilities/     # Tauri permissions
│   └── src/
│       └── main.rs       # Rust entry point
├── package.json
└── README.md
```

## Development

The frontend runs on the Vite dev server (port 1420), and Tauri spawns a native webview to display the app. During development, hot module replacement works seamlessly.

```bash
npm run tauri dev
```

## Building

```bash
# Build for the current platform
npm run tauri build

# The output binary will be in src-tauri/target/release/
```

## License

MIT — {{author}}
