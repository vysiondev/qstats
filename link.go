package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gocql/gocql"
	"github.com/vysiondev/qstats-go/bot_constants"
	"github.com/vysiondev/qstats-go/err"
	"github.com/vysiondev/qstats-go/utils"
	"strings"
)

func (b *BaseHandler) ExecuteLinkCommand(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, d RunDetails) error {
	if len(d.Args) == 0 {
		return &err.SafeError{Message: "You need to provide a Quaver user's username to link an account to."}
	}
	u, reqErr := SearchUser(strings.Join(d.Args, " "))
	if reqErr != nil {
		return reqErr
	}

	if dbErr := b.Db.Query(`UPDATE users SET quaverid = ?, prefer_7k = ? WHERE discordid = ?`,
		u.Users[0].ID,
		d.Is7K,
		m.Author.ID,
	).Consistency(gocql.One).Exec(); dbErr != nil {
		return dbErr
	}

	_, e := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s Linked to Quaver user **%s**. Your profile (accessed from \"profile\" command) will now show **%s** stats by default.", bot_constants.OkEmoji, u.Users[0].Username, utils.GetKeymodeString(d.Is7K)))
	if e != nil {
		return e
	}
	return nil
}
