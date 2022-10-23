package tg

type Config struct {
	Token      string   `toml:"token"`
	AgentsName []string `toml:"agents"`
	Addresses  []string `toml:"addresses"`
	Offices    map[int]*Office
}
type Office struct {
	AgentName string
	Address   string
	ChatID    int64
}

func NewConfig() *Config {
	return &Config{
		Token:      "",
		AgentsName: nil,
		Addresses:  nil,
		Offices:    make(map[int]*Office, 10),
	}
}
