package main

import (
	"github.com/bwmarrin/discordgo"
)

func (_ *BaseHandler) HandleGuildCreate(s *discordgo.Session, r *discordgo.GuildCreate) {
	s.State.Guilds = append(s.State.Guilds, &discordgo.Guild{ID: r.ID})
}

func (_ *BaseHandler) HandleGuildDelete(s *discordgo.Session, r *discordgo.GuildDelete) {
	for i, g := range s.State.Guilds {
		if g.ID == r.ID {
			s.State.Guilds = remove(s.State.Guilds, i)
		}
	}
}

func remove(s []*discordgo.Guild, i int) []*discordgo.Guild {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}
