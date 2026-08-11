package riot

import (
	"encoding/json"
	"fmt"
)

// Root timeline structs
type TimelineResponse struct {
	Metadata struct {
		MatchID string `json:"matchId"`
	} `json:"metadata"`
	Info TimelineInfo `json:"info"`
}

type TimelineInfo struct {
	FrameInterval int64           `json:"frameInterval"`
	Frames        []TimelineFrame `json:"frames"`
}

type TimelineFrame struct {
	Events    []json.RawMessage `json:"events"`
	Timestamp int64             `json:"timestamp"`
}

type BaseEvent struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
}

type ItemPurchasedEvent struct {
	BaseEvent
	ItemID        int `json:"itemId"`
	ParticipantID int `json:"participantId"`
}

type ChampionKillEvent struct {
	BaseEvent
	KillerID int   `json:"killerId"`
	VictimID int   `json:"victimId"`
	Assists  []int `json:"assistingParticipantIds"`
}

type EventWrapper struct {
	Type          string `json:"type"`
	ParticipantID int    `json:"participantId,omitempty"`
	ChampionName  string `json:"championName,omitempty"`
	RiotIDName    string `json:"RiotIdName,omitempty"`
	Data          any    `json:"data"`
}

type MatchResponse struct {
	Info MatchInfo `json:"info"`
}

type MatchInfo struct {
	Participants []Participant `json:"participants"`
}

type Participant struct {
	ParticipantID int    `json:"participantId"`
	Puuid         string `json:"puuid"`
	SummonerID    string `json:"summonerId"`
	RiotIDName    string `json:"riotIdGameName"`
	RiotIDTag     string `json:"riotIdTagline"`
	ChampionName  string `json:"championName"`
	TeamID        int    `json:"teamId"`
	Role          string `json:"individualPosition"`

	Tier    string  `json:"tier,omitempty"` // e.g., "DIAMOND"
	Rank    string  `json:"rank,omitempty"` // e.g., "I"
	Wins    int     `json:"wins,omitempty"`
	Losses  int     `json:"losses,omitempty"`
	Winrate float64 `json:"winrate,omitempty"` // e.g., 54.2
}

type LeagueEntry struct {
	QueueType string `json:"queueType"` // RANKED_SOLO_5x5
	Tier      string `json:"tier"`
	Rank      string `json:"rank"`
	Wins      int    `json:"wins"`
	Losses    int    `json:"losses"`
}

func ParseTimelineEvents(rawEvents []json.RawMessage) ([]EventWrapper, error) {
	var parsedEvents []EventWrapper

	for _, raw := range rawEvents {
		var base BaseEvent
		if err := json.Unmarshal(raw, &base); err != nil {
			return nil, fmt.Errorf("failed to parse base event: %w", err)
		}

		var eventData any

		switch base.Type {
		case "ITEM_PURCHASED":
			var e ItemPurchasedEvent
			if err := json.Unmarshal(raw, &e); err != nil {
				return nil, err
			}
			eventData = e

		case "CHAMPION_KILL":
			var e ChampionKillEvent
			if err := json.Unmarshal(raw, &e); err != nil {
				return nil, err
			}
			eventData = e

		default:
			var generic map[string]any
			if err := json.Unmarshal(raw, &generic); err != nil {
				return nil, err
			}
			eventData = generic
		}

		parsedEvents = append(parsedEvents, EventWrapper{
			Type: base.Type,
			Data: eventData,
		})
	}

	return parsedEvents, nil
}
