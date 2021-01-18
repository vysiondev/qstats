package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/vysiondev/qstats-go/bot_constants"
	"github.com/vysiondev/qstats-go/quaverapi_structs"
	"github.com/vysiondev/qstats-go/utils"
	"strconv"
	"sync"
	"time"
)

func GetLetterGrade(letterGradeStr string) string {
	switch letterGradeStr {
	case "F":
		return bot_constants.GradeFEmoji
	case "D":
		return bot_constants.GradeDEmoji
	case "C":
		return bot_constants.GradeCEmoji
	case "B":
		return bot_constants.GradeBEmoji
	case "A":
		return bot_constants.GradeAEmoji
	case "S":
		return bot_constants.GradeSEmoji
	case "SS":
		return bot_constants.GradeSSEmoji
	case "X":
		return bot_constants.GradeXEmoji
	}
	return "?"
}
func GetRankedStatus(status int) string {
	switch status {
	case 0:
		return bot_constants.RsNotSubmitted
	case 1:
		return bot_constants.RsUnranked
	case 2:
		return bot_constants.RsRanked
	case 3:
		return bot_constants.RsDan
	}
	return "?"
}

func getScoreMapData(score quaverapi_structs.Score) (quaverapi_structs.QuaverMapResponse, error) {
	mapData, err := GetMap(int(score.Map.ID))
	if err != nil {
		return quaverapi_structs.QuaverMapResponse{}, err
	}
	return mapData, nil
}

func calculateScoreField(score quaverapi_structs.Score, c chan<- MessageEmbedFieldChan, i int, mapData quaverapi_structs.QuaverMap, numberOverrideByPage int, hideNumbers bool, onExit func()) {
	go func() {
		defer onExit()
		maxCombo := (mapData.CountHitobjectLong * 2) + mapData.CountHitobjectNormal
		var numStr string
		if hideNumbers {
			numStr = ""
		} else {
			numStr = "#**" + strconv.Itoa(numberOverrideByPage+i) + "** • "
		}
		rankedStatus := GetRankedStatus(int(mapData.RankedStatus))
		fullCombo := ""
		if score.MaxCombo >= maxCombo && maxCombo != 0 {
			fullCombo = " • " + bot_constants.FullCombo
		}
		diffRating := "?"
		if mapData.DifficultyRating != 0.0 {
			diffRating = fmt.Sprintf("%0.2f", mapData.DifficultyRating)
		}
		replayAvailable := bot_constants.ReplayEmoji + " ---"
		if score.PersonalBest {
			replayAvailable = fmt.Sprintf("[%s Replay](https://quavergame.com/download/replay/%d)", bot_constants.ReplayEmoji, score.ID)
		}
		t, err := time.Parse(time.RFC3339, score.Time)
		if err != nil {
			t = time.Now()
		}
		c <- MessageEmbedFieldChan{
			Field: discordgo.MessageEmbedField{
				Name: numStr + rankedStatus + fullCombo,
				Value: fmt.Sprintf("**[%s [%s]](https://quavergame.com/mapset/map/%d)** [%s] +**%s**\n%s **%0.2f** QR • %0.2f%%\n%0.2f ratio • x%s / %s • [%s / %s / %s / %s / %s / %s]\nSet %s • %s",
					mapData.Title,
					mapData.DifficultyName,
					mapData.ID,
					diffRating,
					score.ModsString,
					GetLetterGrade(score.Grade),
					score.PerformanceRating,
					score.Accuracy,
					score.Ratio,
					utils.AddCommas(score.MaxCombo),
					utils.AddCommas(maxCombo),
					utils.AddCommas(score.CountMarv),
					utils.AddCommas(score.CountPerf),
					utils.AddCommas(score.CountGreat),
					utils.AddCommas(score.CountGood),
					utils.AddCommas(score.CountOkay),
					utils.AddCommas(score.CountMiss),
					utils.TimeElapsed(time.Now(), t, false),
					replayAvailable,
				),
				Inline: false,
			},
			Index: i,
		}
	}()
}

type MessageEmbedFieldChan struct {
	Field discordgo.MessageEmbedField
	Index int
}

func CreateScoreFields(scores []quaverapi_structs.Score, numberOverrideByPage int, hideNumbers bool) ([]discordgo.MessageEmbedField, error) {
	var mapDataArray []quaverapi_structs.QuaverMapResponse
	fieldChan := make(chan MessageEmbedFieldChan, len(scores))

	for _, s := range scores {
		d, e := getScoreMapData(s)
		if e != nil {
			return nil, e
		}
		mapDataArray = append(mapDataArray, d)
	}

	var wg sync.WaitGroup
	// CPU intensive, so we spawn goroutines.
	for i, mapData := range mapDataArray {
		wg.Add(1)
		calculateScoreField(scores[i], fieldChan, i, mapData.Map, numberOverrideByPage, hideNumbers, func() { wg.Done() })
	}

	go func() {
		defer close(fieldChan)
		wg.Wait()
	}()

	fields := make([]discordgo.MessageEmbedField, len(scores))
	for field := range fieldChan {
		fields[field.Index] = field.Field
	}

	return fields, nil
}
