package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/vysiondev/qstats-go/err"
	"github.com/vysiondev/qstats-go/quaverapi_structs"
	"github.com/vysiondev/qstats-go/utils"
	"strconv"
	"strings"
)

func (b *BaseHandler) SetMapsInConversation(guildID string, is7K bool, mapIDs []int, ctx context.Context) error {
	// Index 0 will always be the keymode as a string.
	var mapIDsAsString []string
	mapIDsAsString = append(mapIDsAsString, utils.GetKeymodeIntAsStr(is7K))
	for i, s := range mapIDs {
		if i >= 5 {
			break
		}
		mapIDsAsString = append(mapIDsAsString, strconv.Itoa(s))
	}
	_, redisSetErr := b.RedisClient.Set(ctx, "map_"+guildID, strings.Join(mapIDsAsString, ","), 0).Result()
	if redisSetErr != nil {
		return redisSetErr
	}
	return nil
}

func (b *BaseHandler) GetMapsInConversation(guildID string, ctx context.Context) (*quaverapi_structs.MapsInConversation, error) {
	v, e := b.RedisClient.Get(ctx, "map_"+guildID).Result()
	if e == redis.Nil {
		return nil, &err.SafeError{Message: "No maps were found in the server conversation."}
	} else if e != nil {
		return nil, e
	}
	// Element at index 0 is the indication of keymode
	splitValue := strings.Split(v, ",")
	return &quaverapi_structs.MapsInConversation{
		GuildID: guildID,
		Maps:    splitValue[1:],
		Is7K:    splitValue[0] == "2",
	}, nil
}
