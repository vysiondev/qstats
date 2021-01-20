package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/vysiondev/qstats-go/bot_constants"
	"github.com/vysiondev/qstats-go/utils"
	"strconv"
	"strings"
)

func calcPlayerIndex(page int) int {
	return (page - 1) * 25
}

func (b *BaseHandler) ExecuteLeaderboardCommand(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, rd RunDetails) error {
	country := ""
	if len(rd.Args) > 0 {
		country = strings.ToLower(rd.Args[0])
	}
	lb, err := GetLeaderboard(rd.Page, country, rd.Is7K)
	if err != nil {
		return err
	}

	lbStr := ""
	for i, player := range lb {
		var playerQRAdvantage float64
		if i+1 != len(lb) {
			playerQRAdvantage = lb[i].Stats.OverallPerformanceRating - lb[i+1].Stats.OverallPerformanceRating
		} else {
			playerQRAdvantage = 0.0
		}

		playerCountry := utils.GetCountryStr(player.Country)

		numberOne := ""
		if player.Stats.Rank == 1 {
			numberOne = ":trophy: "
		}
		lbStr += fmt.Sprintf("%s:flag_%s: `%d.` %s • %0.2f QR • %0.2f%% (+%0.2f)\n",
			numberOne,
			strings.ToLower(playerCountry),
			calcPlayerIndex(rd.Page)+i+1,
			player.Username,
			player.Stats.OverallPerformanceRating,
			player.Stats.OverallAccuracy,
			playerQRAdvantage,
		)
	}

	lbTypeStr := ":earth_americas: Global"
	if len(rd.Args) > 0 {
		lbTypeStr = fmt.Sprintf(":flag_%s: %s", strings.ToLower(utils.GetCountryStr(rd.Args[0])), strings.ToUpper(rd.Args[0]))
	}

	embed := utils.CreateEmbed()
	embed.AddTitle(fmt.Sprintf("%s %s leaderboards (#%d - #%d)", lbTypeStr, utils.GetKeymodeString(rd.Is7K), calcPlayerIndex(rd.Page)+1, calcPlayerIndex(rd.Page)+25))
	embed.AddDescription(lbStr)
	embed.AddFooter("+ indicates QR advantage")
	embed.AddTitleURL(fmt.Sprintf(bot_constants.QuaverMainSite+"/leaderboard/?mode=%s&page=%s&country=%s", utils.GetKeymodeIntAsStr(rd.Is7K), strconv.Itoa(rd.Page-1), country))

	_, e := s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
	if e != nil {
		return e
	}
	return nil
}
