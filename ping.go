package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *BaseHandler) ExecutePingCommand(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ RunDetails) error {
	_, e := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s Heartbeat acknowledgement took **%d**ms last interval / %s Shard **%d**",
		b.Config.Emoji.Ping,
		s.LastHeartbeatAck.Sub(s.LastHeartbeatSent).Milliseconds(),
		b.Config.Emoji.Ws,
		s.ShardID),
	)
	if e != nil {
		return e
	}
	return nil
}
