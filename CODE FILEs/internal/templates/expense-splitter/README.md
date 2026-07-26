# {{project_name}}

> Split expenses with friends. Built with Flutter, Supabase, and UPI.

[![Flutter](https://img.shields.io/badge/Flutter-3.x-blue.svg)](https://flutter.dev)
[![Dart](https://img.shields.io/badge/Dart-3.x-blue.svg)](https://dart.dev)
[![Supabase](https://img.shields.io/badge/Supabase-Database-green.svg)](https://supabase.com)

## Overview

{{project_name}} is a cross-platform mobile app for splitting expenses among groups. Add expenses, assign participants, and settle up via UPI — all backed by Supabase for real-time sync.

## Features

- Create and manage expense groups
- Add expenses with equal or custom splits
- Track who owes whom
- UPI payment integration for settlements
- Real-time sync via Supabase
- Clean Material Design 3 UI

## Getting Started

### Prerequisites

- Flutter 3.x SDK
- Dart 3.x
- A [Supabase](https://supabase.com) project

### Setup

1. Install dependencies:

```bash
flutter pub get
```

2. Configure environment:

```bash
cp .env.example .env
# Fill in your Supabase URL and anon key
```

3. Run database migrations from `supabase/schema.sql` in your Supabase dashboard.

4. Launch the app:

```bash
flutter run
```

## Project Structure

```
lib/
├── main.dart                  # Entry point
├── app.dart                   # MaterialApp setup, routing, theme
├── core/
│   └── constants.dart         # App-wide constants, API keys
├── features/
│   └── expenses/
│       ├── expense_page.dart  # Expense list screen
│       └── add_expense.dart   # Add expense form
└── ...

supabase/
└── schema.sql                 # Database schema
```

## Architecture

The app follows a feature-first folder structure:

- `core/` — shared utilities, constants, and themes
- `features/` — each feature is self-contained with its own UI, data, and logic

## License

MIT — {{author}}
