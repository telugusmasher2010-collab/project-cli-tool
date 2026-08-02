# {{project_name}}

> Map-first mobile app built with React Native, Expo, Expo Router, and react-native-maps

[![Expo](https://img.shields.io/badge/Expo-SDK%2054-black)](https://expo.dev)
[![React Native](https://img.shields.io/badge/React%20Native-0.81-blue)](https://reactnative.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue.svg)](https://www.typescriptlang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Overview

{{project_name}} is a cross-platform mobile app built with React Native and Expo that puts interactive maps at the center of the experience. It uses Expo Router for file-based navigation and react-native-maps for native map rendering on iOS and Android.

## Features

- Cross-platform iOS and Android support via Expo
- File-based routing with Expo Router
- Native map rendering with react-native-maps
- TypeScript with strict mode out of the box
- Runs in Expo Go — no native build required to get started

## Tech Stack

| Layer | Technology |
|-------|------------|
| Framework | Expo SDK 54 |
| UI | React Native 0.81 |
| Language | TypeScript |
| Navigation | Expo Router |
| Maps | react-native-maps |

## Installation

### Prerequisites

- Node.js 20.19+ (as required by Expo SDK 54)
- npm or your preferred package manager
- The [Expo Go](https://expo.dev/go) app on your device, or an Android/iOS simulator

### Setup

1. Install dependencies:

```bash
npm install
```

> For app store builds, configure your Google Maps API key on the `react-native-maps` config plugin in `app.json` (see the [react-native-maps docs](https://docs.expo.dev/versions/latest/sdk/map-view/)).

## Running

Start the Metro bundler:

```bash
npm start
```

Then press `a` to open the app on an Android emulator, `i` for the iOS simulator, or scan the QR code with the Expo Go app on your device.

## Folder Structure

```
{{project_name}}/
├── app/
│   ├── _layout.tsx       # Root layout (Expo Router)
│   └── index.tsx         # Home screen with the map
├── assets/               # Images, fonts, and other assets
├── components/           # Reusable React components
├── app.json              # Expo configuration
├── package.json
├── tsconfig.json
└── README.md
```

## Scripts

| Script | Description |
|--------|-------------|
| `npm start` | Start the Expo development server (Metro) |
| `npm run android` | Start the dev server and open on Android |
| `npm run ios` | Start the dev server and open on iOS |
| `npm run lint` | Lint the codebase with ESLint |

## License

MIT — {{author}}
