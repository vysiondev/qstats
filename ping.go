package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/vysiondev/qstats-go/bot_constants"
)

func (b *BaseHandler) ExecutePingCommand(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ RunDetails) error {
	_, e := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s Heartbeat acknowledgement took **%d**ms last interval / %s Shard **%d**",
		bot_constants.PingEmoji,
		s.LastHeartbeatAck.Sub(s.LastHeartbeatSent).Milliseconds(),
		bot_constants.WsEmoji,
		s.ShardID),
	)
	if e != nil {
		return e
	}
	return nil
}
