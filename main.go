package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"
)

// -------------------- Data Models --------------------

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"full_name"`
}

type Game struct {
	ID           int    `json:"id"`
	Date         string `json:"date"`
	HomeTeam     Team   `json:"home_team"`
	VisitorTeam  Team   `json:"visitor_team"`
	HomeScore    int    `json:"home_team_score"`
	VisitorScore int    `json:"visitor_team_score"`
	Excitement   float64
	Factors      GameFactors
}

type APIResponse struct {
	Data []Game      `json:"data"`
	Meta interface{} `json:"meta"`
}

// -------------------- Game Scoring --------------------

type GameFactors struct {
	PredictedHomeProb float64
	Stakes            float64
	Rivalry           float64
	StarPower         float64
	LineMove          float64
	FormDiff          float64
	TimeBoost         float64
}

func closeness(prob float64) float64 {
	if prob > 1 {
		prob = 1
	} else if prob < 0 {
		prob = 0
	}
	return 1 - abs(0.5-prob)*2
}

func computeScore(f GameFactors) float64 {
	weights := map[string]float64{
		"closeness": 0.30,
		"stakes":    0.20,
		"rivalry":   0.15,
		"stars":     0.15,
		"linemove":  0.10,
		"form":      0.05,
		"time":      0.05,
	}
	c := closeness(f.PredictedHomeProb)
	score := c*weights["closeness"] +
		f.Stakes*weights["stakes"] +
		f.Rivalry*weights["rivalry"] +
		f.StarPower*weights["stars"] +
		f.LineMove*weights["linemove"] +
		f.FormDiff*weights["form"] +
		f.TimeBoost*weights["time"]
	return score
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func randFloat() float64 {
	return float64(time.Now().UnixNano()%100) / 100.0
}

// -------------------- Fetch Games --------------------

func fetchGames(date string) ([]Game, error) {
	apiKey := os.Getenv("BALLEDONTLIE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("BALLEDONTLIE_API_KEY not set")
	}

	url := fmt.Sprintf("https://api.balldontlie.io/v1/games?dates[]=%s", date)
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

	return apiResp.Data, nil
}

// -------------------- Main --------------------

func main() {
	date := time.Now().Format("2006-01-02") // today
	fmt.Println("Fetching NBA games for", date)

	games, err := fetchGames(date)
	if err != nil {
		fmt.Println("Error fetching games:", err)
		return
	}

	if len(games) == 0 {
		fmt.Println("No games found for this date.")
		return
	}

	// Assign mock factors & compute excitement
	for i := range games {
		factors := GameFactors{
			PredictedHomeProb: 0.5 + 0.1*randFloat(),
			Stakes:            randFloat(),
			Rivalry:           randFloat(),
			StarPower:         randFloat(),
			LineMove:          randFloat(),
			FormDiff:          randFloat(),
			TimeBoost:         randFloat(),
		}
		games[i].Factors = factors
		games[i].Excitement = computeScore(factors)
	}

	// Sort by excitement descending
	sort.Slice(games, func(i, j int) bool {
		return games[i].Excitement > games[j].Excitement
	})

	// Print top games with verbose output
	topN := 5
	if len(games) < topN {
		topN = len(games)
	}

	fmt.Printf("\nTop %d games for %s:\n\n", topN, date)
	for i := 0; i < topN; i++ {
		g := games[i]
		fmt.Printf("%d. %s vs %s\n", i+1, g.VisitorTeam.Name, g.HomeTeam.Name)
		fmt.Printf("   Scores: Away %d - Home %d\n", g.VisitorScore, g.HomeScore)
		fmt.Printf("   Excitement Score: %.3f\n", g.Excitement)
		fmt.Println("   Factor Breakdown:")
		fmt.Printf("     Closeness: %.2f\n", closeness(g.Factors.PredictedHomeProb))
		fmt.Printf("     Stakes: %.2f\n", g.Factors.Stakes)
		fmt.Printf("     Rivalry: %.2f\n", g.Factors.Rivalry)
		fmt.Printf("     Star Power: %.2f\n", g.Factors.StarPower)
		fmt.Printf("     Line Move: %.2f\n", g.Factors.LineMove)
		fmt.Printf("     Form Diff: %.2f\n", g.Factors.FormDiff)
		fmt.Printf("     Time Boost: %.2f\n\n", g.Factors.TimeBoost)
	}
}