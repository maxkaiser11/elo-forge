package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Account struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagline"`
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	escapedName := url.PathEscape(cfg.GameName)
	escapedTag := url.PathEscape(cfg.TagLine)

	// 1. Fetch Account
	accountUrl := fmt.Sprintf("https://%s.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s", cfg.Region, escapedName, escapedTag)
	account, err := fetchRiotData[Account](accountUrl, cfg.RiotAPIKey)
	if err != nil {
		log.Fatalf("Failed to fetch account: %v", err)
	}

	// 2. Fetch Match IDs
	matchUrl := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids", cfg.Region, account.Puuid)
	matches, err := fetchRiotData[[]string](matchUrl, cfg.RiotAPIKey)
	if err != nil {
		log.Fatalf("Failed to fetch matches: %v", err)
	}
	if len(matches) == 0 {
		log.Fatal("No matches found for this user.")
	}

	lastMatch := matches[0]
	fmt.Println("Latest match ID:", lastMatch)

	// 3. Fetch Timeline
	timelineUrl := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s/timeline", cfg.Region, lastMatch)
	timeline, err := fetchRiotData[TimelineResponse](timelineUrl, cfg.RiotAPIKey)
	if err != nil {
		log.Fatalf("Failed to fetch timeline: %v", err)
	}

	// Extract raw events across all frames
	var rawEvents []json.RawMessage
	for _, frame := range timeline.Info.Frames {
		rawEvents = append(rawEvents, frame.Events...)
	}

	// Parse dynamic event structs (from timeline.go)
	events, err := parseTimelineEvents(rawEvents)
	if err != nil {
		log.Fatalf("Failed to parse events: %v", err)
	}

	parsedFilename := fmt.Sprintf("parsed_events_%s.json", lastMatch)
	if err := saveToFile(parsedFilename, events); err != nil {
		log.Printf("Warning: Failed to save parsed events file: %w", err)
	}

	for _, event := range events {
		switch e := event.Data.(type) {
		case ItemPurchasedEvent:
			fmt.Printf("[Item] Player %d bought item %d at %dms\n", e.ParticipantID, e.ItemID, e.Timestamp)
		case ChampionKillEvent:
			fmt.Printf("[Kill] Player %d killed Player %d with assists: %v\n", e.KillerID, e.VictimID, e.Assists)
		}
	}
}

func saveToFile(filename string, data any) error {
	prettyJSON, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	err = os.WriteFile(filename, prettyJSON, 0o644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", filename, err)
	}

	fmt.Printf("Saved: %s\n", filename)
	return nil
}

// Universal Generic Helper for Riot API GET requests
func fetchRiotData[T any](requestURL, apiKey string) (T, error) {
	var result T

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return result, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Riot-Token", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("Riot API error (status %d): %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("unmarshaling JSON: %w", err)
	}

	return result, nil
}
