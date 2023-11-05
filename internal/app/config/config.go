package config

type Config struct {
	Port       uint
	ConnString string
}

// TODO: make config with envs/yaml
func New() *Config {
	return &Config{
		Port:       7001,
		ConnString: "host=postgres port=5432 user=postgres password=postgres",
	}
}
