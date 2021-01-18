package main

type Configuration struct {
	Prometheus PrometheusConfig
	Cassandra  CassandraConfig
	Links      LinksConfig
	Bot        BotConfig
	Webhook    WebhookConfig
	Redis      RedisConfig
	Metrics    bool
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
