package config

import (
	"bytes"
	"fmt"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Target struct {
	Host       string `mapstructure:"host"`
	Protocol   string `mapstructure:"protocol"`
	Deployment string `mapstructure:"deployment"`
	Timeout    struct {
		Forward int `mapstructure:"forward"`
		Ping    int `mapstructure:"ping"`
	} `mapstructure:"timeout"`
}

type Config struct {
	Host   string `mapstructure:"host"`
	Target Target `mapstructure:"target"`
}

var defaultConfig = []byte(`
Host: :8080
Target:
  Host: example.com
  Protocol: tcp
  Deployment: example
	Timeout:
	  Ping: 1000
	  Forward: 5000
`)

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/")
	viper.AddConfigPath(".")
	viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	viper.ReadInConfig()
	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Panic().
			Err(err).
			Msg("Fatal error reading config file")
	}
	log.Info().Str("config", fmt.Sprintf("%+v", cfg)).
		Msg("config loaded")
	return &cfg
}
