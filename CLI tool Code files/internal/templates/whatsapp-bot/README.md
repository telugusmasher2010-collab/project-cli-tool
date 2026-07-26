# {{project_name}}

> WhatsApp Bot built with Node.js, Baileys, SQLite, and Fastify

[![Node.js](https://img.shields.io/badge/Node.js-20+-green)](https://nodejs.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Overview

{{project_name}} is a production-ready WhatsApp bot framework using the Baileys library for WhatsApp Web protocol interaction, SQLite for persistent storage, and Fastify for a lightweight HTTP API.

## Features

- WhatsApp Web multi-device support via Baileys
- Persistent session management
- SQLite database for messages, contacts, and state
- Fastify HTTP API for bot management and webhooks
- Plugin-based message handler architecture
- Graceful shutdown and reconnection logic

## Prerequisites

- Node.js 20+
- npm or yarn

## Getting Started

1. Install dependencies:

```bash
npm install
```

2. Copy the config example and fill in your settings:

```bash
cp config.example.json config.json
```

3. Start the bot:

```bash
npm start
```

4. Scan the QR code displayed in the terminal with your WhatsApp mobile app.

## Project Structure

```
{{project_name}}/
├── index.js              # Entry point
├── config.json           # Bot configuration (git-ignored)
├── config.example.json   # Example configuration
├── package.json
├── database/
│   └── schema.sql        # SQLite schema
├── src/
│   ├── bot.js            # Baileys client setup
│   ├── handlers/         # Message handler plugins
│   └── db.js             # Database connection
└── README.md
```

## Configuration

Edit `config.json` with your settings. See `config.example.json` for all available options.

| Key | Description | Default |
|-----|-------------|---------|
| `botName` | Display name for the bot | `{{project_name}}` |
| `owner` | Owner JID (phone@c.us) | - |
| `prefix` | Command prefix | `!` |
| `database.path` | SQLite database file path | `./data/bot.db` |

## Development

```bash
# Start with auto-reload
npm run dev
```

## License

MIT — {{author}}
