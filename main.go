package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	guildCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "qstats_go_guild_count",
		Help: "# of guilds on this container",
	})
	heapUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "qstats_go_heap_usage",
		Help: "Memory in use by this process in bytes",
	},
	)
)

func main() {
	prometheus.MustRegister(guildCount)
	prometheus.MustRegister(heapUsage)

	viper.SetConfigName("config")
	viper.AddConfigPath("./conf/")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	var configuration Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Set undefined variables
	viper.SetDefault("prometheus.port", "8080")
	viper.SetDefault("links.supportserver", "https://discord.gg/F7RBKh2")
	viper.SetDefault("bot.prefix", "q;")
	viper.SetDefault("bot.cooldown", 4000)
	viper.SetDefault("bot.shards", 1)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("metrics", false)

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	if len(configuration.Bot.OwnerID) == 0 || len(configuration.Bot.Token) == 0 {
		log.Fatal("You did not set your ID as the owner ID (bot.ownerid) or bot token (bot.token) in your config!")
	}

	cluster := gocql.NewCluster(configuration.Cassandra.Hosts...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: configuration.Cassandra.Authentication.Username,
		Password: configuration.Cassandra.Authentication.Password,
	}
	cluster.Keyspace = configuration.Cassandra.Keyspace
	cluster.Consistency = gocql.Quorum
	session, e := cluster.CreateSession()
	if e != nil {
		log.Fatal("Failed to open Cassandra session:", e)
	}
	defer session.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     configuration.Redis.URI,
		Password: configuration.Redis.Password,
		DB:       configuration.Redis.Db,
	})

	holder := SessionHolder{}
	err = holder.CreateDiscordSessions(0, configuration.Bot.Shards-1, configuration.Bot.Shards, session, redisClient, configuration)
	if err != nil {
		log.Fatal("Couldn't start bot:", err)
	}

	// user wants metrics
	if configuration.Metrics {
		go func() {
			for {
				time.Sleep(time.Second * 5)
				for _, s := range holder.Sessions {
					guildCount.Set(float64(len(s.State.Guilds)))
				}
				var mem runtime.MemStats
				runtime.ReadMemStats(&mem)
				heapUsage.Set(float64(mem.Alloc))
			}
		}()

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			log.Fatal(http.ListenAndServe(":"+configuration.Prometheus.Port, nil))
		}()
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	for _, s := range holder.Sessions {
		_ = s.Close()
	}
}
