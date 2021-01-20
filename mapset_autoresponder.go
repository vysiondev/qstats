package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/vysiondev/qstats-go/bot_constants"
	"github.com/vysiondev/qstats-go/utils"
	"sort"
	"strconv"
	"time"
)

func (b *BaseHandler) AutorespondToMapset(s *discordgo.Session, m *discordgo.Message, mapID int) error {
	thisMap, err := GetMapset(mapID)
	if err != nil {
		return err
	}

	var keymodeStr string
	for _, m := range thisMap.Mapset.Maps {
		if m.GameMode == 1 {
			keymodeStr += "4K"
			break
		}
	}
	for _, m := range thisMap.Mapset.Maps {
		if m.GameMode == 2 {
			if len(keymodeStr) > 0 {
				keymodeStr += "+"
			}
			keymodeStr += "7K"
			break
		}
	}

	sort.SliceStable(thisMap.Mapset.Maps, func(i, j int) bool {
		return thisMap.Mapset.Maps[i].DifficultyRating > thisMap.Mapset.Maps[j].DifficultyRating
	})

	var diffStr string
	for i, m := range thisMap.Mapset.Maps {
		mapIs7K := false
		if m.GameMode == 2 {
			mapIs7K = true
		}
		if i < 5 {
			diffStr += fmt.Sprintf("- [[%s] %0.2f - %s](%s/mapset/map/%d)\n",
				utils.GetKeymodeString(mapIs7K),
				m.DifficultyRating,
				utils.RemoveFormattingCharacters(m.DifficultyName),
				bot_constants.QuaverMainSite,
				m.ID,
			)
		} else {
			diffStr += fmt.Sprintf("*+%d more easier difficulties...*", len(thisMap.Mapset.Maps)-i)
			break
		}
	}

	e := utils.CreateEmbed()

	t1, err := time.Parse(time.RFC3339, thisMap.Mapset.DateSubmitted)
	if err != nil {
		t1 = time.Now()
	}
	t2, err := time.Parse(time.RFC3339, thisMap.Mapset.DateLastUpdated)
	if err != nil {
		t2 = time.Now()
	}
	d := time.Duration(thisMap.Mapset.Maps[0].Length) * time.Millisecond

	e.AddDescription(fmt.Sprintf("%s • %s • Mapped by [%s](%s/user/%d)\n**%0.2f** BPM • Duration: **%s**\n[%s Download this mapset](%s/download/mapset/%d)",
		b.GetRankedStatus(int(thisMap.Mapset.Maps[0].RankedStatus)),
		keymodeStr,
		thisMap.Mapset.CreatorUsername,
		bot_constants.QuaverMainSite,
		thisMap.Mapset.CreatorID,
		thisMap.Mapset.Maps[0].BPM,
		utils.FmtDuration(d),
		b.Config.Emoji.Download,
		bot_constants.QuaverMainSite,
		thisMap.Mapset.ID,
	))
	e.AddTitle(fmt.Sprintf("%s - %s",
		utils.RemoveFormattingCharacters(thisMap.Mapset.Artist),
		utils.RemoveFormattingCharacters(thisMap.Mapset.Title),
	))
	e.AddField("Difficulties", diffStr)
	e.AddTitleURL(bot_constants.QuaverMainSite + "/mapset/" + strconv.Itoa(int(thisMap.Mapset.ID)))
	e.AddFooter(fmt.Sprintf("Submitted %s • Last updated %s",
		utils.TimeElapsed(time.Now(), t1, false),
		utils.TimeElapsed(time.Now(), t2, false),
	))
	e.AddImage("https://cdn.quavergame.com/mapsets/" + strconv.Itoa(int(thisMap.Mapset.ID)) + ".jpg")

	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
	if sendErr != nil {
		return sendErr
	}
	return nil
}
