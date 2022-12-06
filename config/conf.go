package config

import (
	"bytes"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Target struct {
	Host       string
	Protocol   string
	Deployment string
	Timeout    struct {
		Forward int
		Ping    int
	}
}

type Config struct {
	Host   string
	Target Target
}

var defaultConfig = []byte(`
Host: localhost:8080
Target:
  Host: example.com
  Protocol: tcp
  Deployment: example
	Timeout:
	  Ping: 1000
	  Forward: 5000
`)

func LoadConfig() (cfg *Config) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/")
	viper.AddConfigPath(".")
	viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Panic().
			Err(err).
			Msg("Fatal error reading config file")
	}
	return
}
