# Best Game of the Day

**A simple Go CLI tool that finds the most exciting upcoming games from NBA and NFL.**

This tool fetches games for the upcoming week (tomorrow through 7 days ahead) from NBA and NFL, calculates a simple excitement score based on team ranking proximity, and displays the top games.

---

## How It Works

1. Fetches NBA and NFL games for the upcoming week (excluding today)
2. Calculates an excitement score for each game using a rank-proxy closeness heuristic
3. Sorts games by excitement score and displays the top 10

The excitement score is based on how close the teams' rank proxies are - games between similarly ranked teams are considered more exciting.

---

## Installation

1. **Clone the repository**:

```bash
git clone https://github.com/username/best-game-of-the-day.git
cd best-game-of-the-day
```

2. **Install Go dependencies**:

```bash
go mod tidy
```

3. **Set API key**:

```bash
export BALLEDONTLIE_API_KEY="your_key_here"
```

Get your API key from [balldontlie.io](https://www.balldontlie.io/)

---

## Usage

Simply run:

```bash
go run main.go
```

The tool will:
- Fetch NBA and NFL games for the upcoming week
- Calculate excitement scores
- Display the top 10 games

Example output:

```
Fetching NBA and NFL games from 2025-01-15 to 2025-01-22 (excluding today)

Top 10 games for upcoming week based on rank-proxy closeness heuristic:

1. [NBA] Los Angeles Lakers (Rank Proxy 15) vs Boston Celtics (Rank Proxy 16)
   Scores: Away 0 - Home 0 | Excitement: 0.94

2. [NFL] Kansas City Chiefs (Rank Proxy 30) vs Buffalo Bills (Rank Proxy 31)
   Scores: Away 0 - Home 0 | Excitement: 0.94
...
```

---

## Data Sources

* **NBA**: [balldontlie API](https://www.balldontlie.io/)
* **NFL**: [balldontlie NFL API](https://www.balldontlie.io/)

---

## Scoring Algorithm

The excitement score is calculated using a simple closeness heuristic:

```go
closenessScore = 1 - |homeRank - visitorRank| / maxRankDiff
```

Where:
- `homeRank` and `visitorRank` are team ID-based rank proxies
- `maxRankDiff` is 16 for NBA and 32 for NFL

Games with teams that have similar rank proxies receive higher excitement scores.

---

## License

MIT License © 2025
