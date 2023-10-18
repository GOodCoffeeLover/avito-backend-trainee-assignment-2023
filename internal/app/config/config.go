package config

type Config struct {
	Port       uint
	ConnString string
}

func New() *Config {
	return &Config{
		Port:       7001,
		ConnString: "host=localhost port=5432 user=postgres password=postgres",
	}
}
