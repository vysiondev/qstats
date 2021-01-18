package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/vysiondev/qstats-go/utils"
)

func (b *BaseHandler) ExecuteAboutCommand(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ RunDetails) error {
	em := utils.CreateEmbed()
	em.AddDescription(fmt.Sprintf("Hey there! I'm QStats, a Discord bot specially designed to fetch Quaver scores, map data, and more for your server's convenience.\n\nI'm currently maintained by [Vysion](https://vys.moe) with help on what to add from [yahweh](https://quavergame.com/user/383). You can contribute as well by visiting the [GitHub](%s) page!\n\n- [My website](%s)\n- [Support server](%s)\n- [Discord Bot List listing](https://top.gg/bot/%s)",
		b.Config.Links.Github,
		b.Config.Links.Website,
		b.Config.Links.SupportServer,
		s.State.User.ID,
	))
	em.AddTitle("About me")
	em.AddImage("https://i.imgur.com/zF34ZJH.png")
	em.AddFooter(fmt.Sprintf("Powered by discordgo %s", discordgo.VERSION))
	_, e := s.ChannelMessageSendEmbed(m.ChannelID, em.MessageEmbed)
	if e != nil {
		return e
	}
	return nil
}
