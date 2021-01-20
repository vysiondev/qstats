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

func (b *BaseHandler) GetLetterGrade(letterGradeStr string) string {
	switch letterGradeStr {
	case "F":
		return b.Config.Emoji.Grade.F
	case "D":
		return b.Config.Emoji.Grade.D
	case "C":
		return b.Config.Emoji.Grade.C
	case "B":
		return b.Config.Emoji.Grade.B
	case "A":
		return b.Config.Emoji.Grade.A
	case "S":
		return b.Config.Emoji.Grade.S
	case "SS":
		return b.Config.Emoji.Grade.SS
	case "X":
		return b.Config.Emoji.Grade.X
	}
	return "?"
}
func (b *BaseHandler) GetRankedStatus(status int) string {
	switch status {
	case 0:
		return b.Config.Emoji.RankedStatus.NotSubmitted
	case 1:
		return b.Config.Emoji.RankedStatus.Unranked
	case 2:
		return b.Config.Emoji.RankedStatus.Ranked
	case 3:
		return b.Config.Emoji.RankedStatus.Dan
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

func (b *BaseHandler) calculateScoreField(score quaverapi_structs.Score, c chan<- MessageEmbedFieldChan, i int, mapData quaverapi_structs.QuaverMap, numberOverrideByPage int, hideNumbers bool, onExit func()) {
	go func() {
		defer onExit()
		maxCombo := (mapData.CountHitobjectLong * 2) + mapData.CountHitobjectNormal
		var numStr string
		if hideNumbers {
			numStr = ""
		} else {
			numStr = "#**" + strconv.Itoa(numberOverrideByPage+i) + "** • "
		}
		rankedStatus := b.GetRankedStatus(int(mapData.RankedStatus))
		fullCombo := ""
		if score.MaxCombo >= maxCombo && maxCombo != 0 {
			fullCombo = " • " + b.Config.Emoji.FullCombo
		}
		diffRating := "?"
		if mapData.DifficultyRating != 0.0 {
			diffRating = fmt.Sprintf("%0.2f", mapData.DifficultyRating)
		}
		replayAvailable := b.Config.Emoji.Download + " ---"
		if score.PersonalBest {
			replayAvailable = fmt.Sprintf("[%s Replay](%s/download/replay/%d)", b.Config.Emoji.Download, bot_constants.QuaverMainSite, score.ID)
		}
		t, err := time.Parse(time.RFC3339, score.Time)
		if err != nil {
			t = time.Now()
		}
		c <- MessageEmbedFieldChan{
			Field: discordgo.MessageEmbedField{
				Name: numStr + rankedStatus + fullCombo,
				Value: fmt.Sprintf("**[%s [%s]](%s/mapset/map/%d)** [%s] +**%s**\n%s **%0.2f** QR • %0.2f%%\n%0.2f ratio • x%s / %s • [%s / %s / %s / %s / %s / %s]\nSet %s • %s",
					utils.RemoveFormattingCharacters(mapData.Title),
					utils.RemoveFormattingCharacters(mapData.DifficultyName),
					bot_constants.QuaverMainSite,
					mapData.ID,
					diffRating,
					score.ModsString,
					b.GetLetterGrade(score.Grade),
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

func (b *BaseHandler) CreateScoreFields(scores []quaverapi_structs.Score, numberOverrideByPage int, hideNumbers bool) ([]discordgo.MessageEmbedField, error) {
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
		b.calculateScoreField(scores[i], fieldChan, i, mapData.Map, numberOverrideByPage, hideNumbers, func() { wg.Done() })
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
