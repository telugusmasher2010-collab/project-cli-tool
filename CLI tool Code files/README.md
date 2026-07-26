# proj-init

Scaffold any project stack in 3 seconds.

## Install

```bash
go install github.com/telugusmasher2010-collab/project-cli-tool@latest
```

## Usage

```bash
# Interactive mode
proj-init init

# List available templates
proj-init list

# Version info
proj-init version
proj-init --version
```

## Templates

| Template | Stack |
|----------|-------|
| tauri-llm | Tauri v2 + Rust + React + local LLM |
| whatsapp-bot | Node.js + Baileys + SQLite + Fastify |
| expense-splitter | Flutter + Dart + Supabase + UPI |
| next-webapp | Next.js 15 + Prisma + Tailwind + Auth |
| react-native-map | React Native + Expo + MapLibre |
| cli-go | Minimal Go CLI with cobra |

## Development

```bash
go build -o proj-init.exe .
go test ./...
```

## License

MIT
