package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *BaseHandler) HandleReady(s *discordgo.Session, r *discordgo.Ready) {
	b.SendMessageToWebhook(fmt.Sprintf("Shard %d ready: serving %d guilds",
		s.ShardID,
		len(r.Guilds),
	))
}
