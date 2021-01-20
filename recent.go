package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/vysiondev/qstats-go/bot_constants"
	"github.com/vysiondev/qstats-go/utils"
	"strconv"
)

func (b *BaseHandler) ExecuteRecentCommand(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, rd RunDetails) error {
	user, err := GetUserData(rd.QuaverID)
	if err != nil {
		return err
	}
	scores, err := GetScores(rd.QuaverID, rd.Is7K, false, rd.Page, true)
	if err != nil {
		return err
	}

	fields, err := b.CreateScoreFields(scores.Scores, rd.Page, false)

	mapIDs := make([]int, len(scores.Scores))
	for i, s := range scores.Scores {
		mapIDs[i] = int(s.ID)
	}

	headerStr := fmt.Sprintf(":flag_%s: #%d most recent %s score for %s",
		utils.GetCountryStr(user.User.Info.Country),
		rd.Page,
		utils.GetKeymodeString(rd.Is7K),
		user.User.Info.Username,
	)
	if !rd.PageGiven {
		headerStr = fmt.Sprintf(":flag_%s: %s's most recent score for %s",
			utils.GetCountryStr(user.User.Info.Country),
			user.User.Info.Username,
			utils.GetKeymodeString(rd.Is7K),
		)
	}
	e := utils.CreateEmbed()
	e.AddTitle(headerStr)
	e.AddTitleURL(bot_constants.QuaverMainSite + "/user/" + strconv.Itoa(rd.QuaverID))
	e.AddThumbnail(user.User.Info.AvatarURL)
	for _, f := range fields {
		e.AddField(f.Name, f.Value)
	}

	_, sendErr := s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
	if sendErr != nil {
		return sendErr
	}
	micErr := b.SetMapsInConversation(m.GuildID, rd.Is7K, mapIDs, ctx)
	if micErr != nil {
		return err
	}
	return nil
}
