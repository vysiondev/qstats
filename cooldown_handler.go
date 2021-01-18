package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/vysiondev/qstats-go/err"
	"strconv"
	"time"
)

func (b *BaseHandler) HandleCooldown(id string, customMessage string, ctx context.Context) error {
	v, e := b.RedisClient.Get(ctx, id).Result()
	if e == redis.Nil {
		// Set new cooldown key since user never had it set.
		_, e = b.RedisClient.Set(ctx, id, makeTimestamp(), time.Millisecond*time.Duration(b.Config.Bot.Cooldown)).Result()
		if e != nil {
			return &err.SafeError{Message: "Cannot set new cooldown. " + e.Error()}
		}
		return nil
	} else if e != nil {
		return &err.SafeError{Message: "Cannot determine cooldown length. " + e.Error()}
	}
	// Cooldown exists, return a cooldown error
	i, e := strconv.ParseInt(v, 10, 64)
	if e != nil {
		return &err.SafeError{Message: "Cooldown value is not a UNIX timestamp."}
	}
	return &err.CooldownError{
		Message:  customMessage,
		TimeLeft: (i - makeTimestamp()) + int64(b.Config.Bot.Cooldown),
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
