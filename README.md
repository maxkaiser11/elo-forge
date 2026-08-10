# League AI Coaching API

A Go-based backend service designed to pull match timeline data from the Riot Games API, unmarshal dynamic timeline events, and prepare structured game state payloads for AI coaching algorithms.

## 📁 Project Structure

```text
.
├── cmd/
│   └── api/             # Application entry point & setup
│       ├── main.go
│       └── config.go    # Environment configuration loader
├── internal/
│   └── riot/            # Riot API types & timeline event parser
│       └── timeline.go
├── .env.example         # Template for environment variables
├── .gitignore
├── README.md
├── go.mod
└── go.sum
```

---

## ⚡ Quick Start

### 1. Prerequisites
* **Go 1.20+** installed
* A valid **Riot API Key** from the [Riot Developer Portal](https://developer.riotgames.com/)

### 2. Environment Setup
Create a `.env` file in the root directory:

```env
RIOT_API_KEY=RGAPI-your-actual-api-key-here
RIOT_REGION=europe
GAME_NAME=YourRiotName
TAGLINE=YourTagline
```

### 3. Install Dependencies & Run
```bash
# Tidy dependencies
go mod tidy

# Run the API service
go run ./cmd/api
```

---

## 🛠 Features

* **Generic HTTP Client:** Type-safe API fetching using Go generics (`[T any]`).
* **Dynamic Event Parser:** Uses a two-pass `json.RawMessage` unmarshaler to safely handle varying Riot timeline event types (`ITEM_PURCHASED`, `CHAMPION_KILL`, etc.) without data loss.
* **Formatted Data Export:** Automatically marshals raw timeline responses and parsed typed events into formatted `.json` files for local inspection.
