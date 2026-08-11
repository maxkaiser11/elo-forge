package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	riot "github.com/maxkaiser11/elo-forge/internal"
)

type Account struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagline"`
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// 1. Fetch Data
	account, err := fetchAccount(cfg)
	if err != nil {
		log.Fatalf("Account fetch failed: %v", err)
	}

	lastMatchID, err := fetchLatestMatchID(cfg, account.Puuid)
	if err != nil {
		log.Fatalf("Match ID fetch failed: %v", err)
	}

	matchDetails, err := fetchMatchDetails(cfg, lastMatchID)
	if err != nil {
		log.Fatalf("Match details fetch failed: %v", err)
	}

	// 2. Map Participants
	participantMap, searchedPlayerID := riot.BuildParticipantMap(matchDetails.Info.Participants, account.Puuid)
	riot.EnrichParticipantsWithRanks(participantMap, func(puuid string) (string, string, int, int, float64, error) {
		return fetchPlayerRank(cfg, puuid)
	})
	if player, ok := participantMap[searchedPlayerID]; ok {
		fmt.Printf("Searched Player: #%d %s (%s) - %s | Rank: %s %s (%.1f%% WR)\n",
			searchedPlayerID, player.ChampionName, player.RiotIDName, player.Role,
			player.Tier, player.Rank, player.Winrate)
	}

	// 3. Process Timeline & Enrich
	timeline, err := fetchTimeline(cfg, lastMatchID)
	if err != nil {
		log.Fatalf("Timeline fetch failed: %v", err)
	}

	events, err := riot.ParseTimelineEvents(riot.ExtractRawEvents(timeline))
	if err != nil {
		log.Fatalf("Timeline parsing failed: %v", err)
	}

	logSearchedPlayerHighlights(events, participantMap, searchedPlayerID)

	// 4. Export JSON Output
	enrichedEvents := riot.EnrichEvents(events, participantMap)
	filename := fmt.Sprintf("parsed_events_%s.json", lastMatchID)
	if err := saveToFile(filename, enrichedEvents); err != nil {
		log.Printf("Warning: Failed to save JSON file: %v", err)
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
