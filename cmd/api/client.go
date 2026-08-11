package main

import (
	"fmt"

	riot "github.com/maxkaiser11/elo-forge/internal"
)

func fetchAccount(cfg *Config) (Account, error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s", cfg.Region, cfg.GameName, cfg.Tagline)
	return fetchRiotData[Account](url, cfg.RiotAPIKey)
}

func fetchLatestMatchID(cfg *Config, puuid string) (string, error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids", cfg.Region, puuid)
	matches, err := fetchRiotData[[]string](url, cfg.RiotAPIKey)
	if err != nil || len(matches) == 0 {
		return "", fmt.Errorf("no matches found: %w", err)
	}
	return matches[0], nil
}

func fetchMatchDetails(cfg *Config, matchID string) (riot.MatchResponse, error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s", cfg.Region, matchID)
	return fetchRiotData[riot.MatchResponse](url, cfg.RiotAPIKey)
}

func fetchTimeline(cfg *Config, matchID string) (riot.TimelineResponse, error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s/timeline", cfg.Region, matchID)
	return fetchRiotData[riot.TimelineResponse](url, cfg.RiotAPIKey)
}
