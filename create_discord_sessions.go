package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
	"log"
	"sync"
)

type SessionHolder struct {
	Sessions []*discordgo.Session
}

// Derived from: https://github.com/iopred/bruxism/blob/master/discord.go#L300
// i was having a shit ton of issues figuring out how to spawn multiple shards, this code seemed to do it properly.
func (s *SessionHolder) CreateDiscordSessions(minShard int, maxShard int, shardCount int, session *gocql.Session, r *redis.Client, c Configuration) error {
	numSessions := (maxShard - minShard) + 1
	s.Sessions = make([]*discordgo.Session, numSessions)
	wg := sync.WaitGroup{}

	baseHandler := NewBaseHandler(session, r, c)
	baseHandler.CommandList = baseHandler.FetchCommandList()

	for i := 0; i < numSessions; i++ {
		log.Printf("Starting shard %d\n", minShard+i)
		session, err := discordgo.New("Bot " + baseHandler.Config.Bot.Token)
		if err != nil {
			return err
		}
		s.Sessions[i] = session

		session.StateEnabled = false
		session.SyncEvents = false
		session.ShardID = minShard + i
		session.ShardCount = shardCount
		session.Identify.Properties.Browser = "Discord iOS"

		session.AddHandler(baseHandler.HandleReady)
		session.AddHandler(baseHandler.HandleMessageCreate)
		session.AddHandler(baseHandler.HandleGuildCreate)
		session.AddHandler(baseHandler.HandleGuildDelete)

		wg.Add(1)
		go func(session *discordgo.Session) {
			defer wg.Done()
			err := session.Open()
			if err != nil {
				log.Printf("error opening shard %s", err)
			}
		}(s.Sessions[i])
	}
	wg.Wait()
	return nil
}
