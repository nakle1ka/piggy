package postgres

type PostgresConfig struct {
	Host     string
	User     string
	Password string
	DBname   string
	Port     int
	Sslmode  bool
}

type opt = func(p *PostgresConfig)

func defaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBname:   "postgres",
		Sslmode:  false,
	}
}

func WithHost(host string) opt {
	return func(p *PostgresConfig) {
		p.Host = host
	}
}

func WithUser(user string) opt {
	return func(p *PostgresConfig) {
		p.User = user
	}
}

func WithPassword(password string) opt {
	return func(p *PostgresConfig) {
		p.Password = password
	}
}

func WithDbname(dbname string) opt {
	return func(p *PostgresConfig) {
		p.DBname = dbname
	}
}

func WithPort(port int) opt {
	return func(p *PostgresConfig) {
		p.Port = port
	}
}

func WithSslmode(sslmode bool) opt {
	return func(p *PostgresConfig) {
		p.Sslmode = sslmode
	}
}
