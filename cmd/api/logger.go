package main

import (
	"fmt"

	"github.com/maxkaiser11/elo-forge/internal"
)

func logSearchedPlayerHighlights(events []riot.EventWrapper, pMap map[int]riot.Participant, searchedID int) {
	for _, event := range events {
		if kill, ok := event.Data.(riot.ChampionKillEvent); ok {
			killer := pMap[kill.KillerID]
			victim := pMap[kill.VictimID]

			if kill.KillerID == searchedID {
				fmt.Printf("🎯 [KILL] You (%s) killed %s (%s) at %dms\n",
					killer.ChampionName, victim.ChampionName, victim.RiotIDName, kill.Timestamp)
			} else if kill.VictimID == searchedID {
				fmt.Printf("⚠️ [DEATH] You (%s) were killed by %s (%s) at %dms\n",
					victim.ChampionName, killer.ChampionName, killer.RiotIDName, kill.Timestamp)
			}
		}
	}
}
