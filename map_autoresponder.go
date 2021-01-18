package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/vysiondev/qstats-go/utils"
	"strconv"
	"time"
)

func (b *BaseHandler) AutorespondToMap(s *discordgo.Session, m *discordgo.Message, mapID int, ctx context.Context) error {
	thisMap, err := GetMap(mapID)
	if err != nil {
		return err
	}

	thisMapIs7K := false
	if thisMap.Map.GameMode == 2 {
		thisMapIs7K = true
	}

	micErr := b.SetMapsInConversation(m.GuildID, thisMapIs7K, []int{int(thisMap.Map.ID)}, ctx)
	if micErr != nil {
		return err
	}

	e := utils.CreateEmbed()

	t1, err := time.Parse(time.RFC3339, thisMap.Map.DateSubmitted)
	if err != nil {
		t1 = time.Now()
	}
	t2, err := time.Parse(time.RFC3339, thisMap.Map.DateLastUpdated)
	if err != nil {
		t2 = time.Now()
	}
	d := time.Duration(thisMap.Map.Length) * time.Millisecond

	e.AddDescription(fmt.Sprintf("%s • Mapped by [%s](https://quavergame.com/user/%d)\nDifficulty: **%0.2f** QR • **%0.2f** BPM • Duration: **%s**\n**%0.2f**%% LN • **%sx** max combo • **%s** total plays\n[Download this mapset](https://quavergame.com/download/mapset/%d)",
		GetRankedStatus(int(thisMap.Map.RankedStatus)),
		thisMap.Map.CreatorUsername,
		thisMap.Map.CreatorID,
		thisMap.Map.DifficultyRating,
		thisMap.Map.BPM,
		utils.FmtDuration(d),
		(float64(thisMap.Map.CountHitobjectLong)/float64(thisMap.Map.CountHitobjectLong+thisMap.Map.CountHitobjectNormal))*100,
		utils.AddCommas((thisMap.Map.CountHitobjectLong*2)+thisMap.Map.CountHitobjectNormal),
		utils.AddCommas(thisMap.Map.PlayCount),
		thisMap.Map.ID,
	))
	e.AddTitle(fmt.Sprintf("[%s] %s - %s [%s]",
		utils.GetKeymodeString(thisMapIs7K),
		thisMap.Map.Artist,
		thisMap.Map.Title,
		thisMap.Map.DifficultyName,
	))
	e.AddTitleURL("https://quavergame.com/mapset/map/" + strconv.Itoa(int(thisMap.Map.ID)))
	e.AddFooter(fmt.Sprintf("Submitted %s • Last updated %s",
		utils.TimeElapsed(time.Now(), t1, false),
		utils.TimeElapsed(time.Now(), t2, false),
	))
	e.AddImage("https://cdn.quavergame.com/mapsets/" + strconv.Itoa(int(thisMap.Map.MapsetID)) + ".jpg")

	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
	if sendErr != nil {
		return sendErr
	}
	return nil
}
