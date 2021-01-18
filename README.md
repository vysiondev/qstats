# QStats

The ultimate Quaver Discord bot. Retrieve scores, profiles and more without ever leaving Discord.

### Run

- You need a Redis database set up somewhere as well as a Cassandra database cluster. Optionally, you can setup Prometheus (see further down).
    - Cassandra: create a keyspace, and create a table using the following query inside it.
    - `CREATE TABLE users ( discordid text, prefer_7k boolean, quaverid int, primary key(discordid) );`
- You need to rename conf/example.yml to conf/config.yml, and set all the necessary values.
- Build the image (preferably tag it as "qstats")
- Run `docker run -d -p 127.0.0.1:8080:8080 --name qstatsbot qstats`.
    - No metrics: don't specify `-p` option

When launched, get Prometheus metrics by scraping localhost:8080/metrics in the scope of the container.

If you don't want metrics you can set `metrics: false` in the config.