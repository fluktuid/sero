package config

import (
	"bytes"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Host       string
	TargetHost string
	Protocol   string
}

var defaultConfig = []byte(`
Host: localhost:8080
TargetHost: example.com
Protocol: tcp
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
