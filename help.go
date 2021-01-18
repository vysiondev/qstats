package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/vysiondev/qstats-go/err"
	"github.com/vysiondev/qstats-go/utils"
	"strings"
)

func (b *BaseHandler) ExecuteHelpCommand(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, rd RunDetails) error {
	embed := utils.CreateEmbed()

	cmdStr := make([]string, len(b.CommandList))
	for i, c := range b.CommandList {
		cmdStr[i] = " `" + c.Name + "` |"
	}

	if len(rd.Args) == 0 {
		embed.AddTitle("Help")
		embed.AddDescription(fmt.Sprintf("Commands are prefixed with `%s`.\n\n- To view the stats in 7K, append `--7k` to the end of the message.\n- To get more help with a command, type `%shelp [command name]`.\n- Type `-p [page number]` in your message to navigate to different pages.\n- [Join the support server](%s) for more help.\n\n**Command list**\n|%s",
			b.Config.Bot.Prefix,
			b.Config.Bot.Prefix,
			b.Config.Links.SupportServer,
			strings.Join(cmdStr, " "),
		))
	} else {
		commandToUse := b.FindCommandIndex(rd.Args[0])
		if commandToUse == nil {
			return &err.SafeError{Message: "Command not found"}
		}
		embed.AddTitle(b.CommandList[*commandToUse].Name)
		embed.AddDescription(fmt.Sprintf("%s\n\nYou can also use these shorthands: *%s*",
			b.CommandList[*commandToUse].Description,
			strings.Join(b.CommandList[*commandToUse].Shorthands, ", "),
		))
	}

	_, e := s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
	if e != nil {
		return e
	}
	return nil
}
