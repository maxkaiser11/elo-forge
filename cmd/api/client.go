package main

import (
	"fmt"
	"math"
	"strings"

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

func getPlatformHost(region string) string {
	switch strings.ToLower(region) {
	case "europe":
		return "euw1" // Or "eun1" depending on your target server
	case "americas":
		return "na1"
	case "asia":
		return "kr"
	default:
		return region // Fallback if cfg.Region is already "euw1"
	}
}

func fetchPlayerRank(cfg *Config, puuid string) (tier, rank string, wins, losses int, winrate float64, err error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/league/v4/entries/by-puuid/%s", getPlatformHost(cfg.Region), puuid)
	entries, err := fetchRiotData[[]riot.LeagueEntry](url, cfg.RiotAPIKey)
	if err != nil {
		return "UNRANKED", "", 0, 0, 0.0, err
	}
	if len(entries) == 0 {
		return "UNRANKED", "", 0, 0, 0.0, nil
	}

	for _, e := range entries {
		if e.QueueType == "RANKED_SOLO_5x5" {
			totalGames := e.Wins + e.Losses
			wr := 0.0
			if totalGames > 0 {
				wr = math.Round((float64(e.Wins)/float64(totalGames))*1000) / 10
			}
			return e.Tier, e.Rank, e.Wins, e.Losses, wr, nil
		}
	}
	return "UNRANKED", "", 0, 0, 0.0, nil
}
