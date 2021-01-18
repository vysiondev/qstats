package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gocql/gocql"
	"github.com/vysiondev/qstats-go/bot_constants"
	"github.com/vysiondev/qstats-go/err"
	"github.com/vysiondev/qstats-go/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type User struct {
	DiscordID string
	Prefer7K  bool
	QuaverID  int
}

func getIdFromUrl(content string) int {
	re := regexp.MustCompile("[0-9]+")
	if strings.Contains(content, "https://quavergame.com/mapset/") {
		str := re.FindString(content)
		if len(str) == 0 {
			return -1
		} else {
			returnInt, convertErr := strconv.Atoi(str)
			if convertErr != nil {
				return -1
			} else {
				return returnInt
			}
		}
	} else {
		return -1
	}
}

func (b *BaseHandler) HandleWithContext(ctx context.Context, args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	mapID := getIdFromUrl(m.Content)
	if mapID != -1 {
		if strings.Contains(m.Content, "https://quavergame.com/mapset/") {
			if strings.Contains(m.Content, "https://quavergame.com/mapset/map/") {
				mapPreviewErr := b.HandleCooldown(m.Author.ID, "Cannot display map preview for another %dms.", ctx)
				if mapPreviewErr != nil {
					PrintErrorToChannel(s, m.ChannelID, mapPreviewErr)
					return
				}
				e := b.AutorespondToMap(s, m.Message, mapID, ctx)
				if e != nil {
					PrintErrorToChannel(s, m.ChannelID, e)
					return
				}
				return
			} else {
				mapsetPreviewErr := b.HandleCooldown(m.Author.ID, "Cannot display mapset preview for another %dms.", ctx)
				if mapsetPreviewErr != nil {
					PrintErrorToChannel(s, m.ChannelID, mapsetPreviewErr)
					return
				}
				e := AutorespondToMapset(s, m.Message, mapID)
				if e != nil {
					PrintErrorToChannel(s, m.ChannelID, e)
					return
				}
			}
		}
	}

	if !strings.HasPrefix(args[0], b.Config.Bot.Prefix) {
		return
	}
	args = removePrefix(args, b.Config.Bot.Prefix)

	rdChan := make(chan RunDetails, 1)
	MakeRunDetails(args, rdChan)

	runDetails := <-rdChan
	close(rdChan)

	cmdIndex := b.FindCommandIndex(args[0])
	if cmdIndex == nil {
		return
	}

	if b.CommandList[*cmdIndex].OwnerOnly && m.Author.ID != b.Config.Bot.OwnerID {
		return
	}

	// Immediately process command cooldown handling here.
	cmdUsage := b.HandleCooldown(m.Author.ID, "Wait %dms before using another command.", ctx)
	if cmdUsage != nil {
		PrintErrorToChannel(s, m.ChannelID, cmdUsage)
		return
	}

	var user User
	if dbErr := b.Db.Query(`SELECT * FROM users WHERE discordid = ?`, m.Author.ID).Consistency(gocql.One).Scan(&user.DiscordID, &user.Prefer7K, &user.QuaverID); dbErr != nil {
		if dbErr == gocql.ErrNotFound {
			user = User{QuaverID: 0}

		} else {
			PrintErrorToChannel(s, m.ChannelID, dbErr)
			return
		}
	}

	if b.CommandList[*cmdIndex].RequiresUser == true {
		if len(runDetails.Args) > 0 {
			u, searchUserErr := SearchUser(strings.Join(runDetails.Args, " "))
			if searchUserErr != nil {
				PrintErrorToChannel(s, m.ChannelID, searchUserErr)
				return
			}
			runDetails.QuaverID = u.Users[0].ID
		} else {
			if user.QuaverID != 0 {
				runDetails.QuaverID = user.QuaverID
			} else {
				PrintErrorToChannel(s, m.ChannelID, &err.SafeError{Message: "You need to provide a user or link your Quaver account (q;link) to run this command!"})
				return
			}
		}
	}

	// always use user's preferred keymode when they are linked, command is "profile", and they didn't override 4k/7k
	if (!utils.Contains(runDetails.Options, "7") && !utils.Contains(runDetails.Options, "4")) && user.QuaverID != 0 && b.CommandList[*cmdIndex].Name == "profile" {
		runDetails.Is7K = user.Prefer7K
	} else {
		// else check for option
		runDetails.Is7K = utils.Contains(runDetails.Options, "7")
	}

	res := b.CommandList[*cmdIndex].RunFunction(ctx, s, m, runDetails)
	if res != nil {
		PrintErrorToChannel(s, m.ChannelID, res)
		return
	}
}

// All errors bubble up here. Ignore context.Canceled.
func PrintErrorToChannel(s *discordgo.Session, cid string, thisError error) {
	if thisError == context.DeadlineExceeded || thisError == context.Canceled {
		return
	}
	var msgToSend string
	switch serr := thisError.(type) {
	case *err.ReadError:
		msgToSend = bot_constants.ReadErrEmoji + " **Read Error:** " + serr.Error() + "\n*Either the Quaver API sent an incomplete body, or the API has changed.*"
		break
	case *err.SafeError:
		msgToSend = bot_constants.ErrorEmoji + " " + serr.Error()
		break
	case *err.CooldownError:
		msgToSend = bot_constants.CooldownEmoji + " **Cooldown:** " + fmt.Sprintf(serr.Error(), serr.TimeLeft)
		break
	default:
		msgToSend = bot_constants.FatalErrEmoji + " **Unexpected Error:** " + serr.Error() + "\n*This error should not have happened. Consider reporting it to the developer.*"

	}
	_, sendErr := s.ChannelMessageSend(cid, msgToSend)
	if sendErr != nil {
		switch sendErr2 := sendErr.(type) {
		case *discordgo.RESTError:
			if sendErr2.Message.Code == discordgo.ErrCodeMissingPermissions {
				return
			}
			break
		default:
			fmt.Println(sendErr2.Error())
		}
	}
}

func (b *BaseHandler) HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	if len(m.Message.GuildID) == 0 {
		return
	}
	if len(m.Content) == 0 {
		return
	}

	args := strings.Fields(m.Content)

	// Start global context that will persist throughout the entire command.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	b.HandleWithContext(ctx, args, s, m)

}

func removePrefix(a []string, prefix string) []string {
	a[0] = strings.TrimPrefix(a[0], prefix)
	return a
}
