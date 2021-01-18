package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/vysiondev/qstats-go/utils"
	"strconv"
)

func (b *BaseHandler) ExecuteTopCommand(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, rd RunDetails) error {
	user, err := GetUserData(rd.QuaverID)
	if err != nil {
		return err
	}
	scores, err := GetScores(rd.QuaverID, rd.Is7K, true, rd.Page, rd.PageGiven)
	if err != nil {
		return err
	}

	fields, err := CreateScoreFields(scores.Scores, rd.Page, false)

	mapIDs := make([]int, len(scores.Scores))
	for i, s := range scores.Scores {
		mapIDs[i] = int(s.Map.ID)
	}

	e := utils.CreateEmbed()
	if rd.PageGiven {
		e.AddTitle(fmt.Sprintf(":flag_%s: #%d %s top score for %s",
			utils.GetCountryStr(user.User.Info.Country),
			rd.Page,
			utils.GetKeymodeString(rd.Is7K),
			user.User.Info.Username,
		))
	} else {
		e.AddTitle(fmt.Sprintf(":flag_%s: Top 5 %s scores for %s",
			utils.GetCountryStr(user.User.Info.Country),
			utils.GetKeymodeString(rd.Is7K),
			user.User.Info.Username,
		))
	}
	e.AddTitleURL("https://quavergame.com/user/" + strconv.Itoa(rd.QuaverID))
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
