package main

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

// Event structs
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
	Type string `json:"type"`
	Data any    `json:"data"`
}

func parseTimelineEvents(rawEvents []json.RawMessage) ([]EventWrapper, error) {
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
