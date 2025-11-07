# Sports Game Predictor

**Predict the most exciting games of the day across the four major sports leagues using Go.**

This project is a personal project for software engineers interested in sports analytics. It combines real-time sports data, player/team statistics, rivalry heuristics, and other metrics to rank games by predicted excitement. The goal is to help users decide which games are the “must-watch” for a given day.

---

## Table of Contents

* [Project Overview](#project-overview)
* [Features](#features)
* [Architecture](#architecture)
* [Data Sources](#data-sources)
* [Scoring Algorithm](#scoring-algorithm)
* [Installation](#installation)
* [Usage](#usage)
* [Configuration](#configuration)
* [Extending the Project](#extending-the-project)
* [Contributing](#contributing)
* [License](#license)

---

## Project Overview

Sports Game Predictor fetches daily schedules from multiple sports leagues, calculates a composite “excitement score” for each game, and exposes an API and/or CLI to display the top games for a user-selected day.

The project aims to demonstrate the following skills:

* Building a concurrent Go application that fetches and processes external APIs.
* Using structured Go code with clear domain models.
* Implementing scoring algorithms with weighted factors and normalization.
* Optionally integrating a small frontend or terminal UI for interactive display.
* Practicing caching, persistence, and API design in Go.

---

## Features

* Fetch schedules and basic stats from multiple leagues (NBA, NFL, MLB, NHL).
* Compute excitement scores using normalized metrics:

  * Game competitiveness
  * Team stakes / playoff relevance
  * Rivalries
  * Star player presence
  * Odds movement
  * Recent team form
  * Viewing time convenience
* Expose top games for a given date via REST API or CLI.
* Concurrent data fetching with rate-limiting to respect API limits.
* Configurable weighting for scoring factors.

---

## Architecture

### Core Components

1. **Ingest Workers**

   * Fetch data from league APIs.
   * Handle concurrency and rate limiting.

2. **Processor**

   * Computes factor values for each game.
   * Normalizes data and calculates excitement scores.

3. **Database / Storage**

   * PostgreSQL or SQLite for storing historical games and scores.
   * Redis cache for hot results.

4. **API Layer**

   * Exposes REST endpoints (e.g., `/v1/games?date=YYYY-MM-DD&top=N`).

5. **Optional Frontend**

   * Web UI (React/Vite) or CLI display.

---

## Data Sources

**Free APIs**

* NBA: [balldontlie](https://www.balldontlie.io/)
* NFL: community or official schedule APIs
* MLB: [MLB Stats API](https://statsapi.mlb.com/)
* NHL: Unofficial NHL community API

**Paid / Optional APIs**

* [MySportsFeeds](https://www.mysportsfeeds.com/)
* [SportsDataIO](https://sportsdata.io/)
* [Sportradar](https://sportradar.com/)

> Note: Odds and real-time data often require paid feeds or licensed access. Always respect API usage policies.

---

## Scoring Algorithm

Each game receives a composite **excitement score** using weighted factors normalized between 0 and 1.

### Factors

* **Closeness (30%)** – How balanced the game is (expected win probability near 50%).
* **Team Stakes (20%)** – Playoff relevance or team rankings.
* **Rivalry (15%)** – Historical rivalries or divisional matchups.
* **Star Power (15%)** – Presence of marquee players.
* **Line Movement (10%)** – Betting odds or market movement.
* **Form Difference (5%)** – Recent team performance.
* **Time Boost (5%)** – Viewer-friendly game start times.

### Formula

```text
Score = 0.3*Closeness + 0.2*Stakes + 0.15*Rivalry + 0.15*Stars + 0.1*LineMove + 0.05*FormDiff + 0.05*TimeBoost
```

---

## Installation

1. **Clone the repository**:

```bash
git clone https://github.com/username/sports-game-predictor.git
cd sports-game-predictor
```

2. **Install Go dependencies**:

```bash
go mod tidy
```

3. **Set up database**:

* PostgreSQL or SQLite.
* Apply migrations if using PostgreSQL.

4. **Set API keys (if required)**:

```bash
export NBA_API_KEY="your_key_here"
export MLB_API_KEY="your_key_here"
```

---

## Usage

### Run the CLI

```bash
go run ./cmd/main.go --date 2025-11-06 --league NBA --top 5
```

### Start API Server

```bash
go run ./cmd/server.go
```

API Endpoint:

```
GET /v1/games?date=2025-11-06&league=NBA&top=5
```

Response:

```json
[
  {
    "id": "20251106-NBA-LAL-BOS",
    "home_team": "Los Angeles Lakers",
    "away_team": "Boston Celtics",
    "excitation_score": 0.87,
    "factors": {
      "closeness": 0.95,
      "stakes": 0.8,
      "rivalry": 0.9,
      "star_power": 0.85,
      "line_move": 0.1,
      "form_diff": 0.4,
      "time_boost": 0.6
    }
  }
]
```

---

## Configuration

* Configurable scoring weights via YAML / JSON file.
* Ability to select leagues, date range, and top-N games.
* API endpoint and port settings.

---

## Extending the Project

* Add more leagues or international competitions.
* Integrate machine learning to learn user preferences and optimize weights.
* Add live updates and push notifications for game start reminders.
* Add richer stats (player-level data, injuries, advanced metrics).

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-league`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-league`)
5. Open a Pull Request

---

## License

MIT License © 2025 Your Name

> Enjoy predicting the most exciting games and happy coding!
