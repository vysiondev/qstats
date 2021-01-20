package main

type Configuration struct {
	Prometheus PrometheusConfig
	Cassandra  CassandraConfig
	Links      LinksConfig
	Bot        BotConfig
	Webhook    WebhookConfig
	Redis      RedisConfig
	Metrics    bool
	Emoji      EmojiConfig
}

type EmojiConfig struct {
	Success      string
	Ws           string
	Ping         string
	Error        EmojiErrorConfig
	Cooldown     string
	Offline      string
	Online       string
	Grade        EmojiGradeConfig
	FullCombo    string
	RankedStatus EmojiRankedStatusConfig
	Download     string
}

type EmojiRankedStatusConfig struct {
	NotSubmitted string
	Unranked     string
	Ranked       string
	Dan          string
}

type EmojiGradeConfig struct {
	F  string
	D  string
	C  string
	B  string
	A  string
	S  string
	SS string
	X  string
}

type EmojiErrorConfig struct {
	Normal string
	Read   string
	Fatal  string
}

type PrometheusConfig struct {
	Port string
}

type CassandraConfig struct {
	Hosts          []string
	Keyspace       string
	Authentication CassandraAuthenticationConfig
}

type CassandraAuthenticationConfig struct {
	Username string
	Password string
}

type LinksConfig struct {
	Github        string
	SupportServer string
	Website       string
}

type BotConfig struct {
	Token    string
	OwnerID  string
	Shards   int
	Cooldown int
	Prefix   string
}

type WebhookConfig struct {
	URL string
}

type RedisConfig struct {
	URI      string
	Password string
	Db       int
}
