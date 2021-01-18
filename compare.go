package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/vysiondev/qstats-go/quaverapi_structs"
	"github.com/vysiondev/qstats-go/utils"
	"strconv"
)

func (b *BaseHandler) ExecuteCompareCommand(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, rd RunDetails) error {
	mapsInConversationData, err := b.GetMapsInConversation(m.GuildID, ctx)
	if err != nil {
		return err
	}

	score, scoreErr := FindFirstScoreInMapIDArray(rd.QuaverID, mapsInConversationData.Maps, rd.Is7K)
	if scoreErr != nil {
		return scoreErr
	}

	scoreArray := []quaverapi_structs.Score{score.Score}
	embedFields, fieldErr := CreateScoreFields(scoreArray, rd.Page, true)
	if fieldErr != nil {
		return err
	}

	e := utils.CreateEmbed()
	e.AddTitle(fmt.Sprintf(":flag_%s: %s's %s score comparison (#%d top score)",
		utils.GetCountryStr(score.UserData.User.Info.Country),
		score.UserData.User.Info.Username,
		utils.GetKeymodeString(rd.Is7K),
		score.BestScoreIndex+1,
	))
	e.AddTitleURL("https://quavergame.com/user/" + strconv.Itoa(rd.QuaverID))
	e.AddField(embedFields[0].Name, embedFields[0].Value)
	e.AddThumbnail(score.UserData.User.Info.AvatarURL)
	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
	if sendErr != nil {
		return sendErr
	}
	return nil
}
