package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
)

type BaseHandler struct {
	Db          *gocql.Session
	RedisClient *redis.Client
	CommandList []Command
	Config      Configuration
}

func (b *BaseHandler) FetchCommandList() []Command {
	return []Command{
		{
			Name:         "profile",
			Description:  "Show a player's general stats, including global/country rank, pass/fail count, judgement ratios, total score, and some more interesting data. Use `--7k` to show 7K stats.",
			RequiresUser: true,
			Shorthands:   []string{"p"},
			RunFunction:  b.ExecuteProfileCommand,
		},
		{
			Name:        "ping",
			Description: "Check to see if the bot is alive.",
			RunFunction: b.ExecutePingCommand,
		},
		{
			Name:        "about",
			Description: "Get information about QStats",
			Shorthands:  []string{"info"},
			RunFunction: b.ExecuteAboutCommand,
		},
		{
			Name:        "invite",
			Description: "Invite QStats to your own server",
			Shorthands:  []string{"i"},
			RunFunction: b.ExecuteInviteCommand,
		},
		{
			Name:        "support",
			Description: "Get a link to the support server",
			Shorthands:  []string{"server"},
			RunFunction: b.ExecuteSupportCommand,
		},
		{
			Name:        "link",
			Shorthands:  []string{"l"},
			Description: "Link a Quaver username to your Discord account so commands can be run without having to provide a username. Ex: `q;top vysion` can now just be run as q;top. Use the `--7k` option to show 7K stats by default for all commands.",
			RunFunction: b.ExecuteLinkCommand,
		},
		{
			Name:        "leaderboard",
			Shorthands:  []string{"lb"},
			Description: "Get the Quaver leaderboards in increments of 25 players. You can also target specific countries, such as typing `q;leaderboard us` for United States leaderboards. Use `-p [page number]` to navigate to different pages. Example: `q;leaderboard kr -p 3` will show South Korea's leaderboard for player rankings #51-75. Use `--7k` at the end to show 7K stats.",
			RunFunction: b.ExecuteLeaderboardCommand,
		},
		{
			Name:         "top",
			RequiresUser: true,
			Shorthands:   []string{"t"},
			Description:  "Show a player's top 5 scores by default. You can use `-p [page number]` to navigate further back. Ex: `q;t vysion -p 10` shows this player's 10th best score and so on.",
			RunFunction:  b.ExecuteTopCommand,
		},
		{
			Name:         "recent",
			RequiresUser: true,
			Shorthands:   []string{"r"},
			Description:  "Show a player's most recent score. You can use `-p [page number]` to navigate further back. Ex: `q;r player -p 10` shows this player's 10th most recent score and so on.",
			RunFunction:  b.ExecuteRecentCommand,
		},
		{
			Name:         "compare",
			RequiresUser: true,
			Shorthands:   []string{"c"},
			Description:  "When a map/list of maps are posted in the channel, you can use this command to list your best score on that map.",
			RunFunction:  b.ExecuteCompareCommand,
		},
		{
			Name:        "help",
			Shorthands:  []string{"h"},
			Description: "Get help with commands.",
			RunFunction: b.ExecuteHelpCommand,
		},
	}
}

func NewBaseHandler(db *gocql.Session, r *redis.Client, c Configuration) *BaseHandler {
	return &BaseHandler{
		Db:          db,
		Config:      c,
		RedisClient: r,
	}
}
