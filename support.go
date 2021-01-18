package main

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

func (b *BaseHandler) ExecuteSupportCommand(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ RunDetails) error {
	_, e := s.ChannelMessageSend(m.ChannelID, b.Config.Links.SupportServer)
	if e != nil {
		return e
	}
	return nil
}
