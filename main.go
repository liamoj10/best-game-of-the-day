package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"time"
)

// -------------------- Data Models --------------------

type Team struct {
	ID        int    `json:"id"`        // Alphabetical proxy
	Name      string `json:"full_name"` 
	RankProxy int    // Use as proxy for ranking heuristic
}

type Game struct {
	ID           int    `json:"id"`
	Date         string `json:"date"`
	HomeTeam     Team   `json:"home_team"`
	VisitorTeam  Team   `json:"visitor_team"`
	HomeScore    int    `json:"home_team_score"`
	VisitorScore int    `json:"visitor_team_score"`
	League       string // "NBA" or "NFL"
	Excitement   float64
}

type APIResponse struct {
	Data []Game      `json:"data"`
	Meta interface{} `json:"meta"`
}

// -------------------- Helper Functions --------------------

func closenessScore(homeRank, visitorRank, maxDiff int) float64 {
	diff := math.Abs(float64(homeRank - visitorRank))
	return 1 - diff/float64(maxDiff)
}

// -------------------- Fetch Functions --------------------

func fetchGames(url string, league string) ([]Game, error) {
	apiKey := os.Getenv("BALLEDONTLIE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("BALLEDONTLIE_API_KEY not set")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	for i := range apiResp.Data {
		game := &apiResp.Data[i]
		game.League = league

		// Use ID as ranking proxy
		game.HomeTeam.RankProxy = game.HomeTeam.ID
		game.VisitorTeam.RankProxy = game.VisitorTeam.ID
	}

	return apiResp.Data, nil
}

func fetchNBAGames(startDate, endDate string) ([]Game, error) {
	url := fmt.Sprintf("https://api.balldontlie.io/v1/games?start_date=%s&end_date=%s", startDate, endDate)
	return fetchGames(url, "NBA")
}

func fetchNFLGames(startDate, endDate string) ([]Game, error) {
	url := fmt.Sprintf("https://api.balldontlie.io/nfl/v1/games?dates[]=%s&dates[]=%s", startDate, endDate)
	return fetchGames(url, "NFL")
}

// -------------------- Main --------------------

func main() {
	today := time.Now()
	start := today.AddDate(0, 0, 1) // start from tomorrow
	end := today.AddDate(0, 0, 7)   // 7 days ahead

	startDate := start.Format("2006-01-02")
	endDate := end.Format("2006-01-02")

	fmt.Printf("Fetching NBA and NFL games from %s to %s (excluding today)\n", startDate, endDate)

	nbaGames, err := fetchNBAGames(startDate, endDate)
	if err != nil {
		fmt.Println("Error fetching NBA games:", err)
		nbaGames = []Game{}
	}

	nflGames, err := fetchNFLGames(startDate, endDate)
	if err != nil {
		fmt.Println("Error fetching NFL games:", err)
		nflGames = []Game{}
	}

	allGames := append(nbaGames, nflGames...)

	if len(allGames) == 0 {
		fmt.Println("No games found for the upcoming week.")
		return
	}

	// Compute rank-proxy closeness
	maxRank := map[string]int{
		"NBA": 16, // heuristic max difference
		"NFL": 32,
	}

	for i := range allGames {
		g := &allGames[i]

		homeRank := g.HomeTeam.RankProxy
		visitorRank := g.VisitorTeam.RankProxy

		// Fallback if missing
		if homeRank == 0 {
			homeRank = 1
		}
		if visitorRank == 0 {
			visitorRank = maxRank[g.League]
		}

		g.Excitement = closenessScore(homeRank, visitorRank, maxRank[g.League])
	}

	// Sort descending by excitement
	sort.Slice(allGames, func(i, j int) bool {
		return allGames[i].Excitement > allGames[j].Excitement
	})

	// Print top games
	topN := 10
	if len(allGames) < topN {
		topN = len(allGames)
	}

	fmt.Printf("\nTop %d games for upcoming week based on rank-proxy closeness heuristic:\n\n", topN)
	for i := 0; i < topN; i++ {
		g := allGames[i]
		fmt.Printf("%d. [%s] %s (Rank Proxy %d) vs %s (Rank Proxy %d)\n", i+1,
			g.League,
			g.VisitorTeam.Name, g.VisitorTeam.RankProxy,
			g.HomeTeam.Name, g.HomeTeam.RankProxy,
		)
		fmt.Printf("   Scores: Away %d - Home %d | Excitement: %.2f\n\n",
			g.VisitorScore, g.HomeScore, g.Excitement)
	}
}