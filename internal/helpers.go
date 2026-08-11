package riot

import (
	"encoding/json"
	"log"
	"sync"
)

// Helper to flatten raw events across all timeline frames
func ExtractRawEvents(timeline TimelineResponse) []json.RawMessage {
	var rawEvents []json.RawMessage
	for _, frame := range timeline.Info.Frames {
		rawEvents = append(rawEvents, frame.Events...)
	}
	return rawEvents
}

func EnrichEvents(events []EventWrapper, participantMap map[int]Participant) []map[string]any {
	var enriched []map[string]any

	for _, event := range events {
		entry := map[string]any{
			"type": event.Type,
		}

		switch e := event.Data.(type) {
		case ItemPurchasedEvent:
			entry["timestamp"] = e.Timestamp
			entry["itemId"] = e.ItemID
			entry["participantId"] = e.ParticipantID

			if player, exists := participantMap[e.ParticipantID]; exists {
				entry["championName"] = player.ChampionName
				entry["riotIdName"] = player.RiotIDName
				entry["playerTier"] = player.Tier
				entry["playerRank"] = player.Rank
				entry["playerWinrate"] = player.Winrate
			}

		case ChampionKillEvent:
			entry["timestamp"] = e.Timestamp
			entry["killerId"] = e.KillerID
			entry["victimId"] = e.VictimID
			entry["assists"] = e.Assists

			if killer, exists := participantMap[e.KillerID]; exists {
				entry["killerChampion"] = killer.ChampionName
				entry["killerRiotId"] = killer.RiotIDName
			}
			if victim, exists := participantMap[e.VictimID]; exists {
				entry["victimChampion"] = victim.ChampionName
				entry["victimRiotId"] = victim.RiotIDName
			}

		default:
			// Fallback for unparsed generic events
			if generic, ok := event.Data.(map[string]any); ok {
				for k, v := range generic {
					entry[k] = v
				}
			}
		}

		enriched = append(enriched, entry)
	}

	return enriched
}

func BuildParticipantMap(participants []Participant, targetPuuid string) (map[int]Participant, int) {
	pMap := make(map[int]Participant)
	var targetID int

	for _, p := range participants {
		pMap[p.ParticipantID] = p
		if p.Puuid == targetPuuid {
			targetID = p.ParticipantID
		}
	}

	return pMap, targetID
}

func EnrichParticipantsWithRanks(pMap map[int]Participant, fetchRankFunc func(summonerID string) (string, string, int, int, float64, error)) {
	var wg sync.WaitGroup

	for id, player := range pMap {
		if player.Puuid == "" {
			continue
		}

		wg.Add(1)
		go func(pID int, puuid string) {
			defer wg.Done()

			tier, rank, wins, losses, winrate, err := fetchRankFunc(puuid)
			if err != nil {
				log.Printf("Warning: rank fetch failed for puuid %s: %v", puuid, err)
				return
			}
			p := pMap[pID]
			p.Tier = tier
			p.Rank = rank
			p.Wins = wins
			p.Losses = losses
			p.Winrate = winrate
			pMap[pID] = p // Update map entry
		}(id, player.Puuid)
	}

	wg.Wait()
}
