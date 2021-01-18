package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *BaseHandler) ExecuteInviteCommand(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ RunDetails) error {
	_, e := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&permissions=0&scope=bot", s.State.User.ID))
	if e != nil {
		return e
	}
	return nil
}
