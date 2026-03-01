package redis

type redisConfig struct {
	Host     string
	Port     int
	Password string
	Db       int
}

type opt = func(*redisConfig)

func defaultRedisConfig() redisConfig {
	return redisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Db:       0,
	}
}

func WithHost(h string) opt {
	return func(c *redisConfig) {
		c.Host = h
	}
}

func WithPort(p int) opt {
	return func(c *redisConfig) {
		c.Port = p
	}
}

func WithPassword(p string) opt {
	return func(c *redisConfig) {
		c.Password = p
	}
}

func WithDb(d int) opt {
	return func(c *redisConfig) {
		c.Db = d
	}
}
