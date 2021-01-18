package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/vysiondev/qstats-go/quaverapi_structs"
	"github.com/vysiondev/qstats-go/utils"
	"strings"
)

func calculateJudgementRatio(marv int, perf int, great int, good int, okay int, miss int) [6]float32 {
	var judgePercentagesOfWhole [6]float32
	for i, judge := range [6]int{marv, perf, great, good, okay, miss} {
		if judge == 0 {
			judgePercentagesOfWhole[i] = 0.0
		} else {
			tot := marv + perf + great + good + okay + miss
			if tot == 0 {
				judgePercentagesOfWhole[i] = 0.0
				continue
			}
			judgePercentagesOfWhole[i] = (float32(judge) / float32(tot)) * 100
		}
	}
	return judgePercentagesOfWhole
}

func (b *BaseHandler) ExecuteProfileCommand(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, d RunDetails) error {
	data, err := GetUserData(d.QuaverID)
	if err != nil {
		return err
	}
	var keymodeData *quaverapi_structs.KeymodeStats
	if d.Is7K {
		keymodeData = &data.User.Keys7
	} else {
		keymodeData = &data.User.Keys4
	}
	userOnlineStatus, err := GetUserSpecificOnlineStatus(d.QuaverID)
	if err != nil {
		return err
	}
	judgementRatios := calculateJudgementRatio(
		keymodeData.Stats.TotalMarv,
		keymodeData.Stats.TotalPerf,
		keymodeData.Stats.TotalGreat,
		keymodeData.Stats.TotalGood,
		keymodeData.Stats.TotalOkay,
		keymodeData.Stats.TotalMiss,
	)

	e := utils.CreateEmbed()
	kms := utils.GetKeymodeString(d.Is7K)
	e.AddTitle(fmt.Sprintf(":flag_%s: %s (%s)", strings.ToLower(utils.GetCountryStr(data.User.Info.Country)), data.User.Info.Username, kms))
	e.AddDescription(fmt.Sprintf("%s\n**[Steam](https://steamcommunity.com/profile/%s)** • **[Quaver Profile](https://quavergame.com/user/%d)**\n\nGlobal %s ranking: #**%s** (Country: **#%s**)\nPerformance rating (QR): **%0.2f**\nAccuracy: **%0.2f%%**\nJudgements: M:**%0.2f%%**, P:**%0.2f%%**, Gr:**%0.2f%%**, Go:**%0.2f%%**, O:**%0.2f%%**, M:**%0.2f%%**\n\nPlay count: **%d** (Failed: %d)\nTotal score: **%s**\nRanked score: **%s**",
		userOnlineStatus,
		data.User.Info.SteamID,
		data.User.Info.ID,
		kms,
		utils.AddCommas(int64(keymodeData.GlobalRank)),
		utils.AddCommas(int64(keymodeData.CountryRank)),
		keymodeData.Stats.OverallPerformanceRating,
		keymodeData.Stats.OverallAccuracy,
		judgementRatios[0],
		judgementRatios[1],
		judgementRatios[2],
		judgementRatios[3],
		judgementRatios[4],
		judgementRatios[5],
		keymodeData.Stats.PlayCount,
		keymodeData.Stats.FailCount,
		utils.AddCommas(int64(keymodeData.Stats.TotalScore)),
		utils.AddCommas(int64(keymodeData.Stats.RankedScore)),
	))
	e.AddThumbnail(data.User.Info.AvatarURL)

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
	if err != nil {
		fmt.Println("Failed to send message:", err)
	}
	return nil
}
