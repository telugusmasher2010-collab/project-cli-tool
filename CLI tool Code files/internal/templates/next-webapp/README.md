# {{project_name}}

> Modern web app built with Next.js 15, React 19, and TypeScript (App Router)

[![Next.js](https://img.shields.io/badge/Next.js-15-black)](https://nextjs.org)
[![React](https://img.shields.io/badge/React-19-blue)](https://react.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue.svg)](https://www.typescriptlang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Overview

{{project_name}} is a modern web application built on the Next.js App Router with React 19 and TypeScript. It ships with a production-ready foundation: strict type checking, hot module replacement, static and server rendering, and a clean minimal starter structure.

## Features

- App Router with file-based routing and server components
- React 19 with the latest rendering features
- TypeScript with strict mode out of the box
- Optimized production builds with `next build`
- Clean, minimal global styles with dark mode support

## Tech Stack

| Layer | Technology |
|-------|------------|
| Framework | Next.js 15 |
| UI | React 19 |
| Language | TypeScript |
| Styling | Global CSS |

## Prerequisites

- Node.js 18.18+ (Node 20+ recommended)
- npm (or your preferred package manager)

## Getting Started

1. Install dependencies:

```bash
npm install
```

2. Start the development server:

```bash
npm run dev
```

3. Open [http://localhost:3000](http://localhost:3000) in your browser.

## Project Structure

```
{{project_name}}/
├── app/
│   ├── layout.tsx        # Root layout with metadata
│   ├── page.tsx          # Home page
│   └── globals.css       # Global styles
├── components/           # Reusable React components
├── public/               # Static assets served from /
├── next.config.ts        # Next.js configuration
├── package.json
├── tsconfig.json
└── README.md
```

## Scripts

| Script | Description |
|--------|-------------|
| `npm run dev`   | Start the development server with hot reload |
| `npm run build` | Create an optimized production build |
| `npm run start` | Start the production server |
| `npm run lint`  | Lint the codebase with ESLint |

## Development

Add pages under `app/`, components under `components/`, and static assets under `public/`. The development server hot-reloads on save.

## License

MIT — {{author}}
