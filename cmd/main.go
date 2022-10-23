package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"someProject/pkg/tg"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path",
		"configs/config.toml", "path to config file")
}
func main() {
	flag.Parse()

	config := tg.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Print(err)
	}
	MakeMapInConfig(config.Addresses, config.AgentsName, config)

	tg.Bot(config)
}

func MakeMapInConfig(adr []string, agents []string, c *tg.Config) {
	i := 0
	for _, s := range adr {
		c.Offices[i+1] = &tg.Office{
			AgentName: agents[i],
			Address:   s,
			ChatID:    0,
		}
		i++
	}

}
